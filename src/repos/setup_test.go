package repos

import (
	"os"
	"testing"

	"github.com/68696c6c/goat"
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

func TestMain(m *testing.M) {
	goat.Init()
	db := initRepoTestDB()

	r := NewUserRepoGorm(db, false)
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
