package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"path/filepath"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/udryan10/grpc-rest-example/generated"
	"github.com/udryan10/grpc-rest-example/server"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	rpcPort  = 9090
	httpPort = 8080
)

func main() {

	go rpcServer()
	httpServer()

}

func httpServer() {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	//r := mux.NewRouter()
	r := http.NewServeMux()
	generatedMux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	// Register the generated service handler endpoints
	generated.RegisterMapsServiceHandlerFromEndpoint(ctx, generatedMux, fmt.Sprintf("localhost:%v", rpcPort), opts)

	r.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		filePath, err := filepath.Abs("example.swagger.json")

		if err != nil {
			log.Fatalln("Error loading file:", err)

		}

		in, err := ioutil.ReadFile(filePath)
		w.Write(in)
		w.WriteHeader(http.StatusOK)

	})

	// bind generated mux to our main handler
	r.Handle("/", generatedMux)

	fmt.Printf("http server running on :%v \n", httpPort)

	http.ListenAndServe(fmt.Sprintf(":%v", httpPort), r)
}

func rpcServer() {

	tcpConn, err := net.Listen("tcp", fmt.Sprintf(":%v", rpcPort))
	if err != nil {
		panic("unable to establish tcpConn")
	}
	s := grpc.NewServer()

	generated.RegisterMapsServiceServer(s, server.NewMapServer())

	fmt.Printf("rpc server listening on :%v \n", rpcPort)
	s.Serve(tcpConn)
}
