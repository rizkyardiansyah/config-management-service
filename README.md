# config-management-service

Complete instructions for setting up and running your application.

- makefile, main.go, docker
- testing: unit, integration, automation?
- All prerequisites and dependencies
  • Step-by-step instructions for setting up the environment
  • Exact commands to build, start, and test your project (unit and API/integration)
  o API documentation (preferably in OpenAPI/Swagger format)
  o Schema explanation
  o Notes on your design decisions and trade-offs
  • Feel free to include ideas for improvements, additional features, or creative
  solutions beyond the listed requirements.

## Pre-requisite
    - brew install sqlite (to query in terminal)
    - makefile
    - brew install golangci-lint
    - go install github.com/cosmtrek/air@v1.29.0 for hot reload

## How to Run this Service
    1. make db-reset
    2. make run

## TODO
    - gitignore data/config.db? just re-seed anytime starting run service?
    - go version in go mode increased itself to 1.23 please set 1.21++

## Integration Test
- LOGIN SUCCESS
curl -X POST http://localhost:8089/api/v1/login -d '{"username":"admin","password":"admin123"}' -H "Content-Type: application/json"
- LOGIN FAILED
curl -X POST http://localhost:8089/api/v1/login -d '{"username":"admin","password":"wrongpass"}' -H "Content-Type: application/json"
- CREATE CONFIG
  curl -X POST http://localhost:8089/api/v1/configs \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTkxODc0MjEsImlhdCI6MTc1OTE4NjUyMSwicm9sZSI6ImFkbWluIiwic3ViIjoiOTA2YjBmZTctN2E3ZS00ZDg4LTk5ZWEtZjAyNzg0ODI3MTJjIn0.SYIjXB6lc_2UHGir65V49czkRTt1cNn" \
  -d '{
    "name": "BCA_VA_DAILY_TRESHOLD",
    "type": "object",
    "schema": "{\n  \"type\": \"object\",\n  \"properties\": {\n    \"max_limit\": { \"type\": \"integer\" },\n    \"enabled\": { \"type\": \"boolean\" }\n  },\n  \"required\": [\"max_limit\", \"enabled\"]\n}",
    "input": "{\n  \"max_limit\": 100000,\n  \"enabled\": true\n}",
    "version": 1
  }'


