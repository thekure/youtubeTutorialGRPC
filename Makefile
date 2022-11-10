gen:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative grpc/processor_message.proto

run: 
	go run main.go

test:
	go test -cover -race ./...