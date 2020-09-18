GO_BUILD_ENV=GO111MODULE=on GOFLAGS=-mod=vendor 
GRPC_OUT=pkg/grpc-server
GRPC_IN=pkg/grpc-server

all: build test vendor

build:
	${GO_BUILD_ENV} go build

test:
	go test -v -count=1

vendor:
	go mod tidy; \
	go mod download; \
	go mod vendor;

proto:
	protoc -I ${GRPC_IN} ${GRPC_IN}/grpc-server.proto \
	 --go_out=plugins=grpc:${GRPC_OUT} \
	 --go_opt=paths=source_relative