package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"path/filepath"

	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"
	"github.com/udryan10/grpc-rest-example/generated"
	"github.com/udryan10/grpc-rest-example/server"
)

const (
	rpcPort  = 9090
	httpPort = 8080
)

// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "maps-server",
	Short: "Server of gRPC and gRPC-gateway on an example maps api",
	Run: func(cmd *cobra.Command, args []string) {
		// if we are running both, need to backround one so its doesn't block the starting of the other
		if startRpcServer && startHttpServer {
			fmt.Println("starting rpc server...")
			go rpcServer()
			fmt.Println("starting http server...")
			httpServer()
		} else if startRpcServer {
			fmt.Println("starting rpc server...")
			rpcServer()
		} else if startHttpServer {
			fmt.Println("starting http server...")
			httpServer()
		}
	},
}

var startRpcServer bool
var startHttpServer bool

func init() {
	RootCmd.PersistentFlags().BoolVar(&startRpcServer, "startRpcServer", true, "start rpc server")
	RootCmd.PersistentFlags().BoolVar(&startHttpServer, "startHttpServer", true, "start http server")
}

func httpServer() {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	r := http.NewServeMux()
	// the nasty arguments are to tell the json parser to include fields that have default values. By default it removes them
	generatedMux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true, EmitDefaults: false}))
	opts := []grpc.DialOption{grpc.WithInsecure()}

	// Register the generated service handler endpoints
	generated.RegisterMarkersServiceHandlerFromEndpoint(ctx, generatedMux, fmt.Sprintf("localhost:%v", rpcPort), opts)
	r.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		filePath, err := filepath.Abs("./generated/example.swagger.json")

		if err != nil {
			log.Fatalln("Error loading file:", err)

		}

		in, err := ioutil.ReadFile(filePath)
		w.Write(in)
		w.WriteHeader(http.StatusOK)

	})

	r.HandleFunc("/schema", func(w http.ResponseWriter, r *http.Request) {
		conn, err := grpc.Dial("127.0.0.1:9090", grpc.WithInsecure())
		if err != nil {
			panic("unable to connect to grpc server" + err.Error())
		}
		defer conn.Close()

		// setup client
		client := generated.NewMarkersServiceClient(conn)

		body, _ := ioutil.ReadAll(r.Body)

		var input generated.GraphQLQuery
		json.Unmarshal(body, &input)

		schema, err := client.GetMarkersGraphQLSchema(context.Background(), &input)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println("here?")
		fmt.Println(schema.Query)
		w.Write([]byte(schema.Query))
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

	generated.RegisterMarkersServiceServer(s, server.NewMarkersServer())

	fmt.Printf("rpc server listening on :%v \n", rpcPort)
	s.Serve(tcpConn)
}
