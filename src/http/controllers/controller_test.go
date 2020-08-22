package controllers

import (
	"testing"

	"github.com/liamhendricks/auth-backend/src/models"
	"github.com/stretchr/testify/require"
)

func TestRevoke(t *testing.T) {
	var c []*models.Course
	c1 := &models.Course{
		Name: "course1",
	}
	c2 := &models.Course{
		Name: "course2",
	}
	c = append(c, c1, c2)
	user := models.User{
		Courses: c,
	}

	require.Equal(t, user.Courses[0].Name, "course1")
	require.Equal(t, user.Courses[1].Name, "course2")
	require.Equal(t, len(user.Courses), 2)

	for k, v := range user.Courses {
		if v.Name == "course1" {
			user.Courses = revoke(k, user.Courses)
			break
		}
	}

	require.Equal(t, len(user.Courses), 1)
	require.Equal(t, user.Courses[0].Name, "course2")
}
