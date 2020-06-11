package cmd

import (
	"github.com/68696c6c/goat"
	sys "github.com/68696c6c/goat/src/sys"
	"github.com/68696c6c/goose"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var dryRun *bool

func init() {
	RootCommand.AddCommand(migrateCommand)
	dryRun = migrateCommand.Flags().BoolP("dry-run", "d", false, "Only report what would have been done.")
}

var migrateCommand = &cobra.Command{
	Use:   "migrate [action] [--dry-run]",
	Short: "Runs the SQL migrations.",
	Run: func(cmd *cobra.Command, args []string) {
		connection, err := goat.GetMigrationDB()
		if err != nil {
			goat.ExitError(err)
		}

		schema, err := goat.GetSchema(connection)
		if err != nil {
			goat.ExitError(err)
		}

		// Inform goose of the current environment.
		env := viper.GetString("env")
		if env == sys.EnvironmentLocal.String() || env == sys.EnvironmentTest.String() {
			goose.SetEnvProduction(false)
		} else {
			goose.SetEnvProduction(true)
		}

		// Only allow 'up' and 'install' operations on production.
		allowed := []string{goose.MigrateOperationInstall, goose.MigrateOperationUp}
		goose.SetProductionOperations(allowed)

		// Perform the migration operation.
		_, _, err = goose.HandleMigrate(schema, args, dryRun)
		if err != nil {
			println(err.Error())
			goat.ExitError(err)
		}

		goat.ExitSuccess()
	},
}
