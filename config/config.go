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
	vp.AddConfigPath("../")
	vp.AddConfigPath("../../")
	vp.AddConfigPath("../../../")
	vp.AddConfigPath("../../../../")
	vp.AddConfigPath("../../../../../")
	err := vp.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("error: %s", err))
	}
	return vp
}

func IsProd() bool {
	return getViper().Get("app") == "prod"
}

func GetPort() string {
	return getViper().Get("port").(string)
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
