COVERAGE_OUT := coverage/coverage.out

all: lint test build clean

build: clean_build
	go build -o build/file_transfer_client cmd/main.go

clean: clean_build clean_cache clean_test_cache

clean_build:
	rm -rf build

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