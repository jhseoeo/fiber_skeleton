.PHONY: run build test lint swagger tidy

run:
	go run ./cmd/main.go

build:
	go build -o bin/server ./cmd/main.go

test:
	go test ./...

lint:
	golangci-lint run

swagger:
	swag init -g cmd/main.go

tidy:
	go mod tidy
