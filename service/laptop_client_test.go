package service_test

import (
	"context"
	"io"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	ytTut "github.com/thekure/youtubeTutorialGRPC/grpc"
	"github.com/thekure/youtubeTutorialGRPC/sample"
	"github.com/thekure/youtubeTutorialGRPC/serializer"
	"github.com/thekure/youtubeTutorialGRPC/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestClientCreateLaptop(t *testing.T) {
	t.Parallel()

	laptopServer, serverAddress := startTestLaptopServer(t, service.NewInMemoryLaptopStore())
	laptopClient := newTestLaptopClient(t, serverAddress)

	laptop := sample.NewLaptop()
	expectedID := laptop.Id
	request := &ytTut.CreateLaptopRequest{
		Laptop: laptop,
	}

	response, err := laptopClient.CreateLaptop(context.Background(), request)
	require.NoError(t, err)
	require.NotNil(t, response)
	require.Equal(t, expectedID, response.Id)

	// Check that the laptop is saved to the store
	other, err := laptopServer.Store.Find(response.Id)
	require.NoError(t, err)
	require.NotNil(t, other)

	// check that the saved laptop is the same as the one we sent.
	requireSameLaptop(t, laptop, other)
}

func TestClientSearchLaptop(t *testing.T) {
	t.Parallel()

	filter := &ytTut.Filter{
		MaxPriceUsd: 2000,
		MinCpuCores: 4,
		MinCpuGhz:   2.2,
		MinRam:      &ytTut.Memory{Value: 8, Unit: ytTut.Memory_GIGABYTE},
	}

	store := service.NewInMemoryLaptopStore()
	expectedIDs := make(map[string]bool)

	for i := 0; i < 6; i++ {
		laptop := sample.NewLaptop()

		switch i {
		case 0:
			laptop.PriceUsd = 2500
		case 1:
			laptop.Cpu.NumberCores = 2
		case 2:
			laptop.Cpu.MinGhz = 2.0
		case 3:
			laptop.Ram = &ytTut.Memory{Value: 4096, Unit: ytTut.Memory_MEGABYTE}
		case 4:
			laptop.PriceUsd = 1999
			laptop.Cpu.NumberCores = 4
			laptop.Cpu.MinGhz = 2.5
			laptop.Cpu.MaxGhz = 4.5
			laptop.Ram = &ytTut.Memory{Value: 16, Unit: ytTut.Memory_GIGABYTE}
			expectedIDs[laptop.Id] = true
		case 5:
			laptop.PriceUsd = 2000
			laptop.Cpu.NumberCores = 6
			laptop.Cpu.MinGhz = 2.8
			laptop.Cpu.MaxGhz = 5.0
			laptop.Ram = &ytTut.Memory{Value: 64, Unit: ytTut.Memory_GIGABYTE}
			expectedIDs[laptop.Id] = true
		}

		err := store.Save(laptop)
		require.NoError(t, err)
	}

	_, serverAddress := startTestLaptopServer(t, store)
	laptopClient := newTestLaptopClient(t, serverAddress)

	request := &ytTut.SearchLaptopRequest{Filter: filter}
	stream, err := laptopClient.SearchLaptop(context.Background(), request)
	require.NoError(t, err)

	found := 0
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}

		require.NoError(t, err)
		require.Contains(t, expectedIDs, response.GetLaptop().GetId())

		found += 1
	}

}

func startTestLaptopServer(t *testing.T, store service.LaptopStore) (*service.LaptopServer, string) {
	laptopServer := service.NewLaptopServer(store)

	grpcServer := grpc.NewServer()
	ytTut.RegisterLaptopServiceServer(grpcServer, laptopServer)

	listener, err := net.Listen("tcp", ":0") // 0 means we want it assigned to any random available port
	require.NoError(t, err)

	go grpcServer.Serve(listener) // start listening for requests.
	// This is a "blocking call" and has to be run in a separate go routine, since we
	// want to send request to this server after that (literally what the guy said..)

	return laptopServer, listener.Addr().String()

	// return the laptopServer, and the address string of the listener.
}

func newTestLaptopClient(t *testing.T, serverAddress string) ytTut.LaptopServiceClient {
	// Dial the server address. This is just for testing, so we use the insecure connection.
	insecure := insecure.NewCredentials()

	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(insecure))
	require.NoError(t, err)
	return ytTut.NewLaptopServiceClient(conn)
	// return a new client with the created connection.
}

func requireSameLaptop(t *testing.T, laptop1 *ytTut.Laptop, laptop2 *ytTut.Laptop) {
	// For this test, we have to serialize the laptops to json to properly compare
	// otherwise the objects will be missing some internal fields, and be inequal.

	json1, err := serializer.ProtobufToJSON(laptop1)
	require.NoError(t, err)

	json2, err := serializer.ProtobufToJSON(laptop2)
	require.NoError(t, err)

	require.Equal(t, json1, json2)
}
