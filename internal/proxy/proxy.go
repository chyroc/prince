package proxy

import (
	"bytes"
	"fmt"
	"github.com/chyroc/prince/internal/pb_gen"
	"github.com/chyroc/prince/internal/rpcserver"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
)

type Proxy struct {
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("接受请求 %s %s %+v %s\n", req.Method, req.URL.String(), req.Header, req.RemoteAddr)

	//outReq := new(http.Request)
	//*outReq = *req // 这只是一个浅层拷贝

	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
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
		logrus.Infoln("resp", resp)

		for key, value := range resp.Headers {
			w.Header().Set(key, value)
		}

		w.WriteHeader(int(resp.Status))
		_, err := io.Copy(w, bytes.NewReader(resp.Body))

		return err
	})
}

func New(addr string) {
	fmt.Println("serve on " + addr)
	http.Handle("/", &Proxy{})
	http.ListenAndServe(addr, nil)
}
