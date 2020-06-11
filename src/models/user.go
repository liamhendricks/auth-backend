package models

import (
	"github.com/68696c6c/goat"
	"github.com/icrowley/fake"
)

type User struct {
	goat.Model
	Name string
}

func MakeUser() User {
	return User{
		Name: fake.FullName(),
	}
}
