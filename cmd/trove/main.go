package main

import (
	"log"

	"github.com/0xdod/trove/app"
)

func main() {
	s := app.NewServer()
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
