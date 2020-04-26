package migrate

import (
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/database"
	tenantModel "gitlab.com/s0j0hn/go-rest-boilerplate-echo/database/models/tenant"
)

func MigrateDatabase() {
	databaseClient := database.Connect()
	databaseClient.AutoMigrate(&tenantModel.TenantModel{})
}
