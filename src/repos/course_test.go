package repos

import (
	"testing"

	"github.com/68696c6c/goat/query"
	"github.com/liamhendricks/auth-backend/src/models"
	"github.com/stretchr/testify/require"
)

func TestCoursesRepoSave(t *testing.T) {
	var l []*models.Lesson
	l = append(l, tf.Lessons[0])
	c := &models.Course{
		Name:    "TestCourse",
		Lessons: l,
	}
	errs := tc.CourseRepo.Save(c)
	require.Nil(t, errs)
	theCourse, errs := tc.CourseRepo.GetByID(c.ID, true)
	require.Equal(t, c.ID, theCourse.ID)
	require.Equal(t, c.Lessons[0].ID, tf.Lessons[0].ID)
}

func TestCoursessRepoUpdate(t *testing.T) {
	c := tf.Courses[1]
	c.Name = "NewName"
	errs := tc.CourseRepo.Save(c)
	require.Nil(t, errs)
	theCourse, errs := tc.CourseRepo.GetByID(c.ID, false)
	require.Equal(t, c.ID, theCourse.ID)
	require.Equal(t, theCourse.Name, "NewName")
}

func TestCoursesRepoGetByID(t *testing.T) {
	id := tf.Courses[0].ID
	Course, errs := tc.CourseRepo.GetByID(id, true)
	require.Nil(t, errs)
	require.Equal(t, id, Course.ID)
}

func TestCoursesRepoGetAll(t *testing.T) {
	c, errs := tc.CourseRepo.GetAll(&query.Query{})
	require.Empty(t, errs)
	require.GreaterOrEqual(t, len(c), 1)
}

func TestGetUserCourses(t *testing.T) {
	id := tf.Courses[0].ID
	course, errs := tc.CourseRepo.GetByID(id, true)
	require.Empty(t, errs)
	require.NotNil(t, course.Lessons)
	for _, l := range course.Lessons {
		println(l)
	}
}

func TestCoursesRepoDelete(t *testing.T) {
	id := tf.Courses[1].ID
	errs := tc.CourseRepo.Delete(id)
	require.Empty(t, errs)
	_, errs = tc.CourseRepo.GetByID(id, true)
	require.NotNil(t, errs)
}
