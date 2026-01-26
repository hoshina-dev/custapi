.PHONY: install run test deps lint generate swagger format build

install:
	go mod download

run:
	go run cmd/main.go

test:
	go test ./...

deps:
	go mod tidy

lint:
	golangci-lint run ./...

generate:
	go generate ./...

swagger:
	go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/main.go -o docs --parseDependency --parseInternal

format:
	go run github.com/swaggo/swag/cmd/swag@v1.16.6 fmt
	go fmt ./...
	gofmt -s -w .

build:
	go build -o bin/custapi cmd/main.go

.DEFAULT_GOAL = run
