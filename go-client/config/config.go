package config

import (
	"github.com/spf13/viper"
)

type StConfig struct {
	DllServer string
	Port      string
}

var Config StConfig

func SetUp() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	viper.Unmarshal(&Config)
}
