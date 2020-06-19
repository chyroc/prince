package cmd

import (
	"fmt"

	"github.com/chyroc/prince/internal/proxyserver"
	"github.com/chyroc/prince/internal/rpcserver"
)

func RunServer(transferHost, proxyHost string) error {
	fmt.Println("[server] 启动服务端 ...")

	if transferHost == "" {
		return fmt.Errorf("请使用 --transfer_host=host:port 指定服务端启动的转发端口")
	}
	if proxyHost == "" {
		return fmt.Errorf("请使用 --proxy_host=host:port 指定服务启动的代理端口")
	}

	go proxyserver.Run(proxyHost) // 提供代理的 http 接口

	return rpcserver.Run(transferHost)
}
