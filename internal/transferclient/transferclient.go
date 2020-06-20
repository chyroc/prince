package transferclient

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/chyroc/prince/internal/helper"
	"github.com/chyroc/prince/internal/transferserver"
)

func RunTransferClient(transferHost string) {
	conn, err := net.Dial("tcp", transferHost)
	if err != nil {
		fmt.Printf("[client] connecting to %s failed: %s\n", transferHost, err)
		return
	}
	defer conn.Close()

	fmt.Printf("[client] connecting to %s\n", transferHost)

	if err := handleWrite(conn); err != nil {
		fmt.Printf("[client] handler client connect failed: %s\n", err)
	}
}

func handleWrite(conn net.Conn) error {
	buf := bufio.NewReader(conn)

	bs, _, err := buf.ReadLine()
	if err != nil {
		return fmt.Errorf("[client] read line failed: %s\n", err)
	}

	fmt.Printf("[client] get request: %s\n", string(bs))

	task := new(transferserver.Task)
	if err := json.Unmarshal(bs, task); err != nil {
		return fmt.Errorf("[client] unmarshal request %s failed: %s\n", bs, err)
	}

	switch task.Type {
	case "http":
		return proxyHTTP(task, conn)
	case "https":
		return proxyHTTPs(task, conn)
	default:
		return fmt.Errorf("[client] %s is invalid proxy type", task.Type)
	}
}

func proxyHTTP(task *transferserver.Task, conn net.Conn) error {
	defer conn.Close()

	req, err := http.NewRequest(task.Method, task.URL, bytes.NewReader(task.Body))
	if err != nil {
		return fmt.Errorf("[client] new req failed: %s", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("[client] do req failed: %s", err)
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		if len(v) > 0 {
			_, _ = conn.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v[0])))
		}
	}
	_, _ = conn.Write([]byte("\r\n"))
	bs, _ := ioutil.ReadAll(resp.Body)
	_, err = conn.Write(bs)
	return err
}

func proxyHTTPs(task *transferserver.Task, conn net.Conn) error {
	dest_conn, err := net.DialTimeout("tcp", task.Host, 10*time.Second)
	if err != nil {
		return fmt.Errorf("[client] dial to %s failed: %s", task.Host, err)
	}

	helper.TransferWait(conn, dest_conn)

	return nil
}
