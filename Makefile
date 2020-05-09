#! /usr/bin/make
PROJECT_NAME := "boilerplate"
PKG := "gitlab.com/s0j0hn/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

all:test build
	@echo DONE!

dep:
	@echo DOWONLOADING MODULES...
	@go mod download

build: swagger
	@echo GENERATING CODE...
	@go build -o dist/server

swagger:
	@echo GENERATING SWAGGER...
	@swag init

postgres:
	@echo GENERATING CODE...
	@docker stack deploy --compose-file docker-compose.yml postgres

test:
	@echo UNIT TESTING...
	@go test ./... -v -cover -coverprofile=coverage.cov

msan:
	@echo MEMORY TESTING...
	@go test -msan -short $(go list ./... | grep -v /vendor/)

race:
	@echo RACE TESTING...
	@go test -race -short $(go list ./... | grep -v /vendor/)

coverage: test
	@echo COVERAGE TESTING...
	@go tool cover -func=coverage.cov

serve: swagger build
	@echo SERVING...
	@./dist/server