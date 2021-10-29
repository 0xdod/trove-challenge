package main

import (
	"log"

	"github.com/0xdod/trove/server"
)

func main() {
	s := server.New()
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
