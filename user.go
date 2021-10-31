package trove

import (
	"context"
	"time"

	"github.com/mitchellh/mapstructure"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	FirstName    string    `json:"first_name" gorm:"not null;size:50"`
	LastName     string    `json:"last_name" gorm:"not null;size:50"`
	Email        string    `json:"email" gorm:"unique;size:255"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserPatch struct {
	FirstName *string      `json:"first_name" mapstructure:"first_name,omitempty"`
	LastName  *string      `json:"last_name" mapstructure:"last_name,omitempty"`
	Email     *string      `json:"email" mapstructure:"email,omitempty" validate:"email"` // TODO remove and change to a secure flow
	Password  *string      `json:"password" mapstructure:"password,omitempty"`            // TODO: remove and change to reset password flow
	options   FilterOption `json:"-"`
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

func (u *User) ApplyPatch(p UserPatch) {

	if p.FirstName != nil {
		u.FirstName = *p.FirstName
	}

	if p.LastName != nil {
		u.LastName = *p.LastName
	}

	if p.Email != nil {
		u.Email = *p.Email
	}

	if p.Password != nil {
		u.SetPassword(*p.Password)
	}
}

func (p UserPatch) ToMap() map[string]interface{} {
	out := make(map[string]interface{})
	if err := mapstructure.Decode(p, &out); err != nil {
		return nil
	}
	return out
}
