package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/0xdod/trove"
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

		newLoan := &trove.Loan{
			Amount:       loanReq.Amount,
			Duration:     loanReq.Duration,
			UserID:       user.ID,
			InterestRate: 15,
		}

		isEligible, err := s.UserIsEligibileForLoan(user, newLoan)

		if err != nil {
			s.serverErrorResponse(w, err)
			return
		}

		if !isEligible {
			s.writeJSON(w, http.StatusOK, RM{
				Status:  "fail",
				Message: "loan request declined, you are not eligible for this amount.",
			})
			return
		}

		if err := s.LoanService.CreateLoan(r.Context(), newLoan); err != nil {
			s.serverErrorResponse(w, err)
			return
		}

		_ = s.writeJSON(w, http.StatusOK, RM{
			Status:  "success",
			Message: "loan request approved",
			Data: M{
				"amount":             loanReq.Amount,
				"repayment_duration": fmt.Sprintf("%d months", loanReq.Duration),
				"next_repayment":     newLoan.PaymentDue().Format("Mon Jan 2 15:04:05 2006"),
				"monthly_repayment":  newLoan.ProratedPayment(),
			},
		})
	})
}

func (s *Server) listLoans(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())

	if user == nil {
		s.serverErrorResponse(w, errors.New("user not authenticated"))
		return
	}

	loans, err := s.LoanService.GetLoansByUser(r.Context(), user.ID)

	if err != nil {
		s.serverErrorResponse(w, err)
		return
	}

	s.writeJSON(w, http.StatusOK, RM{"success", "loans retrieved", loans})
}

func (s *Server) UserIsEligibileForLoan(user *trove.User, loan *trove.Loan) (bool, error) {
	// check previous loans
	ctx := context.Background()
	loans, err := s.LoanService.GetLoansByUser(ctx, user.ID)

	if err != nil {
		return false, err
	}

	val, err := s.PortfolioService.GetPortfolioValue(ctx, user.ID)

	if err != nil {
		return false, err
	}

	totalLoanBal := 0.00

	for _, loan := range loans {
		totalLoanBal += loan.Balance()
	}

	loanLimit := (60.0 / 100.0) * val

	// total loan balance up to 60% of portfolio value?
	return loanLimit > totalLoanBal, nil
}
