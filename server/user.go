package server

import (
	"log"
	"net/http"

	"github.com/0xdod/trove"
)

type userSignupRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (usr *userSignupRequest) User() *trove.User {
	user := trove.User{
		FirstName: usr.FirstName,
		LastName:  usr.LastName,
		Email:     usr.Email,
	}

	return &user
}

func (s *Server) registerUserRoutes() {
	r := s.router
	r.HandleFunc("/users", s.RegisterUser)
}

func (s *Server) RegisterUser(w http.ResponseWriter, r *http.Request) {
	userDetails := &userSignupRequest{}

	if err := s.readJSON(r.Body, &userDetails); err != nil {
		// return error message
	}

	// do some validation
	user := userDetails.User()
	_ = user.SetPassword(userDetails.Password)

	// do business logic thingy, in this case signup
	s.UserService.Create(r.Context(), user)
	// return response
	err := s.writeJSON(w, http.StatusCreated, user)

	if err != nil {
		log.Println(err)
	}
}
