SHELL:=/bin/bash
GO_BUILD_ENV=GO111MODULE=on GOFLAGS=-mod=vendor 
PROTO_IN=proto/crypto/v1
PROTO_OUT=proto/crypto/v1
THIRD_PARTY=third_party

CRYPTO=pkg/crypto/v1/
TEST=test_client/v1
SWAGGER=swagger

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
.PHONY: vendor proto

all: vendor build test

build: ## build image.
	${GO_BUILD_ENV} go build

test: ## run test.
	go test -v -count=1 --race

vendor: ## get all dependencies.
	go mod tidy
	go mod download
	go mod vendor

proto: ## generate server and client stubs
	@echo "${BLUE}|-clean up old files from ${CRYPTO} ${TEST} ${SWAGGER}${RESET}"
	@find ${CRYPTO} -name "*.pb.*" | xargs rm -v
	@find ${TEST} -name "*.pb.*" | xargs rm -v
	@find ${SWAGGER} -name "*.swagger.json" | xargs rm -v
	@echo "${GREEN}DONE${RESET}"

	@echo "${BLUE}|-generate new files to ${PROTO_OUT}${RESET}"
	protoc \
	 -I ${THIRD_PARTY} \
	 -I ${GOPATH}/src \
	 -I ./vendor \
	 --proto_path ${PROTO_IN} \
	 --go_out ${PROTO_OUT} --go_opt=plugins=grpc \
	 --grpc-gateway_out ${PROTO_OUT} --grpc-gateway_opt logtostderr=true \
	 --openapiv2_out ${PROTO_OUT} \
      ${PROTO_IN}/*.proto
	@echo "${GREEN}DONE${RESET}"

	@echo "${BLUE}|-copy file to ${CRYPTO}${RESET}"
	@cp -v ${PROTO_OUT}/*.go ${CRYPTO}
	@echo "${GREEN}DONE${RESET}"

	@echo "${BLUE}|-copy file to ${TEST}${RESET}"
	@cp -v ${PROTO_OUT}/*.go ${TEST}
	@echo "${GREEN}DONE${RESET}"

	@echo "${BLUE}|-copy file to ${SWAGGER}${RESET}"
	@cp -v ${PROTO_OUT}/*.swagger.json ${SWAGGER}
	@echo "${GREEN}DONE${RESET}"
	
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

