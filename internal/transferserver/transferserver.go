package transferserver

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

func RunTransferServer(transferHost string) error {
	fmt.Printf("[transfer] 启动服务 %s ...\n", transferHost)

	l, err := net.Listen("tcp", transferHost)
	if err != nil {
		return fmt.Errorf("listening on %s failed: %s\n", transferHost, err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("[transfer] accepting failed: %s\n", err)
			continue
		}
		fmt.Printf("[transfer] received  %s -> %s \n", conn.RemoteAddr(), conn.LocalAddr())
		go handleRequest(conn)
	}
}

type Task struct {
	Type     string              `json:"type"` // http. https
	Method   string              `json:"method"`
	Host     string              `json:"host"`
	URL      string              `json:"url"`
	Headers  http.Header         `json:"headers"`
	Body     []byte              `json:"body"`
	Callback func(conn net.Conn) `json:"-"`
}

func PushTask(task Task) {
	tasks <- task
}

var tasks = make(chan Task, 1000)

func handleRequest(conn net.Conn) {
	defer conn.Close()

	task := <-tasks
	bs, _ := json.Marshal(task)
	_, _ = conn.Write([]byte(string(bs) + "\n"))

	task.Callback(conn)
}
