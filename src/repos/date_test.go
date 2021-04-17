package repos

import (
	"testing"

	"github.com/68696c6c/goat/query"
	"github.com/liamhendricks/auth-backend/src/models"
	"github.com/stretchr/testify/require"
)

func TestDateRepoSave(t *testing.T) {
	id := Tf.Courses[0].ID
	date := models.MakeDate()
	date.CourseID = id
	errs := Tc.DateRepo.Save(&date)
	require.Empty(t, errs)
}

func TestDateRepoGetAll(t *testing.T) {
	date, errs := Tc.DateRepo.GetAll(&query.Query{})
	require.Empty(t, errs)
	require.Greater(t, len(date), 0)
}

func TestDateRepoGetById(t *testing.T) {
	id := Tf.Dates[0].ID
	date, errs := Tc.DateRepo.GetByID(id)
	require.Empty(t, errs)
	require.Equal(t, date.ID, id)
}

func TestDateRepoDelete(t *testing.T) {
	listMember := Tf.Dates[1]
	errs := Tc.DateRepo.Delete(listMember.ID)
	require.Empty(t, errs)
	_, errs = Tc.DateRepo.GetByID(listMember.ID)
	require.NotEmpty(t, errs)
}
