package repos

import (
	"github.com/68696c6c/goat"
	"github.com/jinzhu/gorm"
	"github.com/liamhendricks/auth-backend/src/models"
)

type ResetRepo interface {
	Save(model *models.Reset) (errs []error)
	GetByToken(token goat.ID) (m models.Reset, errs []error)
	Delete(id goat.ID) (errs []error)
}

type ResetRepoGorm struct {
	db *gorm.DB
}

func NewResetRepoGorm(db *gorm.DB) ResetRepoGorm {
	return ResetRepoGorm{
		db: db,
	}
}

func (r ResetRepoGorm) Save(m *models.Reset) (errs []error) {
	if m.Model.ID.Valid() {
		errs = r.db.Save(m).GetErrors()
	} else {
		errs = r.db.Create(m).GetErrors()
	}
	return
}

func (r ResetRepoGorm) GetByToken(token goat.ID) (m models.Reset, errs []error) {
	errs = r.db.First(&m, "token = ?", token).GetErrors()
	return
}

func (l ResetRepoGorm) Delete(id goat.ID) (errs []error) {
	errs = l.db.Delete(&models.Reset{}, "id = ?", id).GetErrors()
	return
}
