all:
	protoc -I/usr/local/include -I. \
		-I${GOPATH}/src \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--swagger_out=logtostderr=true:./generated/ \
		--grpc-gateway_out=logtostderr=true:./generated/ \
		--go_out=plugins=grpc:./generated/ \
		./example.proto
clean:
	rm -f example.swagger.json generated/*