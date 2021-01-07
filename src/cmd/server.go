package cmd

import (
	"github.com/68696c6c/goat"
	"github.com/liamhendricks/auth-backend/src/app"
	"github.com/liamhendricks/auth-backend/src/http/routes"
	"github.com/spf13/cobra"
)

func init() {
	RootCommand.AddCommand(serverCommand)
}

var serverCommand = &cobra.Command{
	Use:   "server",
	Short: "Runs the API web server.",
	Run: func(*cobra.Command, []string) {
		config := app.GetConfig()
		app, err := app.GetApp(config)
		if err != nil {
			goat.ExitError(err)
		}

		r := goat.GetRouter()
		routes.InitRoutes(r, app)

		// start server
		err = r.Run()
		if err != nil {
			panic(err)
		}
	},
}
