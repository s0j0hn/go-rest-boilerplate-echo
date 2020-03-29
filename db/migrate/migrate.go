package main

import (
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/db"
	tenantModel "gitlab.com/s0j0hn/go-rest-boilerplate-echo/db/models/tenant"
)

func main() {
	databaseClient := db.DatabaseConnect()
	databaseClient.AutoMigrate(&tenantModel.TenantModel{})
}
