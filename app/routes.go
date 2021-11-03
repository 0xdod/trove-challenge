package app

func (s *Server) routes() {
	s.router.HandleFunc("/users", s.registerUser).Methods("POST")
	s.router.Handle("/users/{id}", s.MustAuth(s.updateUser)).Methods("PATCH", "PUT")
	s.router.Handle("/users/{id}/portfolio", s.MustAuth(s.getPortfolio)).Methods("GET")
	s.router.Handle("/users/{id}/portfolio/value", s.MustAuth(s.getPortfolioValue)).Methods("GET")
	s.router.Handle("/loan", s.MustAuth(s.processLoan())).Methods("POST")
	s.router.Handle("/auth/token", s.createAuthToken()).Methods("POST")
}
