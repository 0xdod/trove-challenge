package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/0xdod/trove"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (g *gormUser) FindUserByID(ctx context.Context, id int) (*trove.User, error) {
	user := &trove.User{}
	db := g.db.WithContext(ctx)
	err := db.Model(user).Where("id = ?", id).First(user).Error

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (g *gormUser) FindUserByEmail(ctx context.Context, email string) (*trove.User, error) {
	user := &trove.User{}
	db := g.db.WithContext(ctx)
	err := db.Model(user).Where("email = ?", email).First(user).Error

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (g *gormUser) FindUserByToken(ctx context.Context, tokenStr string) (*trove.User, error) {
	token := &trove.Token{}
	db := g.db.WithContext(ctx)
	err := db.Model(&trove.Token{}).Where("plain_text = ?", tokenStr).First(token).Error

	if err != nil {
		return nil, err
	}
	return g.FindUserByID(ctx, token.UserID)
}

func (g *gormUser) FindUsers(_ context.Context, _ trove.UserFilter) ([]*trove.User, error) {
	panic("not implemented") // TODO: Implement
}

func (g *gormUser) UpdateUser(ctx context.Context, id int, up trove.UserPatch) (*trove.User, error) {
	user := &trove.User{}

	if err := g.db.WithContext(ctx).Model(user).Clauses(clause.Returning{}).Where("id = ?", id).Updates(up.ToMap()).Error; err != nil {
		return nil, err
	}

	return user, nil
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
