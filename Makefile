APP_NAME=configsvc

.PHONY: all build run test lint clean

all: build

build:
	go build -o bin/$(APP_NAME) ./cmd/server

run:
	go run ./cmd/server

coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Opening coverage.html in your browser to view the report.."
	open coverage.html

tidy:
	go mod tidy