package app

func (s *Server) routes() {
	s.router.HandleFunc("/users", s.registerUser).Methods("POST")
	s.router.HandleFunc("/users/{id}", s.updateUser).Methods("PATCH", "PUT")
	s.router.HandleFunc("/users/{id}/portfolio", s.getPortfolio).Methods("GET")
	s.router.HandleFunc("/users/{id}/portfolio/value", s.getPortfolioValue).Methods("GET")
	s.router.HandleFunc("/loan", s.processLoan).Methods("POST")
}
