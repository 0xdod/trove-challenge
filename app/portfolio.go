package app

import (
	"log"
	"net/http"
)

func (s *Server) getPortfolio(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())

	if user == nil {
		s.authErrorResponse(w)
		return
	}

	portfolio, err := s.PortfolioService.FindByUser(r.Context(), user.ID)

	if err != nil {
		s.writeJSON(w, http.StatusInternalServerError, RM{"error", "internal error", nil})
		return
	}

	if err := s.writeJSON(w, http.StatusOK, RM{"success", "portfolio positions retrieved", portfolio}); err != nil {
		log.Printf("json error %v", err)
	}
}

func (s *Server) getPortfolioValue(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())

	if user == nil {
		s.authErrorResponse(w)
		return
	}

	value, err := s.PortfolioService.GetPortfolioValue(r.Context(), user.ID)

	if err != nil {
		s.writeJSON(w, http.StatusInternalServerError, RM{"error", "internal error", nil})
		return
	}

	err = s.writeJSON(w, http.StatusOK, RM{"success", "portfolio value retrieved", M{
		"portfolio_value": value,
	}})

	if err != nil {
		log.Printf("json error %v", err)
	}
}
