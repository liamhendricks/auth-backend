package models

import (
	"github.com/68696c6c/goat"
	"github.com/icrowley/fake"
)

type Lesson struct {
	goat.Model
	Name       string  `json:"name" binding:"required"`
	LessonData string  `json:"lesson_data"`
	CourseID   goat.ID `json:"course_id"`
}

func MakeLesson() Lesson {
	return Lesson{
		Name:       fake.Title(),
		LessonData: `{"foo":"bar"}`,
	}
}
