package repos

import (
	"testing"

	"github.com/68696c6c/goat/query"
	"github.com/liamhendricks/auth-backend/src/models"
	"github.com/stretchr/testify/require"
)

func TestCoursesRepoSave(t *testing.T) {
	var l []*models.Lesson
	l = append(l, Tf.Lessons[0])
	c := &models.Course{
		Name:    "TesTcourse",
		Lessons: l,
	}
	errs := Tc.CourseRepo.Save(c)
	require.Nil(t, errs)
	theCourse, errs := Tc.CourseRepo.GetByID(c.ID, true)
	require.Equal(t, c.ID, theCourse.ID)
	require.Equal(t, c.Lessons[0].ID, Tf.Lessons[0].ID)
}

func TestCoursessRepoUpdate(t *testing.T) {
	c := Tf.Courses[1]
	c.Name = "NewName"
	errs := Tc.CourseRepo.Save(c)
	require.Nil(t, errs)
	theCourse, errs := Tc.CourseRepo.GetByID(c.ID, false)
	require.Equal(t, c.ID, theCourse.ID)
	require.Equal(t, theCourse.Name, "NewName")
}

func TestCoursesRepoGetByID(t *testing.T) {
	id := Tf.Courses[0].ID
	max := Tf.Courses[0].Max
	Course, errs := Tc.CourseRepo.GetByID(id, true)
	require.Nil(t, errs)
	require.Equal(t, id, Course.ID)
	require.Equal(t, max, Course.Max)
}

func TestCoursesRepoGetByName(t *testing.T) {
	name := Tf.Courses[0].Name
	Course, errs := Tc.CourseRepo.GetByName(name, true)
	require.Nil(t, errs)
	require.Equal(t, name, Course.Name)
}

func TestCoursesRepoGetAll(t *testing.T) {
	c, errs := Tc.CourseRepo.GetAll(&query.Query{})
	require.Empty(t, errs)
	require.GreaterOrEqual(t, len(c), 1)
}

func TestGetUserCourses(t *testing.T) {
	id := Tf.Courses[0].ID
	course, errs := Tc.CourseRepo.GetByID(id, true)
	require.Empty(t, errs)
	require.NotNil(t, course.Lessons)
	for _, l := range course.Lessons {
		println(l)
	}
}

func TestCoursesRepoDelete(t *testing.T) {
	id := Tf.Courses[1].ID
	errs := Tc.CourseRepo.Delete(id)
	require.Empty(t, errs)
	_, errs = Tc.CourseRepo.GetByID(id, true)
	require.NotNil(t, errs)
}
