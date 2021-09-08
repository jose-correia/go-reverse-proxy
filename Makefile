# Build the Golang binary
build:
	./scripts/compile.sh

# Run the reverse proxy service
run:
	go run cmd/proxy/main.go

# Execute linter
lint: bin/golangci-lint
	go fmt ./...
	go vet ./...

# Generate Mocks used in unit tests
gen-mocks: bin/moq
	./bin/moq -pkg db_mock -out ./mocks/app/clients/httpclient/client.go ./app/clients/httpclient HttpClient
	./bin/moq -pkg db_mock -out ./mocks/app/handlers/proxy/handler.go ./app/handlers/proxy Handler

# Generate mock command
bin/moq:
	go build -o bin/moq github.com/matryer/moq

# Run HTTP downstream servers used to test the reverse proxy 
run-testservives:
	go run cmd/testservice/main.go --id 1 --port :8000 &
	go run cmd/testservice/main.go --id 2 --port :8001 --cache-control-max-age 10 &
	go run cmd/testservice/main.go --id 3 --port :8002 --cache-control-max-age 60 &

# Execute unit tests
unit-test: 
	go test -timeout=10s -race -benchmem -tags=unit ./...

# Execute jmeter test file
jmeter:
	jmeter -n -t jmeter_test.jmx -Jhost=test.com -Jpath=http://127.0.0.1:8080/proxy/ -Jusers=20 -Jseconds=60 -l jmeter_test_results.jtl
	
# Execute load test
load-test: run run-testservives jmeter

# Build docker image
docker-build:
	docker build -t reverse-proxy:latest .

# Run Docker container
docker-run:
	docker run reverse-proxy

# Install the Helm Chart in the current Kubernetes context
helm-install:
	helm install reverse-proxy helm-charts/reverse-proxy --values helm-charts/reverse-proxy/values.yaml

# Upgrade already-installed Helm Chart
helm-upgrade:
	helm upgrade reverse-proxy helm-charts/reverse-proxy --values helm-charts/reverse-proxy/values.yaml

