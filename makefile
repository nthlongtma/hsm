GO_BUILD_ENV=GO111MODULE=on GOFLAGS=-mod=vendor 
CRYPTO_IN=pkg/crypto/proto/v1
CRYPTO_OUT=pkg/crypto/v1
GOOGLE_API=pkg/crypto/proto

.PHONY: vendor

all: vendor build test

build:
	${GO_BUILD_ENV} go build

test:
	go test -v -count=1

vendor:
	go mod tidy
	go mod download
	go mod vendor

proto:
	@echo "remove old file"
	@find ./pkg/crypto/v1/ -name "*.pb.*" | xargs rm

	@echo "gen files"
	@protoc \
	 -I${GOOGLE_API} \
	 --proto_path ${CRYPTO_IN} \
	 --go_out ${CRYPTO_OUT} --go_opt=plugins=grpc \
	 --grpc-gateway_out ${CRYPTO_OUT} --grpc-gateway_opt logtostderr=true \
      crypto.proto

	@echo copy file to test
	@cp -v ${CRYPTO_OUT}/*.pb.* test_client/v1
	
tool:
	go get \
		github.com/golang/protobuf/protoc-gen-go \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 
