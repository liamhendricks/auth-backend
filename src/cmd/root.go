package cmd

import (
	"github.com/68696c6c/goat"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RootCommand = &cobra.Command{
	Use:   "api",
	Short: "Root command for api.",
}

func init() {
	cobra.OnInitialize(goat.Init)
	viper.AutomaticEnv()
}
