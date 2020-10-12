version=latest

all: build lint test release

build: test
		go build

test:
		go test

lint:
		golangci-lint run --fast --exclude-use-default=false

release:
		./scripts/release.sh $(version)