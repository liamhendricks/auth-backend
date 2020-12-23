package repos

import (
	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/query"
	"github.com/jinzhu/gorm"
	"github.com/liamhendricks/auth-backend/src/models"
)

type UserRepo interface {
	Save(model *models.User) (errs []error)
	GetAll(q *query.Query) (m []*models.User, errs []error)
	GetByID(id goat.ID, load bool) (m models.User, errs []error)
	GetByNameOrEmail(username, email string, load bool) (m models.User, errs []error)
	Delete(id goat.ID) (errs []error)
	Clear(m *models.User, assoc string)
}

func NewUserRepoGorm(db *gorm.DB, disable bool) UsersRepoGorm {
	return UsersRepoGorm{
		db:          db,
		disableAuth: disable,
	}
}

type UsersRepoGorm struct {
	db          *gorm.DB
	disableAuth bool
}

func (u UsersRepoGorm) Save(m *models.User) (errs []error) {
	if m.Model.ID.Valid() {
		errs = u.db.Save(m).GetErrors()
	} else {
		errs = u.db.Create(m).GetErrors()
	}
	return
}

func (u UsersRepoGorm) GetAll(q *query.Query) (m []*models.User, errs []error) {
	base := u.db.Model(&models.User{})
	qr, err := q.ApplyToGorm(base)
	if err != nil {
		return m, []error{err}
	}

	errs = qr.Find(&m).GetErrors()
	return m, errs
}

func (u UsersRepoGorm) GetByID(id goat.ID, load bool) (m models.User, errs []error) {
	q := u.db
	if load {
		q = q.Preload("Courses.Lessons.LessonData").Preload("Session").Preload("Reset")
	}
	errs = q.First(&m, "id = ?", id).GetErrors()
	return
}

func (u UsersRepoGorm) GetByNameOrEmail(username, email string, load bool) (m models.User, errs []error) {
	q := u.db
	if load {
		q = q.Preload("Courses").Preload("Courses.Lessons").Preload("Session").Preload("Reset")
	}
	errs = q.First(&m, "name = ? or email = ?", username, email).GetErrors()
	return
}

func (u UsersRepoGorm) Delete(id goat.ID) (errs []error) {
	errs = u.db.Delete(&models.User{}, "id = ?", id).GetErrors()
	return
}

func (u UsersRepoGorm) Clear(m *models.User, assoc string) {
	u.db.Model(m).Association(assoc).Clear()
}
