package trove

import (
	"context"
	"math"
	"time"
)

type Loan struct {
	ID           int       `json:"id,omitempty" gorm:"primaryKey"`
	Amount       float64   `json:"amount"`
	PaidBack     float64   `json:"paid_back"`
	InterestRate float64   `json:"interest_rate"`
	UserID       int       `json:"-"`
	User         User      `json:"-"`
	Duration     int       `json:"duration"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}

func (l *Loan) PaymentDue() time.Time {
	if l.PaymentComplete() {
		return time.Time{}
	}

	return time.Now().AddDate(0, 1, 0)
}

func (l *Loan) PaymentComplete() bool {
	return l.Balance() <= 0
}

func (l *Loan) Balance() float64 {
	return l.Amount - l.PaidBack
}

func (l *Loan) TotalRepayment() float64 {
	// 1 year == 12 months
	// x years == duration months
	durationInYears := float64(l.Duration) / 12
	rateInPercent := l.InterestRate / 100
	// compound interest compounding once per month
	return l.Amount * math.Pow(1+(rateInPercent), durationInYears)
}

func (l *Loan) ProratedPayment() float64 {
	return l.TotalRepayment() / float64(l.Duration)
}

type LoanService interface {
	CreateLoan(context.Context, *Loan) error
	GetLoansByUser(context.Context, int) ([]*Loan, error)
}
