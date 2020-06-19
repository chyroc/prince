package rpcserver

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chyroc/prince/internal/pb_gen"
)

var ServerChan = make(chan pb_gen.HttpProxyRequest, 1000)
var lock sync.RWMutex
var response = make(map[string]*pb_gen.HttpProxyResponse)

type Server struct {
}

func IsClosed(ch <-chan int) bool {
	select {
	case <-ch:
		return true
	default:
	}

	return false
}

func (r *Server) HttpProxy(stream pb_gen.PrinceService_HttpProxyServer) error {
	closeChan := make(chan int)
	go readResponse(stream, closeChan)
	writeRequest(stream, closeChan)
	return nil
}

func readResponse(stream pb_gen.PrinceService_HttpProxyServer, closeChan chan int) {
	for {
		resp, err := stream.Recv()
		if err != nil {
			fmt.Printf("[server][transfer] recv failed: %s\n", err)
			if status.Code(err) == codes.Canceled {
				// 代表客户端断开连接了，需要取消这个链接
				close(closeChan)
				return
			}
			continue
		}

		lock.Lock()
		response[resp.Uuid] = resp
		lock.Unlock()
	}
}

func writeRequest(stream pb_gen.PrinceService_HttpProxyServer, closeChan chan int) {
	for {
		select {
		case v := <-ServerChan:
			if err := stream.Send(&v); err != nil {
				fmt.Printf("[server][transfer] send failed: %s\n", err)
			}
		case <-closeChan:
			return
		}
	}
}

func Send(request pb_gen.HttpProxyRequest, handler func(proxyResponse *pb_gen.HttpProxyResponse) error) error {
	request.Uuid = uuid.New().String()
	ServerChan <- request

	// 10 ms 一次，持续 20s，也就是 20000/10 = 2000 次
	i := 0
	for {
		if i > 2000 {
			return fmt.Errorf("response not found")
		}
		i++
		lock.RLock()
		resp, ok := response[request.Uuid]
		lock.RUnlock()

		if ok {
			lock.Lock()
			delete(response, request.Uuid)
			lock.Unlock()

			return handler(resp)
		}

		time.Sleep(time.Second / 100)
	}
}
