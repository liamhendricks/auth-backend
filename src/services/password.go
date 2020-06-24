package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
)

type PasswordService interface {
	Encrypt(data []byte) ([]byte, error)
	Decrypt(data []byte) ([]byte, error)
}

type PasswordServiceAES struct {
	passPhrase string
}

type PasswordConfig struct {
	PassPhrase string
}

func NewPasswordServiceAES(c PasswordConfig) PasswordServiceAES {
	return PasswordServiceAES{
		passPhrase: c.PassPhrase,
	}
}

func getHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (s PasswordServiceAES) Encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher([]byte(getHash(s.passPhrase)))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func (s PasswordServiceAES) Decrypt(data []byte) ([]byte, error) {
	key := []byte(getHash(s.passPhrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
