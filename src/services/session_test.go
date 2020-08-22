package services

import (
	"testing"

	"github.com/liamhendricks/auth-backend/src/models"
	"github.com/liamhendricks/auth-backend/src/repos"
	"github.com/stretchr/testify/require"
)

func TestSessionServiceStart(t *testing.T) {
	user := models.MakeUser()
	dr := repos.NewSessionRepoDummy(&user)
	ss := NewSessionServiceDB(dr)
	err := ss.Start(&user)
	require.Nil(t, err)
	require.NotNil(t, user.Session.Token)
}

func TestSessionServiceValid(t *testing.T) {
	user := models.MakeUser()
	dr := repos.NewSessionRepoDummy(&user)
	ss := NewSessionServiceDB(dr)
	err := ss.Start(&user)
	require.Nil(t, err)
	require.NotNil(t, user.Session.Token)

	v := ss.Valid(&user, user.Session.Token)
	require.True(t, v)
}

func TestSessionServiceStop(t *testing.T) {
	user := models.MakeUser()
	dr := repos.NewSessionRepoDummy(&user)
	ss := NewSessionServiceDB(dr)
	err := ss.Start(&user)
	require.Nil(t, err)
	require.NotNil(t, user.Session.Token)

	err = ss.Stop(&user)
	require.Nil(t, err)
}
