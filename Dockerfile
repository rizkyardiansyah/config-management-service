# ---------- Build Stage ----------
FROM golang:1.21-bullseye AS builder
WORKDIR /src

RUN apt-get update && apt-get install -y gcc libc6-dev
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o configsvc ./cmd/server
RUN go build -o migrate ./cmd/migrate

# ---------- Runtime Stage ----------
FROM debian:bullseye-slim
WORKDIR /app

RUN apt-get update && apt-get install -y sqlite3 libsqlite3-0 && rm -rf /var/lib/apt/lists/*
COPY --from=builder /src/configsvc .
COPY --from=builder /src/migrate .
COPY --from=builder /src/migrations ./migrations
COPY --from=builder /src/data ./data
COPY --from=builder /src/config ./config

EXPOSE 8089
CMD ["./configsvc"]
# ---------- Build Stage ----------
FROM golang:1.21-bullseye AS builder
WORKDIR /src

RUN apt-get update && apt-get install -y gcc libc6-dev
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o configsvc ./cmd/server
RUN go build -o migrate ./cmd/migrate

# ---------- Runtime Stage ----------
FROM debian:bullseye-slim
WORKDIR /app

RUN apt-get update && apt-get install -y sqlite3 libsqlite3-0 && rm -rf /var/lib/apt/lists/*
COPY --from=builder /src/configsvc .
COPY --from=builder /src/migrate .
COPY --from=builder /src/migrations ./migrations
COPY --from=builder /src/data ./data
COPY --from=builder /src/config ./config

EXPOSE 8089
CMD ["./configsvc"]