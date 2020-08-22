package repos

import (
	"os"
	"testing"

	"github.com/68696c6c/goat"
	gdb "github.com/68696c6c/goat/src/database"
	"github.com/jinzhu/gorm"
	"github.com/liamhendricks/auth-backend/src/models"
)

var Tc TestContainer
var Tf TestFixtures

type TestFixtures struct {
	Users   []models.User
	Lessons []*models.Lesson
	Courses []*models.Course
}

type TestContainer struct {
	UserRepo    UserRepo
	LessonRepo  LessonRepo
	SessionRepo SessionRepo
	CourseRepo  CourseRepo
	ResetRepo   ResetRepo
}

func TestMain(m *testing.M) {
	goat.Init()
	db := initRepoTestDB()

	r := NewUserRepoGorm(db, false)
	l := NewLessonRepoGorm(db)
	sr := NewSessionRepoGorm(db)
	c := NewCourseRepoGorm(db)
	rr := NewResetRepoGorm(db)
	Tc = TestContainer{
		UserRepo:    r,
		LessonRepo:  l,
		CourseRepo:  c,
		SessionRepo: sr,
		ResetRepo:   rr,
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
	var fc []*models.Course
	var pc []*models.Course

	freeCourse := models.MakeCourse(true)
	paidCourse := models.MakeCourse(false)
	persistFixture(db, &freeCourse)
	persistFixture(db, &paidCourse)
	fc = append(fc, &freeCourse)
	pc = append(pc, &paidCourse)

	for i := 0; i < num; i++ {
		lesson := models.MakeLesson()
		if i%2 == 0 {
			lesson.CourseID = freeCourse.ID
		} else {
			lesson.CourseID = paidCourse.ID
		}
		persistFixture(db, &lesson)
		l = append(l, &lesson)
	}

	for i := 0; i < num; i++ {
		user := models.MakeUser()
		if i%2 == 0 {
			user.Courses = fc
		} else {
			user.Courses = pc
		}
		persistFixture(db, &user)
		u = append(u, user)
	}

	Tf.Users = u
	Tf.Lessons = l
	Tf.Courses = fc
	Tf.Courses = append(Tf.Courses, pc...)
}

func persistFixture(db *gorm.DB, m interface{}) {
	if errs := db.Save(m).GetErrors(); len(errs) > 0 {
		panic(errs[0])
	}
}
