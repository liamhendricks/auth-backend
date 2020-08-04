package repos

import (
	"testing"

	"github.com/68696c6c/goat/query"
	"github.com/liamhendricks/auth-backend/src/models"
	"github.com/stretchr/testify/require"
)

func TestLessonsRepoSave(t *testing.T) {
	l := &models.Lesson{
		Name: "TestLesson",
	}
	errs := tc.LessonRepo.Save(l)
	require.Nil(t, errs)
	theLesson, errs := tc.LessonRepo.GetByID(l.ID)
	require.Equal(t, l.ID, theLesson.ID)
}

func TestLessonssRepoUpdate(t *testing.T) {
	l := tf.Lessons[1]
	l.Name = "NewName"
	errs := tc.LessonRepo.Save(l)
	require.Nil(t, errs)
	theLesson, errs := tc.LessonRepo.GetByID(l.ID)
	require.Equal(t, l.ID, theLesson.ID)
	require.Equal(t, theLesson.Name, "NewName")
}

func TestLessonsRepoGetByID(t *testing.T) {
	id := tf.Lessons[0].ID
	lesson, errs := tc.LessonRepo.GetByID(id)
	require.Nil(t, errs)
	require.Equal(t, id, lesson.ID)
}

func TestLessonsRepoGetAll(t *testing.T) {
	l, errs := tc.LessonRepo.GetAll(&query.Query{})
	require.Empty(t, errs)
	require.GreaterOrEqual(t, len(l), 1)
}

func TestLessonsRepoDelete(t *testing.T) {
	id := tf.Lessons[2].ID
	errs := tc.LessonRepo.Delete(id)
	require.Empty(t, errs)
	_, errs = tc.LessonRepo.GetByID(id)
	require.NotNil(t, errs)
}

func TestAttachLessonToUser(t *testing.T) {
	lesson := &models.Lesson{
		Name:       "BrandNewLesson",
		LessonType: models.FreeLesson,
	}

	lessons, _ := tc.LessonRepo.GetAll(&query.Query{})
	totalLessons := len(lessons)

	u := tf.Users[0]
	u.Lessons = append(u.Lessons, lesson)
	errs := tc.UserRepo.Save(&u)
	require.Empty(t, errs)

	user, errs := tc.UserRepo.GetByID(u.ID, true)
	require.Empty(t, errs)
	require.NotNil(t, user.Lessons)

	has := false
	for _, l := range user.Lessons {
		if l.Name == "BrandNewLesson" {
			has = true
		}
	}

	// it should also save the new lesson
	lessons, _ = tc.LessonRepo.GetAll(&query.Query{})
	newTotalLessons := len(lessons)
	require.GreaterOrEqual(t, newTotalLessons, totalLessons)

	require.Equal(t, has, true)
}
