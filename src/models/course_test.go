package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConvertStringToCourseType(t *testing.T) {
	goodType := "Free"
	badType := "bad"

	goodUT, err := CourseTypeFromString(goodType)
	require.Nil(t, err)
	require.Equal(t, FreeCourse, goodUT)

	_, err = CourseTypeFromString(badType)
	require.NotNil(t, err)
}
