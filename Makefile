.PHONY: build test lint install clean

build:
	go build -o bin/keen ./cmd/keen

test:
	go test ./...

lint:
	golangci-lint run

install:
	go install ./cmd/keen

clean:
	rm -rf bin/

build-mcp:
	go build -o bin/keen-mcp ./cmd/keen-mcp

install-mcp:
	go install ./cmd/keen-mcp

build-all: build build-mcp
