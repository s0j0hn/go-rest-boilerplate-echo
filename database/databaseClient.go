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

// Client is gorm database client.
var Client *gorm.DB

// Connect is used to create the database client.
func Connect() *gorm.DB {
	once.Do(func() {
		connection := config.GetDatabaseAccess()
		client, err := gorm.Open(postgres.New(postgres.Config{ DSN: connection }), &gorm.Config{})

		if err != nil {
			log.Fatal("Error GORM connect:", err)
		}

		Client = client
	})

	return Client
}

// ConnectForTests is used to create mock database client in memory.
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
