package rpcclient

import (
	"google.golang.org/grpc"

	"github.com/chyroc/prince/internal/pb_gen"
)

var Client pb_gen.PrinceServiceClient

func Init(host string) {
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	// defer conn.Close()

	Client = pb_gen.NewPrinceServiceClient(conn)
}
