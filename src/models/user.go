package models

import (
	"github.com/68696c6c/goat"
	"github.com/icrowley/fake"
)

type User struct {
	goat.Model
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func MakeUser() User {
	return User{
		Name:     fake.FullName(),
		Email:    fake.EmailAddress(),
		Password: fake.Password(5, 16, true, true, true),
	}
}
