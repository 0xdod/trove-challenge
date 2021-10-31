package trove

import "context"

type PortfolioPosition struct {
	ID            int     `json:"id,omitempty"`
	UserID        int     `json:"user_id,omitempty"`
	User          *User   `json:"user,omitempty"`
	Symbol        string  `json:"symbol,omitempty"`
	TotalQuantity float64 `json:"total_quantity,omitempty"`
	EquityValue   float64 `json:"equity_value,omitempty"`
	PricePerShare float64 `json:"price_per_share,omitempty"`
}

type Portfolio []*PortfolioPosition

type PositionFilter struct{}

type PortfolioService interface {
	FindByID(context.Context, int) (*PortfolioPosition, error)
	FindByUser(context.Context, int) ([]*PortfolioPosition, error)
	Find(context.Context, PositionFilter) ([]*PortfolioPosition, error)
	GetPortfolioValue(context.Context, int) (float64, error)
}
