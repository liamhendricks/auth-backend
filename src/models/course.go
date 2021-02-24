package models

import (
	"errors"

	"github.com/68696c6c/goat"
	"github.com/icrowley/fake"
)

type Course struct {
	goat.Model
	Name       string     `json:"name" binding:"required"`
	Lessons    []*Lesson  `json:"lessons" gorm:"ForeignKey:CourseID"`
	Users      []*User    `gorm:"many2many:user_courses;"`
	Max        int        `json:"max" binding:"required"`
	CourseType CourseType `json:"course_type" binding:"required"`
}

type CourseType string

const FreeCourse CourseType = "Free"
const PaidCourse CourseType = "Paid"

func CourseTypeFromString(s string) (CourseType, error) {
	if s == string(FreeCourse) || s == string(PaidCourse) {
		return CourseType(s), nil
	}
	return CourseType(""), errors.New("invalid course type")
}

func MakeCourse(free bool) Course {
	var t CourseType
	if free {
		t = FreeCourse
	} else {
		t = PaidCourse
	}
	return Course{
		Name:       fake.Title(),
		CourseType: t,
		Max:        20,
	}
}
