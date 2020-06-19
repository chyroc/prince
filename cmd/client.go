package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chyroc/prince/internal/pb_gen"
	"github.com/chyroc/prince/internal/rpcclient"
)

func RunClient(transferHost string) error {
	return runClient(transferHost, 0)
}

func runClient(transferHost string, times int) error {
	fmt.Println("[client] 启动客户端 ...")

	if transferHost == "" {
		return fmt.Errorf("请使用 --transfer_host=host:port 指定连接的服务端转发服务端口")
	}

	// 初始化
	rpcclient.Init(transferHost)

	// 连接到服务器
	stream, err := rpcclient.Client.HttpProxy(context.Background())
	if err != nil {
		return err
	}

	fmt.Printf("[client] 连接到服务端: %s\n", transferHost)

	// 监听服务下发的任务，然后执行
	for {
		req, err := stream.Recv()
		if err != nil {
			fmt.Printf("[client] recv 失败: %s\n", err)
			if status.Code(err) == codes.Unavailable {
				// 服务端断
				return fmt.Errorf("服务端断连")
			}
			continue
		}

		if err := handlerRequest(stream, req); err != nil {
			fmt.Printf("%s\n", err)
			continue
		}
	}
}

func handlerRequest(stream pb_gen.PrinceService_HttpProxyClient, request *pb_gen.HttpProxyRequest) error {
	fmt.Printf("[client] recv: %+v\n", request)

	req2, err := http.NewRequest(request.Method, request.Url, bytes.NewReader(request.Body))
	if err != nil {
		return fmt.Errorf("[client] 初始化请求错误: %s", err)
	}
	resp, err := http.DefaultClient.Do(req2)
	if err != nil {
		return fmt.Errorf("[client] 发出请求错误: %s", err)
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("[client] 读取返回值错误: %s", err)
	}

	headers := make(map[string]string)
	for k, v := range resp.Header {
		if len(v) >= 1 {
			headers[k] = v[0]
		} else {
			headers[k] = ""
		}
	}

	err = stream.Send(&pb_gen.HttpProxyResponse{
		Uuid:    request.Uuid,
		Status:  int32(resp.StatusCode),
		Headers: headers,
		Body:    bs,
	})
	if err != nil {
		return fmt.Errorf("[client] 发送到转发服务失败: %s", err)
	}

	return nil
}
