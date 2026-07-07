.PHONY: test build run goreleaser-check

test:
	go test -race ./...

build:
	go build ./...

run:
	go run ./cmd/demi

goreleaser-check:
	goreleaser release --snapshot --skip=publish,announce,validate
