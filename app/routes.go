package app

func (s *Server) routes() {
	s.router.Handle("/", s.handleIndex())
	s.router.Handle("/login", s.handleLogin())
	s.router.Handle("/signup", s.handleSignup())

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
