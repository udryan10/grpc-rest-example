package server

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/udryan10/grpc-rest-example/generated"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// implements the MapsService
type mapServer struct{}

// NewMapServer - returns a mapServer
func NewMapServer() generated.MapsServiceServer {
	return new(mapServer)
}

func (m *mapServer) GetMaps(context.Context, *generated.EmptyGet) (*generated.Maps, error) {

	filePath, err := filepath.Abs("server/example.json")

	if err != nil {
		log.Fatalln("Error loading file:", err)
		return nil, grpc.Errorf(codes.Unknown, " error loading file ")
	}

	in, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalln("Error reading file:", err)
		return nil, grpc.Errorf(codes.Unknown, " error reading file from disk ")
	}

	// unmarshal json into protobuff
	mapProto := &generated.Maps{}
	if err := jsonpb.UnmarshalString(string(in), mapProto); err != nil {
		return nil, grpc.Errorf(codes.Unknown, " failed to parse json into protobuff ")
	}

	return mapProto, nil
}