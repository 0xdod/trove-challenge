package app

import (
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/0xdod/trove"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

var validate = validator.New()

type userSignupRequest struct {
	FirstName string `json:"first_name" validate:"required,min=2"`
	LastName  string `json:"last_name" validate:"required,min=2"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6,max=60"`
}

func (s *Server) registerUser(w http.ResponseWriter, r *http.Request) {

	signupReq := &userSignupRequest{}

	if err := s.readJSON(r.Body, &signupReq); err != nil {
		// return error message
		s.writeJSON(w, http.StatusUnprocessableEntity, RM{"error", "request error: " + err.Error(), nil})
		return
	}

	// do some request validation
	if err := validate.Struct(signupReq); err != nil {
		// return error message
		s.writeJSON(w, http.StatusBadRequest, RM{"error", "validation error: " + err.Error(), nil})
		return
	}

	user := &trove.User{
		FirstName: signupReq.FirstName,
		LastName:  signupReq.LastName,
		Email:     signupReq.Email,
	}

	if err := user.SetPassword(signupReq.Password); err != nil {
		// error
		s.writeJSON(w, http.StatusInternalServerError, RM{"error", "internal error: " + err.Error(), nil})
		return
	}

	// do business logic thingy, in this case signup
	if err := s.UserService.Create(r.Context(), user); err != nil {
		// do some error thing.
		s.writeJSON(w, http.StatusInternalServerError, RM{"error", "internal error: " + err.Error(), nil})
		return
	}

	// return response
	if err := s.writeJSON(w, http.StatusCreated, RM{"success", "created user account", user}); err != nil {
		log.Println(err)
	}
}

func (s *Server) updateUser(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)
	patch := trove.UserPatch{}

	if err := s.readJSON(r.Body, &patch); err != nil {
		s.writeJSON(w, http.StatusUnprocessableEntity, RM{"error", "request error: " + err.Error(), nil})
		return
	}

	user, err := s.UserService.UpdateUser(r.Context(), id, patch)

	if err != nil {
		s.writeJSON(w, http.StatusInternalServerError, RM{"error", "internal error, update failed", nil})
		return
	}

	if err := s.writeJSON(w, http.StatusOK, RM{"success", "update successful", user}); err != nil {
		s.writeJSON(w, http.StatusInternalServerError, RM{"error", "internal error", nil})
		log.Printf("json error %v", err)
		return
	}
}

func (s *Server) getPortfolio(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(mux.Vars(r)["id"])

	portfolio, err := s.PortfolioService.FindByUser(r.Context(), userID)

	if err != nil {
		s.writeJSON(w, http.StatusInternalServerError, RM{"error", "internal error", nil})
		return
	}

	if err := s.writeJSON(w, http.StatusOK, RM{"success", "portfolio positions retrieved", portfolio}); err != nil {
		s.writeJSON(w, http.StatusInternalServerError, RM{"error", "internal error", nil})
		log.Printf("json error %v", err)
		return
	}
}

func (s *Server) getPortfolioValue(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(mux.Vars(r)["id"])

	value, err := s.PortfolioService.GetPortfolioValue(r.Context(), userID)

	if err != nil {
		s.writeJSON(w, http.StatusInternalServerError, RM{"error", "internal error", nil})
		return
	}

	err = s.writeJSON(w, http.StatusOK, RM{"success", "portfolio value retrieved", M{"portfolio_value": value}})

	if err != nil {
		s.writeJSON(w, http.StatusInternalServerError, RM{"error", "internal error", nil})
		log.Printf("json error %v", err)
		return
	}
}

func (s *Server) processLoan(w http.ResponseWriter, r *http.Request) {
	// max 60% of value
	// payment period 6 - 12 months
	// calculate payment per month depending on duration
	// interest rate is 15%
	type loanRequest struct {
		Amount   float64 `json:"amount,omitempty" validate:"required"`
		Duration int     `json:"duration,omitempty" validate:"required,min=6,max=12"`
	}

	loanReq := &loanRequest{}

	if err := s.readJSON(r.Body, loanReq); err != nil {
		s.writeJSON(w, http.StatusUnprocessableEntity, RM{"error", err.Error(), nil})
		return
	}

	if err := validate.Struct(loanReq); err != nil {
		s.writeJSON(w, http.StatusBadRequest, RM{"error", err.Error(), nil})
		return
	}

	user, err := s.UserService.FindUserByID(r.Context(), 1)

	if err != nil {
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
