package models

import (
	"github.com/68696c6c/goat"
	"github.com/icrowley/fake"
)

type User struct {
	goat.Model
	Name     string    `json:"name" binding:"required"`
	Email    string    `json:"email" binding:"required"`
	Password string    `json:"password" binding:"required"`
	UserType UserType  `json:"user_type" binding:"required"`
	Courses  []*Course `gorm:"many2many:user_courses;"`
	Session  *Session
	Reset    *Reset
}

type UserType string

const AdminUser UserType = "Admin"
const FreeUser UserType = "Free"
const PaidUser UserType = "Paid"

func MakeUser() User {
	return User{
		Name:     fake.FullName(),
		Email:    fake.EmailAddress(),
		Password: fake.Password(5, 16, true, true, true),
		UserType: FreeUser,
	}
}
