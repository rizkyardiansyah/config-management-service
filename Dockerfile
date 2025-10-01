# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /src

# Install build deps for go-sqlite3 (CGO)
RUN apk add --no-cache build-base sqlite-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build with CGO enabled
RUN CGO_ENABLED=1 GOOS=linux go build -o /server ./cmd/server

# Runtime stage
FROM alpine:3.20

WORKDIR /app

# runtime deps: sqlite + libc
RUN apk add --no-cache sqlite sqlite-libs

COPY --from=builder /server /server
COPY config/config.json ./config/config.json

RUN mkdir -p /data

ENV APP_PORT=8089
ENV DB_DRIVER=sqlite
ENV DB_DSN=/data/app.db
ENV JWT_SECRET=dev-secret

EXPOSE 8089

CMD ["/server"]