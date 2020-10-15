SHELL:=/bin/bash
GO_BUILD_ENV=GO111MODULE=on GOFLAGS=-mod=vendor 
CRYPTO_IN=pkg/crypto/proto/v1
CRYPTO_OUT=pkg/crypto/v1
GOOGLE_API=pkg/crypto/proto

# define standard colors
BLACK        := $(shell tput -Txterm setaf 0)
RED          := $(shell tput -Txterm setaf 1)
GREEN        := $(shell tput -Txterm setaf 2)
YELLOW       := $(shell tput -Txterm setaf 3)
LIGHTPURPLE  := $(shell tput -Txterm setaf 4)
PURPLE       := $(shell tput -Txterm setaf 5)
BLUE         := $(shell tput -Txterm setaf 6)
WHITE        := $(shell tput -Txterm setaf 7)

RESET := $(shell tput -Txterm sgr0)

# set target color
TARGET_COLOR := $(BLUE)
POUND = \#

.SILENT:
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
	@echo "${GREEN}|-clean up old files${RESET}"
	@find ./pkg/crypto/v1/ -name "*.pb.*" | xargs rm -v

	@echo "${GREEN}|-generate new files to ${CRYPTO_OUT}${RESET}"
	@protoc \
	 -I ${GOOGLE_API} \
	 -I ${GOPATH}/src \
	 -I ./vendor \
	 --proto_path ${CRYPTO_IN} \
	 --go_out ${CRYPTO_OUT} --go_opt=plugins=grpc \
	 --grpc-gateway_out ${CRYPTO_OUT} --grpc-gateway_opt logtostderr=true \
      crypto.proto

	@echo "${GREEN}|-copy file to test_client/v1${RESET}"
	@cp -v ${CRYPTO_OUT}/*.pb.* test_client/v1
	
tool: ## get all the tools
	${GO_BUILD_ENV} go get \
		github.com/golang/protobuf/protoc-gen-go \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 

help:
	@echo ""
	@echo "    ${WHITE}:: ${RED}Self-documenting Makefile${RESET} ${WHITE}::${RESET}"
	@echo ""
	@echo "Document targets by adding '$(POUND)$(POUND) comment' after the target"
	@echo ""
	@echo "Example:"
	@echo "  | job1:  $(POUND)$(POUND) help for job 1"
	@echo "  | 	@echo \"run stuff for target1\""
	@echo ""
	@echo "${WHITE}-----------------------------------------------------------------${RESET}"
	@grep -E '^[a-zA-Z_0-9%-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "${TARGET_COLOR}%-30s${RESET} %s\n", $$1, $$2}'

