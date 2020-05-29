package app

import (
	"github.com/68696c6c/goat"
	"github.com/jinzhu/gorm"
	"github.com/liamhendricks/auth-backend/src/repos"
	"github.com/sirupsen/logrus"
)

type ServiceContainer struct {
	Config   Config
	DB       *gorm.DB
	Logger   *logrus.Logger
	UserRepo repos.UserRepo
}

var container ServiceContainer

func GetApp(c Config) (ServiceContainer, error) {
	db, err := goat.GetMainDB()
	if err != nil {
		return ServiceContainer{}, err
	}

	l := goat.GetLogger()

	container = ServiceContainer{
		Config: c,
		DB:     db,
		Logger: l,
	}

	return container, nil
}
