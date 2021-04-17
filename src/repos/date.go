package repos

import (
	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/query"
	"github.com/jinzhu/gorm"
	"github.com/liamhendricks/auth-backend/src/models"
)

type DateRepo interface {
	Save(models *models.Date) (errs []error)
	GetByID(id goat.ID) (m models.Date, errs []error)
	GetAll(q *query.Query) (m []*models.Date, errs []error)
	Delete(id goat.ID) (errs []error)
}

type DateRepoGorm struct {
	db *gorm.DB
}

func NewDateRepoGorm(db *gorm.DB) DateRepoGorm {
	return DateRepoGorm{
		db: db,
	}
}

func (d DateRepoGorm) Save(m *models.Date) (errs []error) {
	if m.ID.Valid() {
		errs = d.db.Save(m).GetErrors()
	} else {
		errs = d.db.Create(m).GetErrors()
	}
	return
}

func (d DateRepoGorm) GetByID(id goat.ID) (m models.Date, errs []error) {
	q := d.db
	errs = q.First(&m, "id = ?", id).GetErrors()
	return
}

func (d DateRepoGorm) GetAll(q *query.Query) (m []*models.Date, errs []error) {
	base := d.db.Model(&models.Date{})
	qr, err := q.ApplyToGorm(base)
	if err != nil {
		return m, []error{err}
	}

	errs = qr.Find(&m).GetErrors()
	return m, errs
}

func (d DateRepoGorm) Delete(id goat.ID) (errs []error) {
	errs = d.db.Delete(&models.Date{}, "id = ?", id).GetErrors()
	return
}
