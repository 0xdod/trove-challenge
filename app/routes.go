package app

import (
	"net/http"
)

func (s *Server) routes() {
	htmlRouter := s.router.PathPrefix("/").Subrouter()
	htmlRouter.Use(s.SessionAuth)
	htmlRouter.Use(s.AddDefaultContext)
	{
		htmlRouter.Handle("/", s.loginRequired(s.handleIndex()))
		htmlRouter.Handle("/login", s.handleLogin())
		htmlRouter.Handle("/signup", s.handleSignup())
		htmlRouter.Handle("/profile", s.loginRequired(s.handleProfileUpdate()))
	}

	s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./ui/assets"))))

	apiRouter := s.router.PathPrefix("/api/v1").Subrouter()
	{
		apiRouter.HandleFunc("/users", s.registerUser).Methods("POST")
		apiRouter.Handle("/users/{id}", s.MustAuth(s.updateUser)).Methods("PATCH", "PUT")
		apiRouter.Handle("/portfolio", s.MustAuth(s.getPortfolio)).Methods("GET")
		apiRouter.Handle("/portfolio/value", s.MustAuth(s.getPortfolioValue)).Methods("GET")
		apiRouter.Handle("/loans", s.MustAuth(s.processLoan())).Methods("POST")
		apiRouter.Handle("/loans", s.MustAuth(s.listLoans)).Methods("GET")
		apiRouter.Handle("/auth/token", s.createAuthToken()).Methods("POST")
	}
}
