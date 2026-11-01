.PHONY: run build tidy test

ADDR ?= :38471

run:
	HTTP_ADDR=$(ADDR) go run ./cmd/api

build:
	go build -o bin/api ./cmd/api

tidy:
	go mod tidy

test:
	go test ./...
