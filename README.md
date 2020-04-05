# Go REST API Boilerplate

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