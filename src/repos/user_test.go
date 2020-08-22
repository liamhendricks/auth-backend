package repos

import (
	"fmt"
	"testing"
	"time"

	"github.com/68696c6c/goat"
	"github.com/68696c6c/goat/query"
	"github.com/icrowley/fake"
	"github.com/liamhendricks/auth-backend/src/models"
	"github.com/stretchr/testify/require"
)

func TestUsersRepoSave(t *testing.T) {
	u := &models.User{
		Name: "TestUser",
	}
	errs := Tc.UserRepo.Save(u)
	require.Nil(t, errs)
	theUser, errs := Tc.UserRepo.GetByID(u.ID, false)
	require.Equal(t, u.ID, theUser.ID)
}

func TestUsersRepoUpdate(t *testing.T) {
	u := Tf.Users[1]
	u.Name = "NewName"
	u.Email = "NewEmail"
	errs := Tc.UserRepo.Save(&u)
	require.Nil(t, errs)
	theUser, errs := Tc.UserRepo.GetByID(u.ID, false)
	require.Equal(t, u.ID, theUser.ID)
	require.Equal(t, theUser.Email, "NewEmail")
	require.Equal(t, theUser.Name, "NewName")
	require.Equal(t, u.Password, theUser.Password)
}

func TestUsersRepoGetByID(t *testing.T) {
	id := Tf.Users[0].ID
	user, errs := Tc.UserRepo.GetByID(id, false)
	require.Nil(t, errs)
	require.Equal(t, id, user.ID)
}

func TestUsersRepoUniqueEmail(t *testing.T) {
	email := Tf.Users[0].Email
	newUser := &models.User{
		Name:  fake.FullName(),
		Email: email,
	}
	errs := Tc.UserRepo.Save(newUser)
	require.Equal(t, errs[0].Error(), fmt.Sprintf("Error 1062: Duplicate entry '%s' for key 'email'", email))
}

func TestUsersRepoGetAll(t *testing.T) {
	u, errs := Tc.UserRepo.GetAll(&query.Query{})
	require.Empty(t, errs)
	require.GreaterOrEqual(t, len(u), 1)
}

func TestUsersRepoDelete(t *testing.T) {
	id := Tf.Users[2].ID
	errs := Tc.UserRepo.Delete(id)
	require.Empty(t, errs)
	_, errs = Tc.UserRepo.GetByID(id, false)
	require.NotNil(t, errs)
}

func TestUsersGetAllRelations(t *testing.T) {
	id := Tf.Users[0].ID
	u, errs := Tc.UserRepo.GetByID(id, false)
	require.Nil(t, errs)

	s := &models.Session{
		Token:      goat.NewID(),
		Expiration: time.Now().Add(30 * time.Minute),
	}
	r := &models.Reset{
		Token:      goat.NewID(),
		Expiration: time.Now().Add(30 * time.Minute),
	}

	u.Session = s
	u.Reset = r

	errs = Tc.UserRepo.Save(&u)
	require.Nil(t, errs)
	u, errs = Tc.UserRepo.GetByID(id, true)
	require.Nil(t, errs)
	require.NotNil(t, u.Courses)
	require.NotNil(t, u.Reset)
	require.NotNil(t, u.Session)
}
