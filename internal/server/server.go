package server

import (
	"github.com/webmstk/shorter/internal/config"
	"github.com/webmstk/shorter/internal/server/handlers"
	"github.com/webmstk/shorter/internal/storage"
	"log"
	"net/http"
	"time"
)

func Run() {
	mux := http.NewServeMux()
	s := &http.Server{
		Addr:         config.ServerHost + ":" + config.ServerPort,
		Handler:      mux,
		IdleTimeout:  10 * time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	}

	storage := storage.NewStorage()
	mux.Handle("/", http.HandlerFunc(handlers.HandlerLinks(storage)))

	log.Println("Starting web-server at", config.ServerBaseURL)
	err := s.ListenAndServe()
	log.Fatal(err)
}
