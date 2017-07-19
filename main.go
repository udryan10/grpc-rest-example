package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/jsonpb"
	pb "github.com/udryan10/grpc-rest-example/generated"
)

func main() {

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {

		// Read the existing address book.
		in, err := ioutil.ReadFile("./example.json")
		if err != nil {
			log.Fatalln("Error reading file:", err)
		}

		// unmarshal json into protobuff
		mapProto := &pb.Maps{}
		if err := jsonpb.UnmarshalString(string(in), mapProto); err != nil {
			log.Fatalln("Failed to parse json:", err)
		}

		// .... do something meaningful ... //
		for _, marker := range mapProto.Markers {
			fmt.Println(marker.AwayTeam)
		}

		// marshal protobuff back to json and return
		pbToJson := jsonpb.Marshaler{}
		returnJSON, err := pbToJson.MarshalToString(mapProto)

		if err != nil {
			c.Abort()
		}
		c.String(200, returnJSON)
	})

	r.POST("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"foo": "bar",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
