package main

import (
	"log"
	"net/http"
	"os"

	"github.com/johnfarrell/stylesheets/handlers"
)

func main() {
	mux := handlers.NewMux()
	addr := os.Getenv("PORT")
	if addr == "" {
		addr = "8080"
	}
	addr = ":" + addr
	log.Printf("Starting server on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
