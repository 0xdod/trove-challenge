package trove

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int       `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserPatch struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}

type FilterOption struct {
	Limit  *int `json:"limit"`
	Offset *int `json:"offset"`
}

type UserFilter struct {
	ID           *int    `json:"id"`
	Name         *string `json:"name"`
	Email        *string `json:"email"`
	FilterOption `json:"filter_option"`
}

type UserService interface {
	Create(context.Context, *User) error
	FindUserByID(context.Context, int) (*User, error)
	FindUsers(context.Context, UserFilter) ([]*User, error)
	UpdateUser(context.Context, int, UserPatch) (*User, error)
	DeleteUser(context.Context, int) error
}

func (u *User) SetPassword(password string) error {
	pBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	u.PasswordHash = string(pBytes)

	return nil
}

func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}
