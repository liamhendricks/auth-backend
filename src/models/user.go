package models

import (
	"errors"

	"github.com/68696c6c/goat"
	"github.com/icrowley/fake"
)

type User struct {
	goat.Model
	Name     string    `json:"name" binding:"required"`
	Email    string    `json:"email" binding:"required"`
	Password string    `json:"-" binding:"required"`
	UserType UserType  `json:"user_type" binding:"required"`
	Courses  []*Course `gorm:"many2many:user_courses;"`
	Session  *Session
	Reset    *Reset
	Status   UserStatus
}

type UserType string
type UserStatus string

const AdminUser UserType = "Admin"
const FreeUser UserType = "Free"
const PaidUser UserType = "Paid"

const UserOk = "OK"
const UserMissingCourse = "Missing Course"
const UserUnpaid = "Unpaid"

func MakeUser() User {
	return User{
		Name:     fake.FullName(),
		Email:    fake.EmailAddress(),
		Password: fake.Password(5, 16, true, true, true),
		UserType: FreeUser,
	}
}

func UserTypeFromString(t string) (UserType, error) {
	if t == string(AdminUser) {
		return AdminUser, nil
	}
	if t == string(FreeUser) {
		return FreeUser, nil
	}
	if t == string(PaidUser) {
		return PaidUser, nil
	}

	return UserType(""), errors.New("not a valid user type")
}

func (t UserType) IsGreaterThanEqTo(v UserType) bool {
	var pass []UserType
	switch t {
	case FreeUser:
		pass = []UserType{FreeUser}
	case PaidUser:
		pass = []UserType{FreeUser, PaidUser}
	case AdminUser:
		pass = []UserType{FreeUser, PaidUser, AdminUser}
	}

	return userTypeInSlice(v, pass)
}

func userTypeInSlice(a UserType, list []UserType) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}

	return false
}
