GO_BUILD_ENV=GO111MODULE=on GOFLAGS=-mod=vendor 
CRYPTO_IN=pkg/crypto/proto/v1
CRYPTO_OUT=pkg/crypto/v1
GOOGLE_API=pkg/crypto/proto

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
	protoc \
	 -I${GOOGLE_API} \
	 --proto_path ${CRYPTO_IN} \
	 --go_out ${CRYPTO_OUT} --go_opt=plugins=grpc \
	 --grpc-gateway_out ${CRYPTO_OUT} --grpc-gateway_opt logtostderr=true \
      crypto.proto
	
tool:
	go install \
    github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger \
    github.com/golang/protobuf/protoc-gen-go
