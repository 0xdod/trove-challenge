package app

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/0xdod/trove"
	"github.com/0xdod/trove/postgres"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

type Server struct {
	server *http.Server
	router *mux.Router

	trove.UserService
	trove.PortfolioService
	trove.AuthService
	trove.LoanService
}

type M map[string]interface{}

const (
	host     = "localhost"
	user     = "admin"
	password = "admin"
	dbname   = "trove"
	port     = 5432
	sslmode  = "disable"
	timezone = "Africa/Lagos"
)

func NewServer() *Server {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		host, user, password, dbname, port, sslmode, timezone)

	db, err := postgres.Open(dsn)

	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}

	log.Println("connected to database successfully")

	router := mux.NewRouter().StrictSlash(true)
	n := negroni.New()
	n.Use(negroni.NewLogger())
	n.Use(negroni.NewRecovery())
	n.UseHandler(router)
	server := &http.Server{
		Addr:         ":8000",
		Handler:      n,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	s := &Server{
		server:           server,
		router:           router,
		UserService:      postgres.NewUserService(db),
		PortfolioService: postgres.NewPortfolioService(db),
		AuthService:      postgres.NewAuthService(db),
		LoanService:      postgres.NewLoanService(db),
	}

	s.routes()

	return s
}

func (s *Server) Run() error {
	log.Printf("server started on port %s", s.server.Addr)
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

type genericResponseModel struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type RM = genericResponseModel
