package proxyserver

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/chyroc/prince/internal/pb_gen"
	"github.com/chyroc/prince/internal/rpcserver"
)

type Proxy struct {
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("[server][proxy] 接受请求 %s %s %+v %s\n", req.Method, req.URL.String(), req.Header, req.RemoteAddr)

	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Printf("[server][proxy] 读取请求失败: %s\n", err)
		return
	}

	headers := make(map[string]string)
	for k, v := range req.Header {
		if len(v) >= 1 {
			headers[k] = v[0]
		} else {
			headers[k] = ""
		}
	}

	rpcserver.Send(pb_gen.HttpProxyRequest{
		Method:  req.Method,
		Url:     req.URL.String(),
		Headers: headers,
		Body:    bs,
	}, func(resp *pb_gen.HttpProxyResponse) error {
		for key, value := range resp.Headers {
			w.Header().Set(key, value)
		}

		w.WriteHeader(int(resp.Status))
		_, err := io.Copy(w, bytes.NewReader(resp.Body))

		return err
	})
}

func Run(addr string) {
	fmt.Printf("[server][proxy] 代理服务器启动: %s\n", addr)
	http.Handle("/", &Proxy{})
	http.ListenAndServe(addr, nil)
}
