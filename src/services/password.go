package services

import (
	"golang.org/x/crypto/bcrypt"
)

type PasswordService interface {
	Hash(data []byte) ([]byte, error)
	Compare(data, plain string) bool
}

type PasswordServiceBcrypt struct {
	cost int
}

func NewPasswordServiceBcrypt() PasswordServiceBcrypt {
	return PasswordServiceBcrypt{
		cost: bcrypt.MinCost,
	}
}

func (s PasswordServiceBcrypt) Hash(data []byte) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword(data, s.cost)
	if err != nil {
		return nil, err
	}

	return hash, nil
}

func (s PasswordServiceBcrypt) Compare(data, plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(data), []byte(plain))
	if err != nil {
		return false
	}

	return true
}
