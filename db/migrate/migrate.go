package main

import (
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/db"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/db/models/tenant"
)

func main() {
	databaseClient := db.GetClient()
	databaseClient.AutoMigrate(&userModel.TenantModel{})
}
