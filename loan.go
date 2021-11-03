package trove

import (
	"context"
	"time"
)

type Loan struct {
	ID           int       `json:"id,omitempty" gorm:"primaryKey"`
	Amount       float64   `json:"amount,omitempty"`
	InterestRate float64   `json:"interest_rate,omitempty"`
	UserID       int       `json:"user_id,omitempty"`
	User         User      `json:"user,omitempty"`
	DueDate      time.Time `json:"due_date,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}

type LoanService interface {
	CreateLoan(context.Context, *Loan) error
	GetLoansByUser(context.Context, int) (*Loan, error)
}
