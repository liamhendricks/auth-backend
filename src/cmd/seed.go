package cmd

import (
	"github.com/68696c6c/goat"
	"github.com/liamhendricks/auth-backend/src/app"
	"github.com/liamhendricks/auth-backend/src/models"
	"github.com/spf13/cobra"
)

func init() {
	RootCommand.AddCommand(seedCommand)
}

var seedCommand = &cobra.Command{
	Use:   "seed",
	Short: "seeds test data",
	Run: func(cmd *cobra.Command, args []string) {
		config := app.GetConfig()
		app, err := app.GetApp(config)
		if err != nil {
			goat.ExitError(err)
		}

		p, err := app.PasswordService.Hash([]byte("password"))
		if err != nil {
			panic(err)
		}

		user := models.User{
			Name:     "Liam",
			Email:    "admin@lanagloschat.com",
			Password: string(p),
			UserType: "Admin",
			Session:  nil,
		}

		course := models.Course{
			Name:       "TestCourse",
			CourseType: "Free",
		}

		errs := app.CourseRepo.Save(&course)
		if len(errs) > 0 {
			goat.ExitError(err)
		}
		println("ID: ", course.ID.String())

		errs = app.UserRepo.Save(&user)
		if len(errs) > 0 {
			goat.ExitError(err)
		}

	},
}
