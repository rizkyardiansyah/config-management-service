APP_NAME=configsvc
DB_URL=sqlite3://./data/config.db
MIGRATE_CMD=migrate
MIGRATIONS_DIR=./migrations

.PHONY: all build run test lint clean

all: build

build:
	go build -o bin/$(APP_NAME) ./cmd/server

run:
	go run ./cmd/server

coverage:
	@pkgs="$$(go list ./internal/... | grep -v -E 'internal/(config|secrets|models|migrations)')"; \
	go test -cover -coverprofile=coverage.out $$pkgs
	go tool cover -html=coverage.out -o coverage.html
	@echo "Opening coverage.html in your browser to view the report.."
	open coverage.html

tidy:
	go mod tidy

# schema only
db-migrate:
	go run cmd/migrate/main.go

# schema + insert users
db-migrate-seed:
	go run cmd/migrate/main.go --seed

# nuke db + fresh schema + seeds
db-reset:
	go run cmd/migrate/main.go --reset --seed