#!/usr/bin/env -S just --justfile

# Swagger
swag-fmt:
    swag fmt

swag-gen:
    swag init --generalInfo ./internal/server/server.go --output ./docs --parseInternal

swagger: swag-fmt swag-gen


# Go linter
go-lint:
    /home/dima/golangci-lint/golangci-lint run

test-coverage:
    go test -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html

tests-services:
    go test ./internal/services -count=1

tests-handlers:
    go test ./internal/handlers -count=1

race-server:
    go run -race app_server.go