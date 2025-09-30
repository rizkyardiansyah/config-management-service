APP_NAME=configsvc
DB_URL=sqlite3://./data/config.db
MIGRATE_CMD=migrate
MIGRATIONS_DIR=./migrations

.PHONY: all build run coverage tidy db-migrate db-migrate-seed db-reset sqlite-shell test lint
all: build

build:
	go build -o bin/$(APP_NAME) ./cmd/server

run:
	CONFIG_PATH=config/config.json air -c .air.toml

coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	open coverage.html

tidy:
	go mod tidy

# schema only
db-migrate:
	go run cmd/migrate/main.go

# schema + insert users
db-migrate-seed:
	go run cmd/migrate/main.go --seed

# nuke db + fresh schema
db-reset:
	go run cmd/migrate/main.go --reset

# SQL-like shell in terminal
sqlite-shell:
	sqlite3 ./data/config.db

test:
	go test -v ./internal/...

lint:
	golangci-lint run ./...