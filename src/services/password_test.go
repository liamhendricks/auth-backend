package services

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEcryption(t *testing.T) {
	s := NewPasswordServiceAES(PasswordConfig{
		PassPhrase: "test_phrase",
	})
	password := "password"
	data, err := s.Encrypt([]byte(password))
	require.Nil(t, err)
	require.NotEqual(t, data, password)

	plaintext, err := s.Decrypt(data)
	require.Nil(t, err)
	require.Equal(t, password, string(plaintext))
}

func TestEcryptionBadPhrase(t *testing.T) {
	s := NewPasswordServiceAES(PasswordConfig{
		PassPhrase: "first_phrase",
	})
	password := "password"
	data, err := s.Encrypt([]byte(password))
	require.Nil(t, err)
	require.NotEqual(t, data, password)

	s.passPhrase = "second_phrase"
	_, err = s.Decrypt(data)
	require.NotNil(t, err)
}
