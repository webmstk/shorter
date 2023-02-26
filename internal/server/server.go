package server

import (
	"log"

	"github.com/webmstk/shorter/internal/config"
	"github.com/webmstk/shorter/internal/server/engine"
	"github.com/webmstk/shorter/internal/server/handlers"
	"github.com/webmstk/shorter/internal/storage"
)

func Run() {
	linksStorage := storage.NewStorage()
	r := engine.SetupEngine()
	handlers.SetupRouter(r, linksStorage)
	log.Println("Starting web-server at", config.ServerBaseURL)
	err := r.Run(config.Config.ServerAddress)
	log.Fatal(err)
}
