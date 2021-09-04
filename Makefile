all: gen-mocks lint test build

test: 
	go test -timeout=10s -race -benchmem -tags=unit ./...

build:
	./scripts/compile.sh

run:
	go run cmd/proxy/main.go

lint: bin/golangci-lint
	go fmt ./...
	go vet ./...

gen-mocks: bin/moq

bin/moq:
	go build -o bin/moq github.com/matryer/moq

docker-build:
	docker build -t reverse-proxy:latest .

docker-run:
	docker run reverse-proxy

helm-install:
	helm install reverse-proxy helm-charts/reverse-proxy --values helm-charts/reverse-proxy/values.yaml

helm-upgrade:
	helm upgrade reverse-proxy helm-charts/reverse-proxy --values helm-charts/reverse-proxy/values.yaml
