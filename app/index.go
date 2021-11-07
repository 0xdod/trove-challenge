package app

import "net/http"

func (s *Server) handleIndex() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Render(w, "index.html", nil)
	})
}

func (s *Server) handleLogin() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Render(w, "login.html", nil)
	})
}

func (s *Server) handleSignup() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Render(w, "signup.html", nil)
	})
}
