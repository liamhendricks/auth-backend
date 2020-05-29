package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	RootCommand.AddCommand(migrateCommand)
}

var migrateCommand = &cobra.Command{
	Use:   "migrate aaction",
	Short: "Runs the SQL migrations.",
	Run: func(_ *cobra.Command, args []string) {
		fmt.Println("migrate")
	},
}
