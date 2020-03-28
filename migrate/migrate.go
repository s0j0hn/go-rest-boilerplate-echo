package main

import (
	"github.com/tkc/go-echo-server-sandbox/db"
	"github.com/tkc/go-echo-server-sandbox/models/tenant"
)

func main() {
	databaseClient := db.GetClient()
	databaseClient.AutoMigrate(&userModel.TenantModel{})
}
