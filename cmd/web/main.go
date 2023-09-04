package main

import (
	"go_ws/internal/handlers"
	"log"
	"net/http"
)

func main() {
	mux := routes()

	log.Println("Starting channel listener")
	go handlers.ListenToWSChannel()

	log.Println("Starting web server on port 8000")
	_ = http.ListenAndServe(":8000", mux)
}
