package main

import (
	"fmt"
	"github.com/chyroc/prince/internal/proxy"
)

func main() {
	fmt.Println("server")
	proxy.New(":8080") // 提供代理的 http 接口
	// 提供转发的 tcp 接口
}
