package app

import (
	"github.com/68696c6c/goat"
	"github.com/jinzhu/gorm"
	"github.com/liamhendricks/auth-backend/src/repos"
	"github.com/liamhendricks/auth-backend/src/services"
	"github.com/sirupsen/logrus"
)

type ServiceContainer struct {
	Config          Config
	DB              *gorm.DB
	Logger          *logrus.Logger
	UserRepo        repos.UserRepo
	LessonRepo      repos.LessonRepo
	PasswordService services.PasswordService
	Errors          goat.ErrorHandler
}

var container ServiceContainer

func GetApp(c Config) (ServiceContainer, error) {
	db, err := goat.GetMainDB()
	if err != nil {
		return ServiceContainer{}, err
	}

	l := goat.GetLogger()
	ur := repos.NewUserRepoGorm(db, false)
	lr := repos.NewLessonRepoGorm(db)
	eh := goat.NewErrorHandler(l)
	ps := services.NewPasswordServiceAES(c.PasswordConfig)

	container = ServiceContainer{
		Config:          c,
		DB:              db,
		Logger:          l,
		UserRepo:        ur,
		LessonRepo:      lr,
		Errors:          eh,
		PasswordService: ps,
	}

	return container, nil
}
