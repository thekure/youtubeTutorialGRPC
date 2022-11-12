package service

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	ytTut "github.com/thekure/youtubeTutorialGRPC/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LaptopServer is the server that provides laptop services
type LaptopServer struct {
	Store LaptopStore
	ytTut.UnimplementedLaptopServiceServer
}

// Returns a new LaptopServer
func NewLaptopServer(store LaptopStore) *LaptopServer {
	return &LaptopServer{store, ytTut.UnimplementedLaptopServiceServer{}}
}

// Unary RPC to create a new laptop.
func (server *LaptopServer) CreateLaptop(
	ctx context.Context,
	request *ytTut.CreateLaptopRequest,
) (*ytTut.CreateLaptopResponse, error) {

	laptop := request.GetLaptop()
	log.Printf("received a create-laptop request with id: %s", laptop.Id)

	if len(laptop.Id) > 0 {
		// Check if it's a valid UUID or not.
		_, err := uuid.Parse(laptop.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "laptop ID is not a valid UUID: %v", err)
		}
	} else {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot generate a new laptop ID: %v", err)
		}
		laptop.Id = id.String()
	}

	// This is when one would probably add the laptop to a database... But not in this tutorial. We use in memory instead.

	// Checks that the client hasn't cancelled the request before saving to database
	if ctx.Err() == context.Canceled {
		log.Print("request is cancelled")
		return nil, status.Error(codes.Canceled, "request is cancelled")
	}

	// The following if statement checks that the request is handled within 5 seconds,
	// otherwise the laptop is not saved, because the request got canceled from client side
	if ctx.Err() == context.DeadlineExceeded {
		log.Print("deadline is exceeded")
		return nil, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}

	err := server.Store.Save(laptop)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, ErrAlreadyExists) {
			code = codes.AlreadyExists
		}
		return nil, status.Errorf(code, "cannot save laptop to the store: %v", err)
	}

	log.Printf("saved laptop with id: %s", laptop.Id)

	response := &ytTut.CreateLaptopResponse{
		Id: laptop.Id,
	}

	return response, nil
}

// SearchLaptop is a server streaming RPC to search for laptops
func (server *LaptopServer) SearchLaptop(
	request *ytTut.SearchLaptopRequest,
	stream ytTut.LaptopService_SearchLaptopServer,
) error {
	filter := request.GetFilter()
	log.Printf("receive a search-laptop request with filter: %v", filter)

	err := server.Store.Search(
		stream.Context(),
		filter,
		func(laptop *ytTut.Laptop) error {
			response := &ytTut.SearchLaptopResponse{Laptop: laptop}

			err := stream.SendMsg(response)
			if err != nil {
				return err
			}

			log.Printf("sent laptop with id: %s", laptop.GetId())
			return nil
		},
	)

	if err != nil {
		return status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	return nil
}
