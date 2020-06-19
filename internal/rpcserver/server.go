package rpcserver

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/chyroc/prince/internal/pb_gen"
)

func Run(transferHost string) error {
	fmt.Printf("[server][transfer] 启动转发服务: %s\n", transferHost)

	sv := grpc.NewServer()
	pb_gen.RegisterPrinceServiceServer(sv, new(Server))

	listener, err := net.Listen("tcp", transferHost)
	if err != nil {
		return fmt.Errorf("[server][transfer] 服务启动失败: %w", err)
	}
	return sv.Serve(listener)
}
