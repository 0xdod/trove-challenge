package postgres

import (
	"context"
	"fmt"

	"github.com/0xdod/trove"
)

type gormAuthService struct {
	*DB
}

func (g *gormAuthService) CreateToken(ctx context.Context, t *trove.Token) error {
	db := g.db.WithContext(ctx)

	if err := db.Create(t).Error; err != nil {
		return fmt.Errorf("cannot create user: %v", err)
	}

	return nil
}

func NewAuthService(db *DB) trove.AuthService {
	db.db.AutoMigrate(&trove.Token{})
	return &gormAuthService{db}
}
