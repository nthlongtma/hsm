GO_BUILD_ENV=GO111MODULE=on GOFLAGS=-mod=vendor 


all: build test vendor

build:
	${GO_BUILD_ENV} go build

test:
	go test -v -count=1

vendor:
	go mod tidy; \
	go mod download; \
	go mod vendor;