package models

import (
	"time"

	"github.com/68696c6c/goat"
	"github.com/icrowley/fake"
)

type Date struct {
	goat.Model
	CourseID goat.ID `json:"course_id"`
	Name     string  `json:"name"`
	Date     string  `json:"date"`
}

func MakeDate() Date {
	return Date{
		Name: fake.CharactersN(5),
		Date: time.Now().String(),
	}
}
