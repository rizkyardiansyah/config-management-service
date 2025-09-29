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

## How to Run this Service
    1. make db-reset
    2. make run

## TODO
    - gitignore data/config.db? just re-seed anytime starting run service?
    - go version in go mode increased itself to 1.23 please set 1.21++

## Integration Test
- LOGIN SUCCESS
curl -X POST http://localhost:8089/api/login -d '{"username":"admin","password":"admin123"}' -H "Content-Type: application/json"
- LOGIN FAILED
curl -X POST http://localhost:8089/api/login -d '{"username":"admin","password":"wrongpass"}' -H "Content-Type: application/json"
