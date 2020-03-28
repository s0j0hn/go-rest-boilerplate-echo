package main

import (
	"github.com/tkc/go-echo-server-sandbox/db"
	"github.com/tkc/go-echo-server-sandbox/db/models/tenant"
)

func main() {
	databaseClient := db.GetClient()
	databaseClient.AutoMigrate(&userModel.TenantModel{})
}
