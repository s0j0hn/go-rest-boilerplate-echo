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

func IsProd() bool {
	return getViper().Get("app") == "prod"
}

func GetAddress() string {
	return getViper().Get("address").(string)
}

func GetDataBaseAccess() string {
	v := getViper()
	connection := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable",
		v.Get("database.host"),
		v.Get("database.user"),
		v.Get("database.password"),
		v.Get("database.name"),
	)
	return connection
}

func GetRabbitMQAccess() string {
	v := getViper()
	connection := fmt.Sprintf("amqp://%s:%s@%s/",
		v.Get("rabbitmq.user"),
		v.Get("rabbitmq.password"),
		v.Get("rabbitmq.host"),
	)
	return connection
}
