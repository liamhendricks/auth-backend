package repos

import (
	"os"
	"testing"

	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/query"
	gdb "github.com/68696c6c/goat/src/database"
	"github.com/jinzhu/gorm"
	"github.com/liamhendricks/auth-backend/src/models"
)

var tc TestContainer
var tf TestFixtures

type TestFixtures struct {
	Users []models.User
}

type TestContainer struct {
	UserRepo UserRepo
}

type UsersRepoTest struct {
	db *gorm.DB
}

func (u UsersRepoTest) Save(m *models.User) (errs []error) {
	if m.Model.ID.Valid() {
		errs = u.db.Save(m).GetErrors()
	} else {
		errs = u.db.Create(m).GetErrors()
	}
	return
}

func (u UsersRepoTest) GetAll(q *query.Query) (m []*models.User, errs []error) {
	base := u.db.Model(&models.User{})
	qr, err := q.ApplyToGorm(base)
	if err != nil {
		return m, []error{err}
	}

	errs = qr.Find(&m).GetErrors()
	return m, errs
}

func (u UsersRepoTest) GetByID(id goat.ID, load bool) (m models.User, errs []error) {
	errs = u.db.First(&m, "id = ?", id).GetErrors()
	return
}

func TestMain(m *testing.M) {
	goat.Init()
	db := initRepoTestDB()

	r := UsersRepoTest{
		db: db,
	}

	tc = TestContainer{
		UserRepo: r,
	}

	seedTests(5, db)

	os.Exit(m.Run())
}

func initRepoTestDB() *gorm.DB {
	db, err := goat.GetCustomDB(gdb.ConnectionConfig{
		Debug:           true,
		Host:            "auth_db",
		Port:            3306,
		Database:        "auth_test",
		Username:        "root",
		Password:        "secret",
		MultiStatements: false,
	})

	if err != nil {
		panic(err)
	}

	schema, err := goat.GetSchema(db)
	if err != nil {
		panic(err)
	}

	_, _, err = schema.Reset()
	if err != nil {
		panic(err)
	}

	return db
}

func seedTests(num int, db *gorm.DB) {
	var u []models.User

	for i := 0; i < num; i++ {
		user := models.MakeUser()
		persistFixture(db, &user)
		u = append(u, user)
	}

	tf.Users = u
}

func persistFixture(db *gorm.DB, m interface{}) {
	if errs := db.Save(m).GetErrors(); len(errs) > 0 {
		panic(errs[0])
	}
}
