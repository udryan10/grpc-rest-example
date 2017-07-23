package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"google.golang.org/grpc"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/jsonpb"
	"github.com/udryan10/grpc-rest-example/generated"
	context "golang.org/x/net/context"
)

var log = logrus.New()

func main() {

	tests := []func(){
		httpRequestNoClient,
		httpRequestClient,
		httpRequestProtoBuf,
		rpcRequest,
	}
	fmt.Println("beginning tests...\n\n\n")

	// run the tests
	for _, test := range tests {
		fmt.Println("******************************************")
		test()
		fmt.Println("****************************************** \n\n")
	}

}

func httpRequestNoClient() {

	fmt.Println("making http call http://localhost:8080/markers \n")
	body, err := _makeHttp("application/json", "application/json")

	if err != nil {
		panic("error in _makeHttp()" + err.Error())
	}

	fmt.Println("JSON: ")
	fmt.Println("    " + string(body) + "\n")

}

func httpRequestClient() {
	fmt.Println("making http call http://localhost:8080/markers, receiving json and marshalling into protobuff \n")

	body, err := _makeHttp("application/json", "application/json")

	if err != nil {
		panic("error in _makeHttp()" + err.Error())
	}

	// marshall JSON string into protobuff type from generated client
	maps := generated.Markers{}

	jsonpb.UnmarshalString(string(body), &maps)

	fmt.Println("protobuff: ")
	fmt.Println("        " + maps.String() + "\n")
}

func rpcRequest() {
	fmt.Println("making rcp request GetMaps()\n")

	// connect to local grpc sever
	conn, err := grpc.Dial("127.0.0.1:9090", grpc.WithInsecure())
	if err != nil {
		panic("unable to connect to grpc server" + err.Error())
	}
	defer conn.Close()

	// setup client
	client := generated.NewMarkersServiceClient(conn)

	// rpc call
	maps, err := client.GetMarkers(context.Background(), &generated.EmptyGet{})

	if err != nil {
		panic("error in GetMaps()")
	}

	fmt.Println("protobuff: ")
	fmt.Println("        " + maps.String() + "\n")
}

func _makeHttp(contentType, accept string) ([]byte, error) {
	request, _ := http.NewRequest("GET", "http://localhost:8080/markers", nil)

	request.Header.Add("Content-Type", contentType)
	request.Header.Add("Accept", accept)
	client := &http.Client{}
	resp, err := client.Do(request)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}
