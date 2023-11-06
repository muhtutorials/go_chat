package main

import (
	"github.com/bmizerany/pat"
	"go_chat/internal/handlers"
	"net/http"
)

func routes() http.Handler {
	mux := pat.New()

	mux.Get("/", http.HandlerFunc(handlers.Home))
	mux.Get("/ws", http.HandlerFunc(handlers.WSEndpoint))

	fileServer := http.FileServer(http.Dir("../../static"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
