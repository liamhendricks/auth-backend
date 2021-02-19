package app

import (
	"github.com/spf13/viper"
)

type Config struct {
	Env                  string
	StripeEndpointSecret string
	StripeSecretKey      string
}

func GetConfig() Config {
	return Config{
		Env:                  viper.GetString("env"),
		StripeEndpointSecret: viper.GetString("stripe_endpoint_secret"),
		StripeSecretKey:      viper.GetString("stripe_secret_key"),
	}
}
