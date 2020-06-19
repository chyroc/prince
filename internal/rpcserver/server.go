package rpcserver

import (
	"fmt"
	"github.com/chyroc/prince/internal/pb_gen"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	address string
}

func NewServer() *Server {
	return &Server{
		address: ":8001",
	}
}

func (r *Server) Start() error {
	logrus.Infof("start server at %s", r.address)

	sv := grpc.NewServer()
	pb_gen.RegisterPrinceServiceServer(sv, new(Server))

	listener, err := net.Listen("tcp", r.address)
	if err != nil {
		return fmt.Errorf("start server failed: %w", err)
	}
	return sv.Serve(listener)
}
