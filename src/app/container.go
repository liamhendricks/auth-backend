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
	DateRepo        repos.DateRepo
	EmailListRepo   repos.EmailListRepo
	PasswordService services.PasswordService
	SessionService  services.SessionService
	MailService     services.MailService
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
	dr := repos.NewDateRepoGorm(db)
	cr := repos.NewCourseRepoGorm(db)
	rr := repos.NewResetRepoGorm(db)
	sr := repos.NewSessionRepoGorm(db)
	er := repos.NewEmailListRepoGorm(db)
	eh := goat.NewErrorHandler(l)
	ps := services.NewPasswordServiceBcrypt()
	ss := services.NewSessionServiceDB(sr)
	es := services.NewSendgridMailer(c.SendgridFromEmail, c.SendgridFromName, c.SendgridBaseURL, c.SendgridSecretKey, c.ResetTemplateID, c.PurchaseTemplateID, c.SignupTemplateID)

	container = ServiceContainer{
		Config:          c,
		DB:              db,
		Logger:          l,
		UserRepo:        ur,
		LessonRepo:      lr,
		CourseRepo:      cr,
		ResetRepo:       rr,
		EmailListRepo:   er,
		DateRepo:        dr,
		Errors:          eh,
		PasswordService: ps,
		SessionService:  ss,
		MailService:     es,
	}

	return container, nil
}
