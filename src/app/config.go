package app

import (
	"github.com/spf13/viper"
)

type Config struct {
	Env                  string
	StripeEndpointSecret string
	StripeSecretKey      string
	SendgridSecretKey    string
	SendgridFromEmail    string
	SendgridFromName     string
	SendgridBaseURL      string
	ResetTemplateID      string
	PurchaseTemplateID   string
	SignupTemplateID     string
}

func GetConfig() Config {
	return Config{
		Env:                  viper.GetString("env"),
		StripeEndpointSecret: viper.GetString("stripe_endpoint_secret"),
		StripeSecretKey:      viper.GetString("stripe_secret_key"),
		SendgridSecretKey:    viper.GetString("sendgrid_secret_key"),
		SendgridFromEmail:    viper.GetString("sendgrid_from_email"),
		SendgridFromName:     viper.GetString("sendgrid_from_name"),
		SendgridBaseURL:      viper.GetString("sendgrid_base_url"),
		PurchaseTemplateID:   viper.GetString("purchase_template_id"),
		SignupTemplateID:     viper.GetString("signup_template_id"),
		ResetTemplateID:      viper.GetString("reset_template_id"),
	}
}
