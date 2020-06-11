package repos

import (
	"testing"

	"github.com/liamhendricks/auth-backend/src/models"
	"github.com/stretchr/testify/require"
)

func TestUsersRepoSave(t *testing.T) {
	u := &models.User{
		Name: "TestUser",
	}
	errs := tc.UserRepo.Save(u)
	require.Nil(t, errs)
	theUser, errs := tc.UserRepo.GetByID(u.ID, false)
	require.Equal(t, u.ID, theUser.ID)
}

func TestUsersRepoGetByID(t *testing.T) {
	id := tf.Users[0].ID
	user, errs := tc.UserRepo.GetByID(id, false)
	require.Nil(t, errs)
	require.Equal(t, id, user.ID)
}
