.PHONY: test build run goreleaser-check lint check

test:
	go test -race ./...

build:
	go build ./...

run:
	go run ./cmd/demi

goreleaser-check:
	goreleaser release --snapshot --skip=publish,announce,validate

lint:
	golangci-lint run ./...

check: test lint
