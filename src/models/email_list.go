package models

import (
	"github.com/68696c6c/goat"
	"github.com/icrowley/fake"
)

type EmailList struct {
	goat.Model
	Email string `json:"email" binding:"required"`
	Name  string `json:"name" binding:"required"`
}

func MakeEmailList() EmailList {
	return EmailList{
		Name:  fake.FullName(),
		Email: fake.EmailAddress(),
	}
}

func (el *EmailList) TableName() string {
	return "email_list"
}
