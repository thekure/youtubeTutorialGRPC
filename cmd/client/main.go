package main

import (
	"context"
	"flag"
	"io"
	"log"
	"time"

	ytTut "github.com/thekure/youtubeTutorialGRPC/grpc"
	"github.com/thekure/youtubeTutorialGRPC/sample"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func createLaptop(laptopClient ytTut.LaptopServiceClient) {
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

func searchLaptop(laptopClient ytTut.LaptopServiceClient, filter *ytTut.Filter) {
	log.Print("search filter: ", filter)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	request := &ytTut.SearchLaptopRequest{Filter: filter}
	stream, err := laptopClient.SearchLaptop(ctx, request)
	if err != nil {
		log.Fatal("cannot search laptop: ", err)
	}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatal("cannot receive response: ", err)
		}

		laptop := response.GetLaptop()
		log.Print("- found: ", laptop.GetId())
		log.Print(" + brand: ", laptop.GetBrand())
		log.Print(" + name: ", laptop.GetName())
		log.Print(" + cpu cores: ", laptop.GetCpu().GetNumberCores())
		log.Print(" + cpu min ghz: ", laptop.GetCpu().GetMinGhz())
		log.Print(" + ram: ", laptop.GetRam().GetValue(), laptop.GetRam().GetUnit())
		log.Print(" + price: ", laptop.GetPriceUsd(), " USD \n ")
	}

}

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
	for i := 0; i < 10; i++ {
		createLaptop(laptopClient)
	}

	filter := &ytTut.Filter{
		MaxPriceUsd: 3000,
		MinCpuCores: 4,
		MinCpuGhz:   2.5,
		MinRam:      &ytTut.Memory{Value: 8, Unit: ytTut.Memory_GIGABYTE},
	}

	searchLaptop(laptopClient, filter)
}
