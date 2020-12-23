package models

import (
	"github.com/68696c6c/goat"
	"github.com/icrowley/fake"
)

type Lesson struct {
	goat.Model
	Name       string      `json:"name" binding:"required"`
	Ordering   int         `json:"ordering" binding:"required"`
	LessonData *LessonData `json:"lesson_data"`
	CourseID   goat.ID     `json:"course_id"`
}

type LessonData struct {
	goat.Model
	LessonID        goat.ID `json:"lesson_id" binding:"required"`
	MainHeader      string  `json:"main_header"`
	MainDescription string  `json:"main_description"`
}

func MakeLesson(id goat.ID) Lesson {
	return Lesson{
		Name: fake.Title(),
		LessonData: &LessonData{
			MainHeader:      "Main Header",
			MainDescription: "Main Description",
		},
		Ordering: 0,
		CourseID: id,
	}
}
