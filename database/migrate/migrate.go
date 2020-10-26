package migrate

import (
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/database"
	tenantModel "gitlab.com/s0j0hn/go-rest-boilerplate-echo/database/models/tenant"
)

func RunMigrateDatabase() {
	databaseClient := database.Connect()
	err := databaseClient.AutoMigrate(&tenantModel.TenantModel{})
	if err != nil {
		panic(err)
	}

}
