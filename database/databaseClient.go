package database

import (
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/config"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"sync"
)

var once sync.Once

var Client *gorm.DB

func Connect() *gorm.DB {
	once.Do(func() {
		connection := config.GetDataBaseAccess()
		client, err := gorm.Open(postgres.New(postgres.Config{ DSN: connection }), &gorm.Config{})

		if err != nil {
			log.Fatal("Error GORM connect:", err)
		}

		Client = client
	})

	return Client
}

func ConnectForTests() *gorm.DB {
	once.Do(func() {
		client, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})

		if err != nil {
			log.Fatal("Error GORM connect:", err)
		}

		Client = client
	})

	return Client
}
