package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	ytTut "github.com/thekure/youtubeTutorialGRPC/grpc"
	"github.com/thekure/youtubeTutorialGRPC/sample"
	"github.com/thekure/youtubeTutorialGRPC/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestServerCreateLaptop(t *testing.T) {
	t.Parallel()

	laptopNoID := sample.NewLaptop()
	laptopNoID.Id = ""

	laptopInvalidID := sample.NewLaptop()
	laptopInvalidID.Id = "invalid-uuid"

	laptopDuplicateID := sample.NewLaptop()
	storeDuplicateID := service.NewInMemoryLaptopStore()
	err := storeDuplicateID.Save(laptopDuplicateID)
	require.Nil(t, err)

	testCases := []struct {
		name   string
		laptop *ytTut.Laptop
		store  service.LaptopStore
		code   codes.Code
	}{
		{
			name:   "success_with_id",
			laptop: sample.NewLaptop(),
			store:  service.NewInMemoryLaptopStore(),
			code:   codes.OK,
		},
		{
			name:   "success_no_id",
			laptop: laptopNoID,
			store:  service.NewInMemoryLaptopStore(),
			code:   codes.OK,
		},
		{
			name:   "failure_invalid_id",
			laptop: laptopInvalidID,
			store:  service.NewInMemoryLaptopStore(),
			code:   codes.InvalidArgument,
		},
		{
			name:   "failure_duplicate_id",
			laptop: laptopDuplicateID,
			store:  storeDuplicateID,
			code:   codes.AlreadyExists,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			request := &ytTut.CreateLaptopRequest{
				Laptop: tc.laptop,
			}

			server := service.NewLaptopServer(tc.store)
			response, err := server.CreateLaptop(context.Background(), request)
			if tc.code == codes.OK {
				require.NoError(t, err)
				require.NotNil(t, response)
				require.NotEmpty(t, response.Id)
				if len(tc.laptop.Id) > 0 {
					require.Equal(t, tc.laptop.Id, response.Id)
				}
			} else {
				require.Error(t, err)
				require.Nil(t, response)
				status, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tc.code, status.Code())
			}
		})
	}
}
