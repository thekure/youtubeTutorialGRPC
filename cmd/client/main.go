package main

import (
	"context"
	"flag"
	"log"
	"time"

	ytTut "github.com/thekure/youtubeTutorialGRPC/grpc"
	"github.com/thekure/youtubeTutorialGRPC/sample"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func main() {
	serverAddress := flag.String("address", "", "the server address")
	flag.Parse()
	log.Printf("dial server %s", *serverAddress)

	insecure := insecure.NewCredentials()
	conn, err := grpc.Dial(*serverAddress, grpc.WithTransportCredentials(insecure))
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}

	laptopClient := ytTut.NewLaptopServiceClient(conn)
	laptop := sample.NewLaptop()
	// sets id to an empty string. this just makes the prints in server term more clear
	laptop.Id = ""
	request := &ytTut.CreateLaptopRequest{
		Laptop: laptop,
	}

	// set timeout for request using context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// above returns a context and a cancel object.
	// we want to defer the cancel call before exiting the main function
	defer cancel()

	response, err := laptopClient.CreateLaptop(ctx, request)
	if err != nil {
		status, ok := status.FromError(err)
		if ok && status.Code() == codes.AlreadyExists {
			// not a big deal
			log.Print("laptop already exists")
		} else {
			log.Fatal("cannot create laptop: ", err)
		}
		return
	}

	log.Printf("created laptop with id: %s", response.Id)
}
