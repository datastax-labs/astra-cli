all: lint test build

build:
	go build -o bin/astra-cli .

lint:
	go fmt ./...
	golangci-lint run

test:
	go test -v ./...

