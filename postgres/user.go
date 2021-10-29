package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/0xdod/trove"
	"gorm.io/gorm"
)

// GormUser is an implementation of the UserService interface using the GORM library with postgres
type gormUser struct {
	*DB
}

func (g *gormUser) Create(ctx context.Context, u *trove.User) error {
	u.Email = strings.ToLower(u.Email)
	err := createUser(ctx, g.DB.db, u)

	if err != nil {
		return fmt.Errorf("cannot create user: %v", err)
	}

	return nil
}

func (g *gormUser) FindUserByID(_ context.Context, _ int) (*trove.User, error) {
	panic("not implemented") // TODO: Implement
}

func (g *gormUser) FindUsers(_ context.Context, _ trove.UserFilter) ([]*trove.User, error) {
	panic("not implemented") // TODO: Implement
}

func (g *gormUser) UpdateUser(_ context.Context, _ int, _ trove.UserPatch) (*trove.User, error) {
	panic("not implemented") // TODO: Implement
}

func (g *gormUser) DeleteUser(_ context.Context, _ int) error {
	panic("not implemented") // TODO: Implement
}

func NewUserService(db *DB) trove.UserService {
	db.db.AutoMigrate(&trove.User{})
	return &gormUser{db}
}

func createUser(ctx context.Context, db *gorm.DB, u *trove.User) error {
	return db.WithContext(ctx).Create(u).Error
}
