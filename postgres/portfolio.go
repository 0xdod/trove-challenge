package postgres

import (
	"context"

	"github.com/0xdod/trove"
)

type gormPortfolioService struct {
	*DB
}

func (gp *gormPortfolioService) FindByID(_ context.Context, _ int) (*trove.PortfolioPosition, error) {
	panic("not implemented") // TODO: Implement
}

func (gp *gormPortfolioService) FindByUser(ctx context.Context, userID int) ([]*trove.PortfolioPosition, error) {
	panic("not implemented") // TODO: Implement
}

func (gp *gormPortfolioService) Find(_ context.Context, _ trove.PositionFilter) ([]*trove.PortfolioPosition, error) {
	panic("not implemented") // TODO: Implement
}

func NewPortfolioService(db *DB) trove.PortfolioService {
	db.db.AutoMigrate(&trove.PortfolioPosition{})
	return &gormPortfolioService{db}
}
