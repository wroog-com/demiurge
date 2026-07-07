.PHONY: test build run

test:
	go test -race ./...

build:
	go build ./...

run:
	go run ./cmd/demi
