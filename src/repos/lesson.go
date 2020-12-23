package repos

import (
	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/query"
	"github.com/jinzhu/gorm"
	"github.com/liamhendricks/auth-backend/src/models"
)

type LessonRepo interface {
	Save(model *models.Lesson) (errs []error)
	GetAll(q *query.Query) (m []*models.Lesson, errs []error)
	GetByID(id goat.ID) (m models.Lesson, errs []error)
	GetByName(name string) (m models.Lesson, errs []error)
	Delete(id goat.ID) (errs []error)
	Clear(m *models.Lesson, assoc string)
}

type LessonRepoGorm struct {
	db *gorm.DB
}

func NewLessonRepoGorm(db *gorm.DB) LessonRepoGorm {
	return LessonRepoGorm{
		db: db,
	}
}

func (l LessonRepoGorm) Save(m *models.Lesson) (errs []error) {
	if m.Model.ID.Valid() {
		errs = l.db.Save(m).GetErrors()
	} else {
		errs = l.db.Create(m).GetErrors()
	}
	return
}

func (l LessonRepoGorm) GetAll(q *query.Query) (m []*models.Lesson, errs []error) {
	base := l.db.Model(&models.Lesson{})
	qr, err := q.ApplyToGorm(base)
	if err != nil {
		return m, []error{err}
	}

	errs = qr.Find(&m).Preload("LessonData").GetErrors()
	return m, errs
}

func (l LessonRepoGorm) GetByID(id goat.ID) (m models.Lesson, errs []error) {
	errs = l.db.First(&m, "id = ?", id).Preload("LessonData").GetErrors()
	return
}

func (l LessonRepoGorm) GetByName(name string) (m models.Lesson, errs []error) {
	errs = l.db.First(&m, "name = ?", name).GetErrors()
	return
}

func (l LessonRepoGorm) Delete(id goat.ID) (errs []error) {
	errs = l.db.Delete(&models.Lesson{}, "id = ?", id).GetErrors()
	return
}

func (u LessonRepoGorm) Clear(m *models.Lesson, assoc string) {
	u.db.Model(m).Association(assoc).Clear()
}
