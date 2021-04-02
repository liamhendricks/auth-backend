package repos

import (
	"testing"
	"time"

	"github.com/68696c6c/goat"
	"github.com/liamhendricks/auth-backend/src/models"
	"github.com/stretchr/testify/require"
)

func TestResetsRepo(t *testing.T) {
	id := goat.NewID()
	dur := 30 * time.Minute
	exp := time.Now().Add(dur)

	r := &models.Reset{
		Token:      id,
		UserID:     Tf.Users[0].ID,
		Expiration: exp,
	}

	errs := Tc.ResetRepo.Save(r)
	require.Nil(t, errs)

	//should fail since this combo does not exist
	_, errs = Tc.ResetRepo.GetByToken(r.Token)
	require.NotEmpty(t, errs)

	theReset, errs := Tc.ResetRepo.GetByTokenUser(r.Token, r.UserID)
	require.Equal(t, r.ID, theReset.ID)
	require.Equal(t, time.Now().Before(theReset.Expiration), true)

	errs = Tc.ResetRepo.Delete(theReset.ID)
	require.Empty(t, errs)
}
