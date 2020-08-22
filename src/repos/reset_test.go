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
		UserID:     tf.Users[0].ID,
		Expiration: exp,
	}

	errs := tc.ResetRepo.Save(r)
	require.Nil(t, errs)

	_, errs = tc.ResetRepo.GetByTokenUser(r.Token, tf.Users[1].ID)
	require.NotEmpty(t, errs)

	theReset, errs := tc.ResetRepo.GetByTokenUser(r.Token, r.UserID)
	require.Equal(t, r.ID, theReset.ID)
	require.Equal(t, time.Now().Before(theReset.Expiration), true)

	errs = tc.ResetRepo.Delete(theReset.ID)
	require.Empty(t, errs)
}
