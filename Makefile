APP_NAME = configsvc
DB_PATH  := $(abspath ./data/config.db)
MIGRATIONS_DIR = ./migrations

.PHONY: all build run coverage tidy \
        db-migrate db-migrate-seed db-reset \
        db-migrate-docker db-migrate-seed-docker db-reset-docker \
        sqlite-shell test lint docker-up docker-down

all: build

build:
	go build -o bin/$(APP_NAME) ./cmd/server

run:
	DB_DSN=$(DB_PATH) CONFIG_PATH=config/config.json air -c .air.toml

coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	open coverage.html

tidy:
	go mod tidy

# -----------------
# LOCAL DB commands
# -----------------

db-migrate:
	DB_DSN=$(DB_PATH) go run cmd/migrate/main.go

db-migrate-seed:
	DB_DSN=$(DB_PATH) go run cmd/migrate/main.go --seed

db-reset:
	DB_DSN=$(DB_PATH) go run cmd/migrate/main.go --reset

# -----------------
# DOCKER DB commands
# -----------------
# assumes multi-stage Dockerfile builds ./migrate binary

db-migrate-docker:
	docker-compose run --rm $(APP_NAME) ./migrate

db-migrate-seed-docker:
	docker-compose run --rm $(APP_NAME) ./migrate --seed

db-reset-docker:
	docker-compose run --rm $(APP_NAME) ./migrate --reset

# -----------------
# Misc
# -----------------

sqlite-shell:
	sqlite3 $(DB_PATH)

test:
	go test -v ./internal/...

lint:
	golangci-lint run ./...

docker-up:
	@echo "Checking if port 8089 is in use..."
	@PID=$$(lsof -ti :8089) ; \
	if [ -n "$$PID" ]; then \
		echo "Port 8089 is in use by PID $$PID. Killing it..."; \
		kill -9 $$PID; \
	fi
	@echo "Stopping old container if running..."
	-docker stop $(APP_NAME) >/dev/null 2>&1 || true
	-docker rm $(APP_NAME) >/dev/null 2>&1 || true
	@echo "Starting fresh container..."
	docker-compose up --build -d

docker-down:
	docker-compose down