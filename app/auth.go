package app

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/0xdod/trove"
	"github.com/gorilla/sessions"
)

type key string

var store = sessions.NewCookieStore([]byte("my-really-secret-key"))

func createSession(user *trove.User) http.Handler {
	sessionName := "auth"
	userKey := "user.id"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, sessionName)
		session.Values[userKey] = user.ID
		if err := session.Save(r, w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

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
		createSession(user).ServeHTTP(w, r)

		if err != nil {
			s.serverErrorResponse(w, err)
			return
		}

		err = s.writeJSON(w, http.StatusOK, RM{"success", "user token created", M{"token": token}})

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

func (s *Server) loginRequired(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := UserFromContext(r.Context())
		if user.IsAnonymous() {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func (s *Server) SessionAuth(h http.Handler) http.Handler {
	sessionName := "auth"
	userKey := "user.id"

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user *trove.User
		var err error
		session, _ := store.Get(r, sessionName)
		id, ok := session.Values[userKey]

		if !ok {
			user = &trove.User{}
		} else {
			idInt := id.(int)
			user, err = s.UserService.FindUserByID(r.Context(), idInt)

			if err != nil {
				user = &trove.User{}
			}
		}

		ctx := context.WithValue(r.Context(), key("user"), user)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}
