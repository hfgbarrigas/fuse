version=latest

all: build fmt lint test release

fmt:
	go fmt ./...

build: test
		go build

test:
		go test

lint:
		golangci-lint run --fast --exclude-use-default=false

release:
		./scripts/release.sh $(version)