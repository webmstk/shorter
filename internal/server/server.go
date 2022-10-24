package server

import (
	"github.com/webmstk/shorter/internal/server/handlers"
	"log"

	"github.com/webmstk/shorter/internal/config"
	"github.com/webmstk/shorter/internal/storage"
)

func Run() {
	linksStorage := storage.NewStorage()
	r := handlers.SetupRouter(linksStorage)
	log.Println("Starting web-server at", config.ServerBaseURL)
	err := r.Run(config.ServerHost + ":" + config.ServerPort)
	log.Fatal(err)
}
