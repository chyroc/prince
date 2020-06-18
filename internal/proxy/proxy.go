package proxy

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

type Proxy struct {
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("接受请求 %s %s %+v %s\n", req.Method, req.URL.String(), req.Header, req.RemoteAddr)

	outReq := new(http.Request)
	*outReq = *req // 这只是一个浅层拷贝

	//clientIP, _, err := net.SplitHostPort(req.RemoteAddr)
	//if err == nil {
	//	prior, ok := outReq.Header["X-Forwarded-For"]
	//	if ok {
	//		clientIP = strings.Join(prior, ", ") + ", " + clientIP
	//	}
	//	outReq.Header.Set("X-Forwarded-For", clientIP)
	//}


	//transport := http.DefaultTransport
	//res, err := transport.RoundTrip(outReq)
	//if err != nil {
	//	w.WriteHeader(http.StatusBadGateway) // 502
	//	return
	//}

	for key, value := range res.Header {
		for _, v := range value {
			w.Header().Add(key, v)
		}
	}

	w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)
	res.Body.Close()
}

func New(addr string) {
	fmt.Println("serve on " + addr)
	http.Handle("/", &Proxy{})
	http.ListenAndServe(addr, nil)
}
