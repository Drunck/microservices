package token

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

type TokenManager interface {
	NewToken(userId int64, ttl time.Duration) (string, error)
}

type Manager struct {
	signingKey string
}

func NewManager(signingKey string) (*Manager, error) {
	if signingKey == "" {
		return nil, errors.New("empty signing key")
	}
	return &Manager{signingKey: signingKey}, nil
}

func (m *Manager) NewToken(userId int64, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(ttl).Unix(),
		Subject:   strconv.FormatInt(userId, 10),
	})
	return token.SignedString([]byte(m.signingKey))
}
