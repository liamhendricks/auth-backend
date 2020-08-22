package models

import (
	"github.com/68696c6c/goat"
	"github.com/icrowley/fake"
)

type Lesson struct {
	goat.Model
	Name     string  `json:"name" binding:"required"`
	CourseID goat.ID `json:"course_id"`
}

func MakeLesson() Lesson {
	return Lesson{
		Name: fake.Title(),
	}
}
