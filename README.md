# Go REST API Boilerplate
[![Build Status](https://gitlab.com/s0j0hn/go-rest-boilerplate-echo/badges/master/build.svg)](https://gitlab.com/s0j0hn/go-rest-boilerplate-echo/commits/master)
[![Coverage Report](https://gitlab.com/s0j0hn/go-rest-boilerplate-echo/badges/master/coverage.svg)](https://gitlab.com/s0j0hn/go-rest-boilerplate-echo/commits/master)
[![Go Report Card](https://goreportcard.com/badge/gitlab.com/s0j0hn/go-rest-boilerplate-echo)](https://goreportcard.com/report/gitlab.com/s0j0hn/go-rest-boilerplate-echo)
[![License MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](https://img.shields.io/badge/License-MIT-brightgreen.svg)

- github.com/labstack/echo 
- github.com/jinzhu/gorm
- github.com/casbin/casbin/v2


## Requirements

### Go

```
brew install goenv
goenv install 1.14.x
goenv global 1.14.x
goenv rehash
```

## Install project go modules

```sh
make dep
```

## Database (Postgres) Config

`before start: cp config.yaml.example config.yaml`

 ... and update your database config

``` yaml
app: local
port: :8080

database:
  name: <database name>
  name: <user name>
  password: <user passWord>
```

## Gorm Database Migration

``` sh
$ go run ./migrate/migrate.go
```

## Launch the server

``` sh
$ make postgres  // Setup local database with docker compose
$ make serve
```