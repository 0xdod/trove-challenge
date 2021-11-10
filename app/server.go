package app

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/0xdod/trove"
	"github.com/0xdod/trove/postgres"
	"github.com/0xdod/trove/ui"
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
	ui.Renderer
}

type M map[string]interface{}

// const (
// 	host     = "localhost"
// 	user     = "admin"
// 	password = "admin"
// 	dbname   = "trove"
// 	port     = 5432
// 	sslmode  = "disable"
// 	timezone = "Africa/Lagos"
// )

func NewServer(db *postgres.DB) *Server {
	// dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
	// 	host, user, password, dbname, port, sslmode, timezone)
	s := &Server{
		server: &http.Server{
			IdleTimeout:  time.Second,
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		},
		router:           mux.NewRouter().StrictSlash(true),
		UserService:      postgres.NewUserService(db),
		PortfolioService: postgres.NewPortfolioService(db),
		AuthService:      postgres.NewAuthService(db),
		LoanService:      postgres.NewLoanService(db),
		Renderer:         ui.NewRenderer(),
	}

	n := negroni.New()
	n.Use(negroni.NewLogger())
	n.Use(negroni.NewRecovery())
	n.UseHandler(s.router)
	s.server.Handler = n
	s.routes()

	return s
}

func (s *Server) Run(port string) error {
	if !strings.HasPrefix(port, ":") {
		port = fmt.Sprintf(":%s", port)
	}
	s.server.Addr = port
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
