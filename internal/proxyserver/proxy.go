package proxyserver

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/chyroc/prince/internal/helper"
	"github.com/chyroc/prince/internal/transferserver"
)

type Proxy struct{}

func (p *Proxy) serverHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	isFin := make(chan int)

	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("[proxy] read body failed: %s\n", err)
	}

	transferserver.PushTask(transferserver.Task{
		Type:    "http",
		Method:  r.Method,
		Host:    r.URL.Host,
		URL:     r.URL.String(),
		Headers: r.Header,
		Body:    bs,
		Callback: func(conn net.Conn) {
			defer conn.Close()

			buf := bufio.NewReader(conn)
			w.WriteHeader(http.StatusOK)
			for {
				line, _, err := buf.ReadLine()
				if err != nil {
					fmt.Printf("[proxy] read header failed: %s\n", err)
					continue
				}
				if string(line) == "" {
					break
				}
				n := strings.SplitAfterN(string(line), ":", 2)

				k := n[0]
				v := ""
				if len(n) >= 2 {
					v = n[1]
				}
				w.Header().Add(k, v)
			}
			_, _ = io.Copy(w, conn)

			close(isFin)
		},
	})
	<-isFin
}

func (p *Proxy) serverHTTPs(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	isFin := make(chan int)
	transferserver.PushTask(transferserver.Task{
		Type: "https",
		Host: r.URL.Host,
		Callback: func(conn net.Conn) {
			hijacker, ok := w.(http.Hijacker)
			if !ok {
				http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
				return
			}
			//通过Hijack可以让调用方接管连接
			client_conn, _, err := hijacker.Hijack()
			if err != nil {
				http.Error(w, err.Error(), http.StatusServiceUnavailable)
			}

			helper.TransferWait(conn, client_conn)

			close(isFin)
		},
	})
	<-isFin
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("[proxy] received request %s %s %s\n", req.Method, req.URL.String(), req.RemoteAddr, )

	if req.Method == "CONNECT" {
		p.serverHTTPs(w, req)
		return
	}
	p.serverHTTP(w, req)
}

func RunProxyServer(proxyHost string) {
	fmt.Printf("[proxy] 启动服务 %s ...\n", proxyHost)
	http.ListenAndServe(proxyHost, &Proxy{})
}
