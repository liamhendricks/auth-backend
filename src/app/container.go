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
	CourseRepo      repos.CourseRepo
	ResetRepo       repos.ResetRepo
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
	cr := repos.NewCourseRepoGorm(db)
	rr := repos.NewResetRepoGorm(db)
	eh := goat.NewErrorHandler(l)
	ps := services.NewPasswordServiceAES(c.PasswordConfig)

	container = ServiceContainer{
		Config:          c,
		DB:              db,
		Logger:          l,
		UserRepo:        ur,
		LessonRepo:      lr,
		CourseRepo:      cr,
		ResetRepo:       rr,
		Errors:          eh,
		PasswordService: ps,
	}

	return container, nil
}
