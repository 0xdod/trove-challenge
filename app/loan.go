package app

import (
	"math"
	"net/http"
)

func (s *Server) processLoan() http.HandlerFunc {
	type loanRequest struct {
		Amount   float64 `json:"amount,omitempty" validate:"required"`
		Duration int     `json:"duration,omitempty" validate:"required,min=6,max=12"`
	}

	return (func(w http.ResponseWriter, r *http.Request) {
		// max 60% of value
		// payment period 6 - 12 months
		// calculate payment per month depending on duration
		// interest rate is 15%

		loanReq := &loanRequest{}

		if err := s.readJSON(r.Body, loanReq); err != nil {
			s.writeJSON(w, http.StatusUnprocessableEntity, RM{"error", err.Error(), nil})
			return
		}

		if err := validate.Struct(loanReq); err != nil {
			s.writeJSON(w, http.StatusBadRequest, RM{"error", err.Error(), nil})
			return
		}

		user := UserFromContext(r.Context())

		if user == nil {
			s.writeJSON(w, http.StatusInternalServerError, RM{"error", "internal error", nil})
			return
		}

		pValue, err := s.PortfolioService.GetPortfolioValue(r.Context(), user.ID)

		if err != nil {
			s.writeJSON(w, http.StatusInternalServerError, RM{"error", "internal error", nil})
			return
		}

		if !isEligibileForLoan(pValue, loanReq.Amount) {
			s.writeJSON(w, http.StatusOK, RM{
				Status:  "fail",
				Message: "loan request declined, you are not eligible for this amount.",
			})
			return
		}

		proratedPayment := calcProratedPayment(loanReq.Amount, loanReq.Duration)

		s.writeJSON(w, http.StatusOK, RM{
			Status:  "success",
			Message: "loan request approved",
			Data: M{
				"amount":                     loanReq.Amount,
				"repayment_period_in_months": loanReq.Duration,
				"total_repayment":            proratedPayment * float64(loanReq.Duration),
				"monthly_repayment":          proratedPayment,
			},
		})
	})
}

func isEligibileForLoan(portfolioValue, loanAmount float64) bool {
	maximumLoan := (60.00 / 100.00) * portfolioValue
	return loanAmount < maximumLoan
}

func calculateRepayment(amount, rate float64, duration int) float64 {
	// 1 year == 12 months
	// x years == duration months
	durationInYears := float64(duration) / 12
	rateInPercent := rate / 100
	// compound interest compounding once per month
	return amount * math.Pow(1+(rateInPercent/12), 12*durationInYears)
}

func calcProratedPayment(amount float64, period int) float64 {
	interestRate := 15
	totalRepayment := calculateRepayment(amount, float64(interestRate), period)

	return totalRepayment / float64(period)
}
