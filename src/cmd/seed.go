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

		var firstLessons []*models.Lesson
		var secondLessons []*models.Lesson
		var users []*models.User

		lesson := models.Lesson{
			Name:     "How to Draw Pretty Good",
			Ordering: 0,
			LessonData: &models.LessonData{
				MainHeader:      "Module 1: How to Draw Pretty Good",
				MainDescription: "Description1",
			},
		}

		lesson2 := models.Lesson{
			Name:     "How to Draw Extremely Well",
			Ordering: 1,
			LessonData: &models.LessonData{
				MainHeader:      "Module 2: How to Draw Extremely Well",
				MainDescription: "Description1",
			},
		}

		lesson3 := models.Lesson{
			Name:     "How to Paint Pretty Good",
			Ordering: 0,
			LessonData: &models.LessonData{
				MainHeader:      "Module 1: How to Paint Pretty Good",
				MainDescription: "Description1",
			},
		}

		lesson4 := models.Lesson{
			Name:     "How to Paint Extremely Well",
			Ordering: 1,
			LessonData: &models.LessonData{
				MainHeader:      "Module 2: How to Paint Extremely Well",
				MainDescription: "Description1",
			},
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

		firstLessons = append(firstLessons, &lesson)
		firstLessons = append(firstLessons, &lesson2)
		secondLessons = append(secondLessons, &lesson3)
		secondLessons = append(secondLessons, &lesson4)
		users = append(users, &user)

		course := models.Course{
			Name:       "Introduction to Drawing",
			Lessons:    firstLessons,
			Users:      users,
			CourseType: models.PaidCourse,
		}

		course2 := models.Course{
			Name:       "Level Two Painting Class",
			Lessons:    secondLessons,
			Users:      users,
			CourseType: models.PaidCourse,
		}

		errs := app.CourseRepo.Save(&course)
		if len(errs) > 0 {
			panic(errs)
		}

		errs = app.CourseRepo.Save(&course2)
		if len(errs) > 0 {
			panic(errs)
		}
	},
}
