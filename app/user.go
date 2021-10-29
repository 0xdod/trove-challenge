package app

import (
	"log"
	"net/http"

	"github.com/0xdod/trove"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type userSignupRequest struct {
	FirstName string `json:"first_name" validate:"required,min=2"`
	LastName  string `json:"last_name" validate:"required,min=2"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6,max=60"`
}

func (usr *userSignupRequest) ValidatedUser() (*trove.User, error) {
	if err := validate.Struct(usr); err != nil {
		return nil, err
	}

	user := trove.User{
		FirstName: usr.FirstName,
		LastName:  usr.LastName,
		Email:     usr.Email,
	}

	return &user, nil
}

func (s *Server) RegisterUser(w http.ResponseWriter, r *http.Request) {

	userDetails := &userSignupRequest{}

	if err := s.readJSON(r.Body, &userDetails); err != nil {
		// return error message
		s.writeJSON(w, http.StatusUnprocessableEntity, RM{"error", "request error: " + err.Error(), nil})
		return
	}

	// do some request validation
	user, err := userDetails.ValidatedUser()

	if err != nil {
		// return error message
		s.writeJSON(w, http.StatusBadRequest, RM{"error", "validation error: " + err.Error(), nil})
		return
	}

	if err = user.SetPassword(userDetails.Password); err != nil {
		// error
		s.writeJSON(w, http.StatusInternalServerError, RM{"error", "internal error: " + err.Error(), nil})
		return
	}

	// do business logic thingy, in this case signup
	if err = s.UserService.Create(r.Context(), user); err != nil {
		// do some error thing.
		s.writeJSON(w, http.StatusInternalServerError, RM{"error", "internal error: " + err.Error(), nil})
		return
	}

	// return response
	if err := s.writeJSON(w, http.StatusCreated, RM{"success", "created user account", user}); err != nil {
		log.Println(err)
	}
}

func (s *Server) handleSignup() http.Handler {
	type request struct{}

	type response struct{}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}

func (s *Server) coolHandler() http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}
