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
	EmailListRepo   repos.EmailListRepo
	PasswordService services.PasswordService
	SessionService  services.SessionService
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
	sr := repos.NewSessionRepoGorm(db)
	er := repos.NewEmailListRepoGorm(db)
	eh := goat.NewErrorHandler(l)
	ps := services.NewPasswordServiceBcrypt()
	ss := services.NewSessionServiceDB(sr)

	container = ServiceContainer{
		Config:          c,
		DB:              db,
		Logger:          l,
		UserRepo:        ur,
		LessonRepo:      lr,
		CourseRepo:      cr,
		ResetRepo:       rr,
		EmailListRepo:   er,
		Errors:          eh,
		PasswordService: ps,
		SessionService:  ss,
	}

	return container, nil
}
