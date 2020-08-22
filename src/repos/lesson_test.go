package repos

import (
	"testing"

	"github.com/68696c6c/goat/query"
	"github.com/liamhendricks/auth-backend/src/models"
	"github.com/stretchr/testify/require"
)

func TestLessonsRepoSave(t *testing.T) {
	l := &models.Lesson{
		Name:     "TestLesson",
		CourseID: Tf.Courses[0].ID,
	}
	errs := Tc.LessonRepo.Save(l)
	require.Nil(t, errs)
	theLesson, errs := Tc.LessonRepo.GetByID(l.ID)
	require.Equal(t, l.ID, theLesson.ID)
}

func TestLessonssRepoUpdate(t *testing.T) {
	l := Tf.Lessons[1]
	l.Name = "NewName"
	errs := Tc.LessonRepo.Save(l)
	require.Nil(t, errs)
	theLesson, errs := Tc.LessonRepo.GetByID(l.ID)
	require.Equal(t, l.ID, theLesson.ID)
	require.Equal(t, theLesson.Name, "NewName")
}

func TestLessonsRepoGetByID(t *testing.T) {
	id := Tf.Lessons[0].ID
	lesson, errs := Tc.LessonRepo.GetByID(id)
	require.Nil(t, errs)
	require.Equal(t, id, lesson.ID)
}

func TestLessonsRepoGetByName(t *testing.T) {
	name := Tf.Lessons[0].Name
	lesson, errs := Tc.LessonRepo.GetByName(name)
	require.Nil(t, errs)
	require.Equal(t, name, lesson.Name)
}

func TestLessonsRepoGetAll(t *testing.T) {
	l, errs := Tc.LessonRepo.GetAll(&query.Query{})
	require.Empty(t, errs)
	require.GreaterOrEqual(t, len(l), 1)
}

func TestLessonsRepoDelete(t *testing.T) {
	id := Tf.Lessons[2].ID
	errs := Tc.LessonRepo.Delete(id)
	require.Empty(t, errs)
	_, errs = Tc.LessonRepo.GetByID(id)
	require.NotNil(t, errs)
}
