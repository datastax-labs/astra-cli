all: lint test build

build:
	go build -o bin/acm ./...
lint:
	go fmt ./...
	golangci-lint run

test:
	go test -v ./...

