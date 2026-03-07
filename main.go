package main

import (
	"log"
	"net/http"

	"github.com/johnfarrell/stylesheets/handlers"
)

func main() {
	mux := handlers.NewMux()
	addr := ":8080"
	log.Printf("Starting server on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
