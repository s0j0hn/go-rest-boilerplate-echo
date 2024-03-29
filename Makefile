#! /usr/bin/make
PROJECT_NAME := "boilerplate"
PKG := "gitlab.com/s0j0hn/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

all:start-services build test serve
	@echo DONE

dep:
	@echo DOWONLOADING MODULES...
	@go mod download
	@echo DONE

build:
	@echo COMPILING...
	@go build -o dist/server
	@echo DONE

swagger:
	@echo GENERATING SWAGGER...
	@swag init
	@echo DONE

start-services:
	@echo STARTING DOCKER SERVICES...
	@docker stack deploy --compose-file docker-compose.yml goboilerplate

stop-services:
	@echo STOPING DOCKER SERVICES...
	@docker stack rm goboilerplate

test:
	@echo UNIT TESTING...
	@go test ./... -v -cover -coverprofile=coverage.cov
	@echo DONE

msan:
	@echo MEMORY TESTING...
	@go test -msan -short $(go list ./... | grep -v /vendor/)

race:
	@echo RACE TESTING...
	@go test -race -short $(go list ./... | grep -v /vendor/)

coverage:
	@echo COVERAGE TESTING...
	@go tool cover -func=coverage.cov

serve: swagger build
	@echo SERVING...
	@./dist/server

update-docker:
	@echo CREATING IMAGE AND PUSHING INTO GITLAB...
	@docker build -t registry.gitlab.com/s0j0hn/go-rest-boilerplate-echo .
	@docker push registry.gitlab.com/s0j0hn/go-rest-boilerplate-echo
