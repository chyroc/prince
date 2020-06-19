package rpcserver

import (
	"fmt"
	"github.com/chyroc/prince/internal/pb_gen"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

var ServerChan = make(chan pb_gen.HttpProxyRequest, 1000)
var lock sync.RWMutex
var response = make(map[string]*pb_gen.HttpProxyResponse)

func (r *Server) HttpProxy(stream pb_gen.PrinceService_HttpProxyServer) error {
	go readResponse(stream)
	writeRequest(stream)
	return nil
}

func readResponse(stream pb_gen.PrinceService_HttpProxyServer) {
	for {
		resp, err := stream.Recv()
		if err != nil {
			logrus.Error(err)
			continue
		}

		logrus.Infoln(resp)

		lock.Lock()
		response[resp.Uuid] = resp
		lock.Unlock()
	}
}

func writeRequest(stream pb_gen.PrinceService_HttpProxyServer) {
	for v := range ServerChan {
		err := stream.Send(&v)
		logrus.Infoln("server send", err)
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
