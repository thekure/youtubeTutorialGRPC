package serializer

import (
	"github.com/golang/protobuf/jsonpb"
	// "google.golang.org/protobuf/proto"
	"github.com/golang/protobuf/proto"
)

/*
 This method required that I changed from

 	google.golang.org/protobuf/proto

 to

	github.com/golang/protobuf/proto
*/

func ProtobufToJSON(message proto.Message) (string, error) {
	marshaler := jsonpb.Marshaler{
		EnumsAsInts:  false,
		EmitDefaults: true,
		Indent:       "	",
		OrigName:     true,
	}

	return marshaler.MarshalToString(message)
}
