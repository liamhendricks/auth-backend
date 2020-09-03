package services

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEcryption(t *testing.T) {
	s := NewPasswordServiceBcrypt()
	password := "password"
	data, err := s.Hash([]byte(password))
	require.Nil(t, err)
	require.NotEqual(t, data, password)

	equal := s.Compare(string(data), password)
	require.Equal(t, equal, true)
}
