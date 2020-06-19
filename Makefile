build:
	go build -o ./bin/prince main.go

build-linux:
	GOOS=linux GOARCH=amd64 go build -o ./bin/prince main.go

gen_pb:
	protoc --go_out=plugins=grpc:./internal/pb_gen pb.proto
