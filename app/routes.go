package app

func (s *Server) routes() {
	apiRouter := s.router.PathPrefix("/api/v1").Subrouter()
	apiRouter.HandleFunc("/users", s.registerUser).Methods("POST")
	apiRouter.Handle("/users/{id}", s.MustAuth(s.updateUser)).Methods("PATCH", "PUT")
	apiRouter.Handle("/users/{id}/portfolio", s.MustAuth(s.getPortfolio)).Methods("GET")
	apiRouter.Handle("/users/{id}/portfolio/value", s.MustAuth(s.getPortfolioValue)).Methods("GET")
	apiRouter.Handle("/loans", s.MustAuth(s.processLoan())).Methods("POST")
	apiRouter.Handle("/loans", s.MustAuth(s.listLoans)).Methods("GET")
	apiRouter.Handle("/auth/token", s.createAuthToken()).Methods("POST")
}
