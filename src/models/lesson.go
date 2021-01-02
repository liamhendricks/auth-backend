package models

import (
	"github.com/68696c6c/goat"
	"github.com/icrowley/fake"
)

type Lesson struct {
	goat.Model
	Name     string  `json:"name" binding:"required"`
	Ordering int     `json:"ordering" binding:"required"`
	Data     string  `json:"data"`
	CourseID goat.ID `json:"course_id"`
}

func MakeLesson(id goat.ID) Lesson {
	return Lesson{
		Name: fake.Title(),
		Data: `
    {
      "foo": "bar",
    }
    `,
		Ordering: 0,
		CourseID: id,
	}
}
