package server

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/0xdod/trove"
	"github.com/gorilla/mux"
)

type Server struct {
	server *http.Server
	router *mux.Router

	trove.UserService
}

func New() *Server {
	router := mux.NewRouter().StrictSlash(true)
	server := &http.Server{
		Addr:         ":8000",
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	s := &Server{
		server: server,
		router: router,
	}

	s.registerUserRoutes()

	return s
}

func (s *Server) Run() error {
	return s.server.ListenAndServe()
}

func (*Server) readJSON(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

func (*Server) writeJSON(w io.Writer, code int, v interface{}) error {
	if wr, ok := w.(http.ResponseWriter); ok {
		wr.Header().Set("content-type", "application/json")
		wr.WriteHeader(code)
	}

	return json.NewEncoder(w).Encode(v)
}
