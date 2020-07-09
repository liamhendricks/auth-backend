package repos

import (
	"fmt"
	"testing"

	"github.com/68696c6c/goat/query"
	"github.com/icrowley/fake"
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

func TestUsersRepoUpdate(t *testing.T) {
	u := tf.Users[1]
	u.Name = "NewName"
	u.Email = "NewEmail"
	errs := tc.UserRepo.Save(&u)
	require.Nil(t, errs)
	theUser, errs := tc.UserRepo.GetByID(u.ID, false)
	require.Equal(t, u.ID, theUser.ID)
	require.Equal(t, theUser.Email, "NewEmail")
	require.Equal(t, theUser.Name, "NewName")
	require.Equal(t, u.Password, theUser.Password)
}

func TestUsersRepoGetByID(t *testing.T) {
	id := tf.Users[0].ID
	user, errs := tc.UserRepo.GetByID(id, false)
	require.Nil(t, errs)
	require.Equal(t, id, user.ID)
}

func TestUsersRepoUniqueEmail(t *testing.T) {
	email := tf.Users[0].Email
	newUser := &models.User{
		Name:  fake.FullName(),
		Email: email,
	}
	errs := tc.UserRepo.Save(newUser)
	require.Equal(t, errs[0].Error(), fmt.Sprintf("Error 1062: Duplicate entry '%s' for key 'email'", email))
}

func TestUsersRepoGetAll(t *testing.T) {
	u, errs := tc.UserRepo.GetAll(&query.Query{})
	require.Empty(t, errs)
	require.GreaterOrEqual(t, len(u), 0)
}

func TestUsersRepoDelete(t *testing.T) {
	id := tf.Users[2].ID
	errs := tc.UserRepo.Delete(id)
	require.Empty(t, errs)
	_, errs = tc.UserRepo.GetByID(id, false)
	require.NotNil(t, errs)
}
