package app

import "net/http"

func (s *Server) dummyHandlerFunc(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) coolHandler() http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}
