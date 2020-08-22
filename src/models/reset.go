package models

import (
	"time"

	"github.com/68696c6c/goat"
)

type Reset struct {
	goat.Model
	Token      goat.ID   `json:"token" binding:"required"`
	UserID     goat.ID   `json:"user_id" binding:"required"`
	Expiration time.Time `json:"expiration" binding:"required"`
}
