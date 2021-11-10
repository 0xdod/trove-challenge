package main

import (
	"fmt"
	"log"
	"os"

	"github.com/0xdod/trove/app"
	"github.com/0xdod/trove/postgres"
)

type config struct {
	port  string
	dbURL string
}

func main() {
	cfg := config{
		port:  getPort(),
		dbURL: getDbURL(),
	}

	db, err := postgres.Open(cfg.dbURL)

	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	log.Println("connected to database successfully")

	s := app.NewServer(db)
	if err := s.Run(cfg.port); err != nil {
		log.Fatal(err)
	}
}

func getPort() string {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	return port
}

func getDbURL() string {
	dbURL, exists := os.LookupEnv("POSTGRESQL_URL")

	if !exists {
		host := os.Getenv("DB_HOST")
		user := os.Getenv("DB_USER")
		pass := os.Getenv("DB_PASSWORD")
		name := os.Getenv("DB_NAME")
		port := os.Getenv("DB_PORT")
		dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, pass, host, port, name)
	}

	return dbURL
}
