package repos

import (
	"github.com/68696c6c/goat"
	"github.com/jinzhu/gorm"
	"github.com/liamhendricks/auth-backend/src/models"
)

type SessionRepo interface {
	Save(model *models.Session) (errs []error)
	GetByTokenUser(token, userID goat.ID) (m models.Session, errs []error)
	Delete(id goat.ID) (errs []error)
}

type SessionRepoGorm struct {
	db *gorm.DB
}

func NewSessionRepoGorm(db *gorm.DB) SessionRepoGorm {
	return SessionRepoGorm{
		db: db,
	}
}

func (sr SessionRepoGorm) Save(m *models.Session) (errs []error) {
	if m.Model.ID.Valid() {
		errs = sr.db.Save(m).GetErrors()
	} else {
		errs = sr.db.Create(m).GetErrors()
	}
	return
}

func (sr SessionRepoGorm) GetByTokenUser(token, userID goat.ID) (m models.Session, errs []error) {
	errs = sr.db.First(&m, "token = ? and user_id = ?", token, userID).GetErrors()
	return
}

func (sr SessionRepoGorm) Delete(id goat.ID) (errs []error) {
	errs = sr.db.Delete(&models.Session{}, "id = ?", id).GetErrors()
	return
}

type SessionRepoDummy struct {
	User *models.User
}

func NewSessionRepoDummy(user *models.User) SessionRepoDummy {
	return SessionRepoDummy{
		User: user,
	}
}

func (sr SessionRepoDummy) Save(m *models.Session) (errs []error) {
	return nil
}

func (sr SessionRepoDummy) GetByTokenUser(token, userID goat.ID) (m models.Session, errs []error) {
	return *sr.User.Session, nil
}

func (sr SessionRepoDummy) Delete(id goat.ID) (errs []error) {
	return nil
}
