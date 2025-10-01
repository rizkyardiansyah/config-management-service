# Config Service

Configuration management service with versioning, rollback, and JSON schema validation.  
Built with **Go + Gin + GORM (SQLite)**.

---

## ⚙️ Pre-requisites

Before you start, make sure you have:

- **Go** (≥ 1.21)
- **SQLite CLI** (for debugging DB):
  ```bash
  brew install sqlite
  ```
- **GNU Make** (most systems already have it; otherwise install via `brew install make`)
- **Air** for hot reload during development:
  ```bash
  go install github.com/cosmtrek/air@v1.29.0
  ```
- (Optional) **Docker** & **docker-compose** for containerized runs

---

## 🚀 Running the Project

You can run the service **locally** or via **Docker**.

---

### ▶️ Local (Go + Air hot reload)

1. Run migrations:
   ```bash
   make db-migrate-seed
   ```
   This sets up schema + inserts a default admin user.

2. Start the service with hot reload:
   ```bash
   make run
   ```

3. API is available at:
   ```
   http://localhost:8089/api/v1
   
   ```   

---

### 🐳 Docker

1. Build and run:
   ```bash
   make docker-up
   ```
   Please ensure your Docker Software runs. If it is closed, open it again.


2. Run migrations:
   ```bash
   make db-migrate-seed
   ```
   This sets up schema + inserts a default admin user.


3. API is available at:
   ```
   http://localhost:8089/api/v1
   ```

4. Stop:
   ```bash
   make docker-down
   ```

---

## 🧪 Development

- Run all tests:
  ```bash
  make test
  ```

- Run with coverage report:
  ```bash
  make coverage
  ```

- Lint:
  ```bash
  make lint
  ```

---

## 🗃 Database

- The service uses **SQLite**.
- DB file is stored at `./data/config.db` (shared with Docker container).

### Reset database

```bash
make db-reset
```

### Open SQLite shell

```bash
make sqlite-shell
```

---

## 🔑 Authentication

- Auth API issues **JWT tokens** (`/login`).
- Pass token in requests as:
  ```http
  Authorization: Bearer <JWT_TOKEN>
  ```
  - Token TTL is set in (`config/config.json`)
  ```
    {
      "Port": 8089,
      "AccessTokenTTLInDays": 7,
      "RefreshTokenTTLInMinutes": 60
    }
  ```

---

## 📖 API Docs

- Sanity Test Collection 
  ```
  https://.postman.co/workspace/Go-APIs~c6d86e5b-b59d-4337-8493-cc425531230b/collection/20664073-95f151ea-c13f-409b-bc92-a0582cd00198?action=share&creator=20664073
  ```
- Open API Collection 
  ```
  https://.postman.co/workspace/Go-APIs~c6d86e5b-b59d-4337-8493-cc425531230b/collection/20664073-95f151ea-c13f-409b-bc92-a0582cd00198?action=share&creator=20664073
  ```
---

## 📦 Makefile Reference

### Local

- `make build` → build binary into `bin/configsvc`
- `make run` → run server with hot reload (Air)
- `make db-migrate` → run DB migrations (schema only)
- `make db-migrate-seed` → migrations + seed admin user
- `make db-reset` → nuke DB + fresh schema
- `make sqlite-shell` → open SQLite REPL
- `make test` → run unit tests
- `make coverage` → run tests + show coverage report
- `make lint` → run linter

### Docker

- `make docker-up` → build & start service in Docker
- `make docker-down` → stop containers  
