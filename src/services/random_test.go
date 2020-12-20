package services

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomString(t *testing.T) {
	s1 := RandomString(10)
	require.Len(t, s1, 10)
}
