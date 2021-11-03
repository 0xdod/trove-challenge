package app

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/0xdod/trove"
)

type key string

func (s *Server) createAuthToken() http.Handler {
	type Request struct {
		Email    string `json:"email,omitempty" validate:"required,email"`
		Password string `json:"password,omitempty" validate:"required"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := Request{}

		if err := s.readJSON(r.Body, &request); err != nil {
			s.badRequestResponse(w, err)
			return
		}

		if err := validate.Struct(request); err != nil {
			s.badRequestResponse(w, err)
			return
		}

		user, err := s.UserService.FindUserByEmail(r.Context(), request.Email)

		if err != nil {
			s.serverErrorResponse(w, err)
			return
		}

		if !user.VerifyPassword(request.Password) {
			s.writeJSON(w, http.StatusUnauthorized, RM{"error", "invalid user credentials", nil})
			return
		}

		token, err := trove.NewAuthToken(user.ID, 24*time.Hour)

		if err != nil {
			s.serverErrorResponse(w, err)
			return
		}

		err = s.AuthService.CreateToken(r.Context(), token)

		if err != nil {
			s.serverErrorResponse(w, err)
			return
		}

		err = s.writeJSON(w, http.StatusOK, RM{"success", "user tokens created", M{"token": token}})

		if err != nil {
			log.Printf("json error: %v", err)
		}
	})
}

func (s *Server) MustAuth(h http.HandlerFunc) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), " ")

		if len(authHeader) < 2 {
			s.writeJSON(w, http.StatusUnauthorized, RM{"error", "authorization failed, token not provided", nil})
			return
		}

		user, err := s.UserService.FindUserByToken(r.Context(), authHeader[1])

		if err != nil {
			s.writeJSON(w, http.StatusUnauthorized, RM{"error", "authorization failed", nil})
			return
		}

		ctx := context.WithValue(r.Context(), key("user"), user)
		r = r.WithContext(ctx)
		h(w, r)
	})
}
