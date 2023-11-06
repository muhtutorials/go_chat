package main

import (
	"fmt"
	"go_chat/internal/handlers"
	"log"
	"net/http"
)

const port = "8080"

func main() {
	mux := routes()

	log.Println("Starting channel listener")
	go handlers.ListenToWSChannel()

	log.Println("Starting web server on port:", port)
	_ = http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
}
