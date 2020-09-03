package services

import (
	"errors"
	"time"

	"github.com/68696c6c/goat"
	"github.com/liamhendricks/auth-backend/src/models"
	"github.com/liamhendricks/auth-backend/src/repos"
)

type SessionService interface {
	Start(user *models.User) error
	Refresh(user *models.User, length time.Duration) error
	Stop(user *models.User) error
	Valid(user *models.User, token goat.ID) bool
}

type SessionServiceDB struct {
	sessionRepo repos.SessionRepo
}

func NewSessionServiceDB(sr repos.SessionRepo) SessionServiceDB {
	return SessionServiceDB{
		sessionRepo: sr,
	}
}

func (ss SessionServiceDB) Start(user *models.User) error {
	var s *models.Session
	if user.Session != nil {
		user.Session.Expiration = time.Now().Add(20 * time.Minute)
		s = user.Session
	} else {
		//TODO: configurable session length
		token := goat.NewID()
		s = &models.Session{
			Token:      token,
			UserID:     user.ID,
			Expiration: time.Now().Add(20 * time.Minute),
		}
	}

	errs := ss.sessionRepo.Save(s)
	if len(errs) > 0 {
		return errors.New("failed saving session")
	}

	user.Session = s
	return nil
}

func (ss SessionServiceDB) Refresh(user *models.User, length time.Duration) error {
	if user.Session == nil {
		return errors.New("no session exists for this user")
	}

	user.Session.Expiration = time.Now().Add(length)
	errs := ss.sessionRepo.Save(user.Session)
	if len(errs) > 0 {
		return errors.New("failed saving session")
	}

	return nil
}

func (ss SessionServiceDB) Stop(user *models.User) error {
	if user.Session == nil {
		return errors.New("no session exists for this user")
	}

	errs := ss.sessionRepo.Delete(user.Session.ID)
	if len(errs) > 0 {
		return errors.New("error stopping session")
	}

	return nil
}

func (ss SessionServiceDB) Valid(user *models.User, token goat.ID) bool {
	if user.Session == nil {
		return false
	}

	if time.Now().After(user.Session.Expiration) {
		return false
	}

	if user.Session.Token != token {
		return false
	}

	return true
}
