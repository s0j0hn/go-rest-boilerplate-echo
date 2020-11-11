package config

import (
	"fmt"
	"github.com/spf13/viper"
)

func getViper() *viper.Viper {
	vp := viper.New()
	vp.SetConfigName("config")
	vp.SetConfigType("yaml")
	vp.AddConfigPath(".")
	err := vp.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("error: %s", err))
	}
	return vp
}

// IsProd to get the env for prod or not.
func IsProd() bool {
	return getViper().Get("app") == "prod"
}

// GetAddress is used to get the webserver host to listen on.
func GetAddress() string {
	return getViper().Get("address").(string)
}

// GetDatabaseAccess is used to get the database credentials.
func GetDatabaseAccess() string {
	v := getViper()
	connection := fmt.Sprintf(
		"host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable",
		v.Get("database.host"),
		v.Get("database.user"),
		v.Get("database.password"),
		v.Get("database.name"),
	)
	return connection
}

// GetRabbitMQAccess is used to get the rabbitmq credentials.
func GetRabbitMQAccess() string {
	v := getViper()
	connection := fmt.Sprintf("amqp://%s:%s@%s/",
		v.Get("rabbitmq.user"),
		v.Get("rabbitmq.password"),
		v.Get("rabbitmq.host"),
	)
	return connection
}

// GetAMQPPushQueue is used to get the env value for queue to push events.
func GetAMQPPushQueue() string {
	v := getViper()
	return fmt.Sprintf("%s", v.Get("rabbitmq.pushqueue"))
}

// GetAMQPQListenQueue is used to get the env value for queue to listen to events.
func GetAMQPQListenQueue() string {
	v := getViper()
	return fmt.Sprintf("%s", v.Get("rabbitmq.listenqueue"))
}