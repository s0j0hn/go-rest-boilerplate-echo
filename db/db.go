package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/tkc/go-echo-server-sandbox/config"
)

func GetClient() *gorm.DB {

	connection := config.GetDataBaseAccess()
	databaseClient, err := gorm.Open(
		"postgres",
		connection,
	)
	if err != nil {
		panic("failed to connect database")
	}

	databaseClient.DB().SetMaxIdleConns(10)
	databaseClient.DB().SetMaxOpenConns(20)
	return databaseClient
}
