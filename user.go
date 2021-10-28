package trove

import (
	"context"
	"time"
)

type User struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
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
