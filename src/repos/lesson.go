package repos

import (
	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/query"
	"github.com/liamhendricks/auth-backend/src/models"
)

type LessonRepo interface {
	Save(model *models.Lesson) (errs []error)
	GetAll(q *query.Query) (m []*models.Lesson, errs []error)
	GetByID(id goat.ID, load bool) (m models.Lesson, errs []error)
	Delete(id goat.ID) (errs []error)
}

type LessonRepoGorm struct {
}

func (l LessonRepoGorm) Save(model *models.Lesson) (errs []error) {
	return
}

func (l LessonRepoGorm) GetAll(q *query.Query) (m []*models.Lesson, errs []error) {
	return
}

func (l LessonRepoGorm) GetByID(id goat.ID, load bool) (m models.Lesson, errs []error) {
	return
}

func (l LessonRepoGorm) Delete(id goat.ID) (errs []error) {
	return
}
