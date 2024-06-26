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
	List    []models.EmailList
	Lessons []*models.Lesson
	Courses []*models.Course
	Dates   []*models.Date
}

type TestContainer struct {
	UserRepo      UserRepo
	EmailListRepo EmailListRepo
	LessonRepo    LessonRepo
	SessionRepo   SessionRepo
	CourseRepo    CourseRepo
	ResetRepo     ResetRepo
	DateRepo      DateRepo
}

func TestMain(m *testing.M) {
	goat.Init()
	db := initRepoTestDB()

	r := NewUserRepoGorm(db, false)
	l := NewLessonRepoGorm(db)
	el := NewEmailListRepoGorm(db)
	sr := NewSessionRepoGorm(db)
	c := NewCourseRepoGorm(db)
	rr := NewResetRepoGorm(db)
	dr := NewDateRepoGorm(db)
	Tc = TestContainer{
		EmailListRepo: el,
		UserRepo:      r,
		LessonRepo:    l,
		CourseRepo:    c,
		SessionRepo:   sr,
		ResetRepo:     rr,
		DateRepo:      dr,
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
	var el []models.EmailList
	var l []*models.Lesson
	var fc []*models.Course
	var pc []*models.Course
	var d []*models.Date

	freeCourse := models.MakeCourse(true)
	paidCourse := models.MakeCourse(false)
	persistFixture(db, &freeCourse)
	persistFixture(db, &paidCourse)
	fc = append(fc, &freeCourse)
	pc = append(pc, &paidCourse)

	for i := 0; i < 5; i++ {
		pd := models.MakeDate()
		pd.CourseID = paidCourse.ID
		persistFixture(db, &pd)

		fd := models.MakeDate()
		fd.CourseID = freeCourse.ID
		persistFixture(db, &fd)
		d = append(d, &pd, &fd)
	}

	for i := 0; i < num; i++ {
		var lesson models.Lesson
		if i%2 == 0 {
			lesson = models.MakeLesson(freeCourse.ID)
		} else {
			lesson = models.MakeLesson(paidCourse.ID)
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

		listMember := models.MakeEmailList()
		persistFixture(db, &listMember)
		el = append(el, listMember)
	}

	Tf.Users = u
	Tf.List = el
	Tf.Lessons = l
	Tf.Courses = fc
	Tf.Dates = d
	Tf.Courses = append(Tf.Courses, pc...)
}

func persistFixture(db *gorm.DB, m interface{}) {
	if errs := db.Save(m).GetErrors(); len(errs) > 0 {
		panic(errs[0])
	}
}
