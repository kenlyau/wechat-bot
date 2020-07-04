package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DllServer string
	Port      string
}

var config Config

func GetConfig() Config {
	return config
}

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	viper.Unmarshal(&config)
}
