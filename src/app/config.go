package app

import (
	"github.com/liamhendricks/auth-backend/src/services"
	"github.com/spf13/viper"
)

type Config struct {
	Env            string
	PasswordConfig services.PasswordConfig
}

func GetConfig() Config {
	return Config{
		Env: viper.GetString("env"),
		PasswordConfig: services.PasswordConfig{
			PassPhrase: viper.GetString("pphrase"),
		},
	}
}
