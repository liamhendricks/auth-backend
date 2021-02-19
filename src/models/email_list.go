package models

import "github.com/68696c6c/goat"

type EmailList struct {
	goat.Model
	Email string `json:"email" binding:"required"`
	Name  string `json:"name" binding:"required"`
}
