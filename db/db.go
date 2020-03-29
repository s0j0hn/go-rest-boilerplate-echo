package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/config"
	"log"
	"sync"
)

var once sync.Once

var Client *gorm.DB

func DatabaseConnect() *gorm.DB {
	once.Do(func() {
		connection := config.GetDataBaseAccess()
		client, err := gorm.Open(
			"postgres",
			connection,
		)

		if err != nil {
			log.Fatal("Error GORM connect:", err)
		}

		client.DB().SetMaxIdleConns(10)
		client.DB().SetMaxOpenConns(20)
		Client = client
	})

	return Client
}
