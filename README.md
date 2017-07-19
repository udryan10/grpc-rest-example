# grpc-rest-example


## build swagger, grpc gateway and protobuff
```
protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --swagger_out=logtostderr=true:. \
  --grpc-gateway_out=logtostderr=true:. \
  --go_out=./generated/ \
  ./example.proto
```