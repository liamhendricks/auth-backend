package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserTypeFromString(t *testing.T) {
	realUserType := "Admin"
	admin, err := UserTypeFromString(realUserType)
	require.Equal(t, admin, AdminUser)
	require.Nil(t, err)

	fakeUserType := "TheAdmin"
	fakeAdmin, err := UserTypeFromString(fakeUserType)
	require.NotNil(t, err)
	require.Equal(t, fakeAdmin, UserType(""))
}

func TestUserGreaterThanEqTo(t *testing.T) {
	u := MakeUser()
	ok := u.UserType.IsGreaterThanEqTo(FreeUser)
	require.Equal(t, true, ok)

	ok = u.UserType.IsGreaterThanEqTo(PaidUser)
	require.Equal(t, false, ok)

	ok = u.UserType.IsGreaterThanEqTo(AdminUser)
	require.Equal(t, false, ok)

	u.UserType = AdminUser
	ok = u.UserType.IsGreaterThanEqTo(PaidUser)
	require.Equal(t, true, ok)
}

func TestUserTypeInSlice(t *testing.T) {
	list := []UserType{AdminUser, PaidUser, FreeUser}
	ok := userTypeInSlice(AdminUser, list)
	require.Equal(t, true, ok)

	testUser := UserType("TEST")

	ok = userTypeInSlice(testUser, list)
	require.Equal(t, false, ok)
}
