package server

import (
	"github.com/gin-gonic/gin"
	"github.com/webmstk/shorter/internal/config"
	"github.com/webmstk/shorter/internal/server/handlers"
	"github.com/webmstk/shorter/internal/storage"
	"log"
)

func Run() {
	gin.SetMode(gin.ReleaseMode)
	linksStorage := storage.NewStorage()
	r := handlers.SetupRouter(linksStorage)
	log.Println("Starting web-server at", config.ServerBaseURL)
	err := r.Run(config.ServerHost + ":" + config.ServerPort)
	log.Fatal(err)
}
