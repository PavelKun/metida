package cryptography

import (
	"fmt"
	"math/rand"
	"time"
)

type RefreshToken interface {
	CreateRefreshToken() (string, error)
}

type TokenRefresh struct{}

func NewRefreshToken() *TokenRefresh {

	return &TokenRefresh{}
}

func (o *TokenRefresh) CreateRefreshToken() (string, error) {
	b := make([]byte, 32)
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	_, err := r.Read(b)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
