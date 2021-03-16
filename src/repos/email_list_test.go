package repos

import (
	"testing"

	"github.com/68696c6c/goat/query"
	"github.com/liamhendricks/auth-backend/src/models"
	"github.com/stretchr/testify/require"
)

func TestEmailListRepoSave(t *testing.T) {
	listMember := models.MakeEmailList()
	errs := Tc.EmailListRepo.Save(&listMember)
	require.Empty(t, errs)
}

func TestEmailListRepoGetAll(t *testing.T) {
	el, errs := Tc.EmailListRepo.GetAll(&query.Query{})
	require.Empty(t, errs)
	require.Greater(t, len(el), 0)
}

func TestEmailListRepoGetById(t *testing.T) {
	id := Tf.List[0].ID
	el, errs := Tc.EmailListRepo.GetByID(id)
	require.Empty(t, errs)
	require.Equal(t, el.ID, id)
}

func TestEmailListRepoDelete(t *testing.T) {
	listMember := Tf.List[1]
	errs := Tc.EmailListRepo.Delete(listMember.ID)
	require.Empty(t, errs)
	_, errs = Tc.EmailListRepo.GetByID(listMember.ID)
	require.NotEmpty(t, errs)
}
