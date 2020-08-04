package models

import (
	"github.com/68696c6c/goat"
	"github.com/icrowley/fake"
)

type Lesson struct {
	goat.Model
	Name       string     `json:"name" binding:"required"`
	LessonType LessonType `json:"lesson_type" binding:"required"`
	Users      []*User    `gorm:"many2many:user_lessons;"`
}

type LessonType string

const FreeLesson LessonType = "Free"
const PaidLesson LessonType = "Paid"

func MakeLesson() Lesson {
	return Lesson{
		Name:       fake.Title(),
		LessonType: FreeLesson,
	}
}
