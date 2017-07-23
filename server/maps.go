package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/graphql-go/graphql"
	"github.com/udryan10/grpc-rest-example/generated"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// implements the MapsService
type markersServer struct{}

// NewMapServer - returns a mapServer
func NewMarkersServer() generated.MarkersServiceServer {
	return new(markersServer)
}

func (m *markersServer) GetMarkers(context.Context, *generated.EmptyGet) (*generated.Markers, error) {

	return loadMapFromDisk()
}

// implements MapsGetter
type graphQLClient struct{}

func (g graphQLClient) GetMarkers() *generated.Markers {
	maps, _ := loadMapFromDisk()
	return maps
}

func (m *markersServer) GetMarkersGraphQL(c context.Context, g *generated.GraphQLQuery) (*generated.GraphQlMarkersWrapper, error) {

	schemaConfig := graphql.SchemaConfig{Query: GraphQLMapsType}

	schema, err := graphql.NewSchema(schemaConfig)

	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	params := graphql.Params{
		Schema:        schema,
		RequestString: g.Query,
	}

	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		return &generated.GraphQlMarkersWrapper{}, r.Errors[0]
	}
	rJSON, _ := json.Marshal(r)

	fmt.Println(string(rJSON))
	// unmarshal json into protobuff
	markerseProto := &generated.GraphQlMarkersWrapper{}
	customJSONMarshaler := jsonpb.Unmarshaler{}

	if err := customJSONMarshaler.Unmarshal(bytes.NewBuffer(rJSON), markerseProto); err != nil {
		return nil, grpc.Errorf(codes.Unknown, fmt.Sprintf("failed to parse json into protobuff: %v", err))
	}
	fmt.Println(markerseProto)

	return markerseProto, nil
}

func (m *markersServer) GetMarkersGraphQLSchema(c context.Context, g *generated.GraphQLQuery) (*generated.GraphQLQuery, error) {
	fmt.Println("here")
	schemaConfig := graphql.SchemaConfig{Query: GraphQLMapsType}

	schema, err := graphql.NewSchema(schemaConfig)

	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	fmt.Println(g.Query)
	params := graphql.Params{
		Schema:        schema,
		RequestString: g.Query,
	}

	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		fmt.Println(r.Errors)
		return &generated.GraphQLQuery{}, r.Errors[0]
	}

	rJSON, _ := json.Marshal(r.Data)

	fmt.Println(string(rJSON))
	return &generated.GraphQLQuery{
		Query: string(rJSON),
	}, nil
}

func loadMapFromDisk() (*generated.Markers, error) {
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
	mapProto := &generated.Markers{}

	customJSONUnmarshaler := jsonpb.Unmarshaler{
		AllowUnknownFields: true,
	}

	if err := customJSONUnmarshaler.Unmarshal(bytes.NewBuffer(in), mapProto); err != nil {
		return nil, grpc.Errorf(codes.Unknown, fmt.Sprintf("failed to parse json into protobuff: %v", err))
	}

	return mapProto, nil
}
