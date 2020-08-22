package repos

import (
	"testing"
	"time"

	"github.com/68696c6c/goat"
	"github.com/liamhendricks/auth-backend/src/models"
	"github.com/stretchr/testify/require"
)

func TestSessionRepo(t *testing.T) {
	id := goat.NewID()
	dur := 30 * time.Minute
	exp := time.Now().Add(dur)

	r := &models.Session{
		Token:      id,
		UserID:     Tf.Users[0].ID,
		Expiration: exp,
	}

	errs := Tc.SessionRepo.Save(r)
	require.Nil(t, errs)

	//should fail since this combo does not exist
	_, errs = Tc.SessionRepo.GetByTokenUser(r.Token, Tf.Users[1].ID)
	require.NotEmpty(t, errs)

	theSession, errs := Tc.SessionRepo.GetByTokenUser(r.Token, r.UserID)
	require.Equal(t, r.ID, theSession.ID)
	require.Equal(t, time.Now().Before(theSession.Expiration), true)

	errs = Tc.SessionRepo.Delete(theSession.ID)
	require.Empty(t, errs)
}
