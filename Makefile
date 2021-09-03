all: gen-mocks lint test build

test: unit-test integration-test
	go mod tidy

unit-test:
	go test -timeout=10s -race -benchmem -tags=unit ./...

integration-test:
	go test -timeout=10s -race -benchmem -tags=integration ./...

build:
	./scripts/compile.sh

lint: bin/golangci-lint
	go fmt ./...
	go vet ./...

gen-mocks: bin/moq

bin/moq:
	go build -o bin/moq github.com/matryer/moq
