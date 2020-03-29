#! /usr/bin/make

export GO111MODULE=on

all:test build
	@echo DONE!

dep:
	@echo DOWONLOADING MODULES...
	@go mod download

build:
	@echo GENERATING CODE...
	@go build -o dist/server

db:
	@echo GENERATING CODE...
	@docker stack deploy --compose-file docker-compose.yml postgres

test:
	@echo TESTING...
	@go test ./... -v

serve: build
	@echo SERVING...
	@./dist/server