package postgres

import (
	"context"

	"github.com/0xdod/trove"
)

type gormLoanService struct {
	*DB
}

func (g *gormLoanService) CreateLoan(ctx context.Context, loan *trove.Loan) error {
	db := g.db.WithContext(ctx)
	return db.Create(loan).Error
}

func (g *gormLoanService) GetLoansByUser(ctx context.Context, userID int) ([]*trove.Loan, error) {
	db := g.db.WithContext(ctx)
	loans := []*trove.Loan{}
	if err := db.Model(&trove.Loan{}).Where("user_id = ?", userID).Find(loans).Error; err != nil {
		return nil, err
	}
	return loans, nil
}

func NewLoanService(db *DB) trove.LoanService {
	return &gormLoanService{db}
}
