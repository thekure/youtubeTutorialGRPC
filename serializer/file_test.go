package serializer_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	ytTut "github.com/thekure/youtubeTutorialGRPC/grpc"
	"github.com/thekure/youtubeTutorialGRPC/sample"
	"github.com/thekure/youtubeTutorialGRPC/serializer"
	"google.golang.org/protobuf/proto"
)

func TestFileSerializer(t *testing.T) {
	// Runs test in parallel, thus potentially detects race conditions.
	t.Parallel()

	binaryFile := "../tmp/laptop.bin"
	jsonFile := "../tmp/laptop.json"

	laptop1 := sample.NewLaptop()
	err := serializer.WriteProtobufToBinaryFile(laptop1, binaryFile)

	/*
		The following requires a package you get like this:
			go get github.com/stretchr/testify
	*/
	require.NoError(t, err)

	laptop2 := &ytTut.Laptop{}
	err = serializer.ReadProtobufFromBinaryFile(binaryFile, laptop2)
	require.NoError(t, err)
	require.True(t, proto.Equal(laptop1, laptop2))

	err = serializer.WriteProtobufToJSONFile(laptop1, jsonFile)
	require.NoError(t, err)
}

/*
	go test ./serializer/file_test.go

	runs only that specific file

	go test ./serializer

	runs the tests within the specified folder

	go run test ./...

	runs tests in all subpackages
*/
