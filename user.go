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
	FirstName       *string `json:"first_name,omitempty" mapstructure:"first_name,omitempty"`
	LastName        *string `json:"last_name,omitempty" mapstructure:"last_name,omitempty"`
	Email           *string `json:"email,omitempty" mapstructure:"email,omitempty" validate:"email"` // TODO remove and change to a secure flow
	OldPassword     *string `json:"old_password,omitempty" mapstructure:"-"`
	NewPassword     *string `json:"new_password,omitempty" mapstructure:"-,omitempty"`
	NewPasswordHash *string `json:"-" mapstructure:"password_hash,omitempty"`
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
	FindUserByEmail(context.Context, string) (*User, error)
	FindUserByToken(context.Context, string) (*User, error)
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

	if p.OldPassword != nil {
		if u.VerifyPassword(*p.OldPassword) {
			u.SetPassword(*p.NewPassword)
		}
	}
}

func (p UserPatch) ToMap() map[string]interface{} {
	out := make(map[string]interface{})
	if err := mapstructure.Decode(p, &out); err != nil {
		return nil
	}
	return out
}

func (u *User) IsAnonymous() bool {
	return *u == User{}
}
