package app

import (
	"github.com/spf13/viper"
)

type Config struct {
	Env string
}

func GetConfig() Config {
	return Config{
		Env: viper.GetString("env"),
	}
}
