package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/config"
	"log"
	"sync"
)

var once sync.Once

var Client *gorm.DB

func Connect() *gorm.DB {
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

func ConnectForTests() *gorm.DB {
	once.Do(func() {
		client, err := gorm.Open("sqlite3", ":memory:")

		if err != nil {
			log.Fatal("Error GORM connect:", err)
		}

		client.DB().SetMaxIdleConns(10)
		client.DB().SetMaxOpenConns(20)
		Client = client
	})

	return Client
}
