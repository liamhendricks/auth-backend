package models

import (
	"time"

	"github.com/68696c6c/goat"
)

type Session struct {
	goat.Model
	Token      goat.ID   `json:"token" binding:"required"`
	UserID     goat.ID   `json:"user_id" binding:"required"`
	Expiration time.Time `json:"expiration" binding:"required"`
}

func MakeSession() *Session {
	return &Session{
		Token:      goat.NewID(),
		Expiration: time.Now().Add(30 * time.Minute),
	}
}
