COVERAGE_OUT := coverage/coverage.out

all: install gen lint test build clean

install:
	go mod tidy
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go get github.com/golang/mock/mockgen@v1.6.0
	go get github.com/golang/mock/mockgen@v1.6.0
	go install github.com/golang/mock/mockgen
	export PATH=$PATH:$(go env GOPATH)/bin

build: clean_build
	go build -o build/file_transfer_server cmd/main.go

clean: clean_build clean_cache clean_test_cache

clean_build:
	rm -rf build

gen:
	go generate ./...

test: clean_test_cache
	@echo "+ $@"
	go test -race -v -coverprofile=$(COVERAGE_OUT) ./...
	go tool cover -html=$(COVERAGE_OUT)
	go tool cover -func=$(COVERAGE_OUT)

clean_test_cache:
	go clean -testcache
	rm -f $(COVERAGE_OUT)

clean_cache:
	go clean -cache

lint:
	@golangci-lint run --timeout 10m