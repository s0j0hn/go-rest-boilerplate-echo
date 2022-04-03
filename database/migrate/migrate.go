package migrate

import (
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/database"
	tenantModel "gitlab.com/s0j0hn/go-rest-boilerplate-echo/database/models/tenant"
)

// RunMigrateDatabase is used to prepare for the database
func RunMigrateDatabase() {
	databaseClient := database.Connect()
	err := databaseClient.AutoMigrate(&tenantModel.ModelTenant{})
	if err != nil {
		panic(err)
	}

}
