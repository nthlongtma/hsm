GO_BUILD_ENV=GO111MODULE=on GOFLAGS=-mod=vendor 
CRYPTO=pkg/crypto

all: vend build test

build:
	${GO_BUILD_ENV} go build

test:
	go test -v -count=1

vend:
	go mod tidy
	go mod download
	go mod vendor

proto:
	cd ${CRYPTO}; \
	pwd; \
	protoc -I ./proto \
	 --go_out . --go_opt=plugins=grpc  --go_opt=paths=source_relative \
	 --grpc-gateway_out . --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative \
     ./proto/crypto/v1/crypto.proto

tool:
	go install \
    github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger \
    github.com/golang/protobuf/protoc-gen-go
