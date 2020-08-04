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
	Users   []models.User
	Lessons []*models.Lesson
}

type TestContainer struct {
	UserRepo   UserRepo
	LessonRepo LessonRepo
}

func TestMain(m *testing.M) {
	goat.Init()
	db := initRepoTestDB()

	r := NewUserRepoGorm(db, false)
	l := NewLessonRepoGorm(db)
	tc = TestContainer{
		UserRepo:   r,
		LessonRepo: l,
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
	var l []*models.Lesson

	for i := 0; i < num; i++ {
		lesson := models.MakeLesson()
		persistFixture(db, &lesson)
		l = append(l, &lesson)
	}

	for i := 0; i < num; i++ {
		user := models.MakeUser()
		persistFixture(db, &user)
		user.Lessons = l
		u = append(u, user)
	}

	tf.Users = u
	tf.Lessons = l
}

func persistFixture(db *gorm.DB, m interface{}) {
	if errs := db.Save(m).GetErrors(); len(errs) > 0 {
		panic(errs[0])
	}
}
