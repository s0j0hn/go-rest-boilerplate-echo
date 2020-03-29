# Go REST API Boilerplate

Working on echo framework
- github.com/labstack/echo
- github.com/jinzhu/gorm


## Requirements

### Go

```
brew install goenv
goenv install 1.14.x
goenv global 1.14.x
goenv rehash
```

## Install proejct go modules

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

## Serve and Hot Reload

``` sh
$ fresh
```