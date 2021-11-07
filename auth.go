package trove

import (
	"context"
	"crypto/sha256"
	"encoding/base32"
	"math/rand"
	"time"
)

const (
	ScopeActivation     = "activation"
	ScopeAuthentication = "authentication"
)

type Token struct {
	PlainText string    `json:"token,omitempty" gorm:"unique"`
	Hash      []byte    `json:"-" gorm:"unique"`
	UserID    int       `json:"-"`
	User      User      `json:"-"`
	Expiry    time.Time `json:"expiry,omitempty"`
	Scope     string    `json:"-"`
}

func NewAuthToken(userID int, ttl time.Duration) (*Token, error) {
	return generateToken(userID, ttl, ScopeAuthentication)
}

func generateToken(userID int, ttl time.Duration, scope string) (*Token, error) {
	t := Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}
	rand.Seed(time.Now().Unix())
	randBytes := make([]byte, 16)

	_, err := rand.Read(randBytes)

	if err != nil {
		return nil, err
	}

	t.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randBytes)

	hash := sha256.Sum256([]byte(t.PlainText))
	t.Hash = hash[:]

	return &t, nil
}

type AuthService interface {
	CreateToken(context.Context, *Token) error
}
