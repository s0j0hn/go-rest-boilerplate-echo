package main

import (
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/database"
	tenantModel "gitlab.com/s0j0hn/go-rest-boilerplate-echo/database/models/tenant"
)

func main() {
	databaseClient := database.Connect()
	databaseClient.AutoMigrate(&tenantModel.TenantModel{})
}
