package trove

import (
	"context"
	"math"
	"time"
)

type Loan struct {
	ID                int       `json:"id,omitempty" gorm:"primaryKey"`
	Amount            float64   `json:"amount,omitempty"`
	PaidBack          float64   `json:"paid_back,omitempty"`
	InterestRate      float64   `json:"interest_rate,omitempty"`
	UserID            int       `json:"-,omitempty"`
	User              User      `json:"-,omitempty"`
	InstallmentNumber int       `json:"-,omitempty"`
	Duration          int       `json:"duration,omitempty"`
	NextPaymentAmount float64   `json:"next_payment_amount,omitempty"`
	DueDate           time.Time `json:"due_date,omitempty"`
	CreatedAt         time.Time `json:"created_at,omitempty"`
	UpdatedAt         time.Time `json:"updated_at,omitempty"`
}

func (l *Loan) PaymentDue() time.Time {
	if l.PaymentComplete() {
		return time.Time{}
	}
	l.InstallmentNumber++
	return l.CreatedAt.AddDate(0, l.InstallmentNumber, 0)
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
