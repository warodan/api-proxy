APP_NAME=api-proxy

.PHONY: run build test

run:
	go run ./cmd/app

build:
	go build -o bin/$(APP_NAME) ./cmd/app

test:
	go test ./... -v
