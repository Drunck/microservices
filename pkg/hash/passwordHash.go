package hash

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type bcryptHasher struct {
	cost int
}

func newBcryptHasher() *bcryptHasher {
	return &bcryptHasher{}
}

func NewBCryptHasher(cost int) *bcryptHasher {
	return &bcryptHasher{cost: cost}
}

func (h *bcryptHasher) Hash(password string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}

	hash := fmt.Sprintf("%s", hashBytes)
	return hash, nil
}
