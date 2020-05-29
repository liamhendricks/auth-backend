package cmd

import (
	"fmt"

	"github.com/68696c6c/goat"
	"github.com/spf13/cobra"
)

func init() {
	RootCommand.AddCommand(makeMigrationCommand)
}

var makeMigrationCommand = &cobra.Command{
	Use:   "make:migration name",
	Short: "Add a SQL migration using the provided file name.",
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		db, err := goat.GetMigrationDB()
		if err != nil {
			goat.ExitError(err)
		}

		schema, err := goat.GetSchema(db)
		if err != nil {
			goat.ExitError(err)
		}

		file, err := schema.CreateMigration(args[0])
		if err != nil {
			goat.ExitError(err)
		}

		fmt.Printf("created migration: %s\n", file)
		goat.ExitSuccess()
	},
}
