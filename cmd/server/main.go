package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	ytTut "github.com/thekure/youtubeTutorialGRPC/grpc"
	"github.com/thekure/youtubeTutorialGRPC/service"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("SERVER PRINT")
	// We need a port for the server. Using flag.Int() we get it from commandline argument.

	port := flag.Int("port", 0, "the server port")
	flag.Parse()
	log.Printf("start server on port %d", *port)

	laptopServer := service.NewLaptopServer(service.NewInMemoryLaptopStore())
	grpcServer := grpc.NewServer()
	ytTut.RegisterLaptopServiceServer(grpcServer, laptopServer)

	// Create address string with the port we got before:
	address := fmt.Sprintf("0.0.0.0:%d", *port)
	// Listen for tcp connections on this server address:
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
