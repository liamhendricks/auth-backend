package repos

import (
	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/query"
	"github.com/jinzhu/gorm"
	"github.com/liamhendricks/auth-backend/src/models"
)

type CourseRepo interface {
	Save(model *models.Course) (errs []error)
	GetAll(q *query.Query) (m []*models.Course, errs []error)
	GetByID(id goat.ID, load bool) (m models.Course, errs []error)
	GetByName(name string, load bool) (m models.Course, errs []error)
	Delete(id goat.ID) (errs []error)
	Clear(m *models.Course, assoc string)
}

type CourseRepoGorm struct {
	db *gorm.DB
}

func NewCourseRepoGorm(db *gorm.DB) CourseRepoGorm {
	return CourseRepoGorm{
		db: db,
	}
}

func (c CourseRepoGorm) Save(m *models.Course) (errs []error) {
	if m.Model.ID.Valid() {
		errs = c.db.Save(m).GetErrors()
	} else {
		errs = c.db.Create(m).GetErrors()
	}
	return
}

func (c CourseRepoGorm) GetAll(q *query.Query) (m []*models.Course, errs []error) {
	base := c.db.Model(&models.Course{})
	qr, err := q.ApplyToGorm(base)
	if err != nil {
		return m, []error{err}
	}

	errs = qr.Find(&m).GetErrors()
	return m, errs
}

func (c CourseRepoGorm) GetByID(id goat.ID, load bool) (m models.Course, errs []error) {
	q := c.db
	if load {
		q = q.Preload("Lessons").Preload("Users").Preload("Dates")
	}
	errs = q.First(&m, "id = ?", id).GetErrors()
	return
}

func (c CourseRepoGorm) GetByName(name string, load bool) (m models.Course, errs []error) {
	q := c.db
	if load {
		q = q.Preload("Lessons").Preload("Users").Preload("Dates")
	}
	errs = q.First(&m, "name = ?", name).GetErrors()
	return
}

func (c CourseRepoGorm) Delete(id goat.ID) (errs []error) {
	errs = c.db.Delete(&models.Course{}, "id = ?", id).GetErrors()
	return
}

func (u CourseRepoGorm) Clear(m *models.Course, assoc string) {
	u.db.Model(m).Association(assoc).Clear()
}
