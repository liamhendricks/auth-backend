package repos

import (
	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/query"
	"github.com/jinzhu/gorm"
	"github.com/liamhendricks/auth-backend/src/models"
)

type EmailListRepo interface {
	Save(model *models.EmailList) (errs []error)
	GetByID(id goat.ID) (m models.EmailList, errs []error)
	GetAll(q *query.Query) (m []*models.EmailList, errs []error)
	Delete(id goat.ID) (errs []error)
}

type EmailListRepoGorm struct {
	db *gorm.DB
}

func NewEmailListRepoGorm(db *gorm.DB) EmailListRepoGorm {
	return EmailListRepoGorm{
		db: db,
	}
}

func (el EmailListRepoGorm) Save(m *models.EmailList) (errs []error) {
	if m.Model.ID.Valid() {
		errs = el.db.Save(m).GetErrors()
	} else {
		errs = el.db.Create(m).GetErrors()
	}
	return
}

func (el EmailListRepoGorm) GetByID(id goat.ID) (m models.EmailList, errs []error) {
	q := el.db
	errs = q.First(&m, "id = ?", id).GetErrors()
	return
}

func (el EmailListRepoGorm) GetAll(q *query.Query) (m []*models.EmailList, errs []error) {
	base := el.db.Model(&models.EmailList{})
	qr, err := q.ApplyToGorm(base)
	if err != nil {
		return m, []error{err}
	}

	errs = qr.Find(&m).GetErrors()
	return m, errs
}

func (el EmailListRepoGorm) Delete(id goat.ID) (errs []error) {
	errs = el.db.Delete(&models.EmailList{}, "id = ?", id).GetErrors()
	return
}
