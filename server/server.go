package server

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	server *http.Server
	router *mux.Router
}

func New() *Server {
	router := mux.NewRouter().StrictSlash(true)
	server := &http.Server{
		Addr:         ":8000",
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	return &Server{server, router}
}

func (s *Server) Run() error {
	return s.server.ListenAndServe()
}
