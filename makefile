GO_BUILD_ENV=GO111MODULE=on GOFLAGS=-mod=vendor 


build:
	${GO_BUILD_ENV} go build

vendor:
	go mod tidy; \
	go mod download; \
	go mod vendor;