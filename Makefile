all:
	protoc -I/usr/local/include -I. \
		-I${GOPATH}/src \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--swagger_out=logtostderr=true:./generated/ \
		--grpc-gateway_out=logtostderr=true:./generated/ \
		--gogoopsee_out=plugins=grpc+graphql,Mgoogle/protobuf/descriptor.proto=github.com/gogo/protobuf/protoc-gen-gogo/descriptor,Mgoogle/api/annotations.proto=google.golang.org/genproto/googleapis/api/annotations:./generated/ \
		./example.proto
clean:
	rm -f example.swagger.json generated/*