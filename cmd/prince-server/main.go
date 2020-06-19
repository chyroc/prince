package main

import (
	"fmt"
	"github.com/chyroc/prince/internal/proxy"
	"github.com/chyroc/prince/internal/rpcserver"
)

func main() {
	fmt.Println("server")
	go proxy.New(":8002") // 提供代理的 http 接口
	rpcserver.NewServer().Start()
	// 提供转发的 tcp 接口
}
