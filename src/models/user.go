package models

import "github.com/68696c6c/goat"

type User struct {
	goat.Model
	Name string
	Tier int
}
