build:
	GOOS=linux GOARCH=amd64 go build -o ./bin/prince-server cmd/prince-server/main.go
	go build -o ./bin/prince-client cmd/prince-client/main.go

gen_pb:
	protoc --go_out=plugins=grpc:./internal/pb_gen pb.proto
