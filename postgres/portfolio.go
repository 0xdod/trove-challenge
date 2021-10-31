package postgres

import (
	"context"
	"log"

	"github.com/0xdod/trove"
)

type gormPortfolioService struct {
	*DB
}

func (gp *gormPortfolioService) FindByID(_ context.Context, _ int) (*trove.PortfolioPosition, error) {
	panic("not implemented") // TODO: Implement
}

func (gp *gormPortfolioService) FindByUser(ctx context.Context, userID int) ([]*trove.PortfolioPosition, error) {
	portfolio := []*trove.PortfolioPosition{}
	db := gp.db.WithContext(ctx)
	err := db.Model(&trove.PortfolioPosition{}).Where("user_id = ?", userID).Find(&portfolio).Error

	if err != nil {
		return nil, err
	}

	return portfolio, nil
}

func (gp *gormPortfolioService) Find(_ context.Context, _ trove.PositionFilter) ([]*trove.PortfolioPosition, error) {
	panic("not implemented") // TODO: Implement
}

func (gp *gormPortfolioService) GetPortfolioValue(ctx context.Context, userID int) (float64, error) {
	type result struct {
		UserID     int
		TotalValue float64
	}

	db := gp.db.WithContext(ctx)
	sql := `SELECT user_id, SUM(total_quantity*price_per_share) as total_value
	        FROM portfolio_positions WHERE user_id = ?
			GROUP BY user_id`
	res := &result{}

	if err := db.Raw(sql, userID).Scan(res).Error; err != nil {
		log.Println(db.Statement.SQL.String())
		return 0, err
	}

	return res.TotalValue, nil
}

func NewPortfolioService(db *DB) trove.PortfolioService {
	db.db.AutoMigrate(&trove.PortfolioPosition{})
	return &gormPortfolioService{db}
}
