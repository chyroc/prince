package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"github.com/chyroc/prince/internal/pb_gen"
	"github.com/chyroc/prince/internal/rpcclient"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

var host string

func init() {
	flag.StringVar(&host, "host", "", "服务端地址")
	flag.Parse()
}

func main() {
	fmt.Println("client")

	if host == "" {
		panic("host is empty")
	}

	rpcclient.Init(host)

	stream, err := rpcclient.Client.HttpProxy(context.Background())
	if err != nil {
		panic(err)
	}

	for {
		req, err := stream.Recv()
		fmt.Println("client recv", req, err)
		if err != nil {
			logrus.Errorln(err)
			continue
		}

		req2, err := http.NewRequest(req.Method, req.Url, bytes.NewReader(req.Body))
		if err != nil {
			logrus.Errorln(err)
			continue
		}
		resp, err := http.DefaultClient.Do(req2)
		if err != nil {
			logrus.Errorln(err)
			continue
		}
		bs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logrus.Errorln(err)
			continue
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
			Uuid:    req.Uuid,
			Status:  int32(resp.StatusCode),
			Headers: headers,
			Body:    bs,
		})
		fmt.Println("client send", err)
		time.Sleep(time.Second)
	}
}
