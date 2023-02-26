package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/webmstk/shorter/internal/config"
	"github.com/webmstk/shorter/internal/server/engine"
	"github.com/webmstk/shorter/internal/storage"
)

func generateShortLink(s string) string {
	s, err := storage.GenerateShortLink(s)
	if err != nil {
		panic("Failed to generate short link")
	}
	return s
}

func setupServer(s storage.Storage) *gin.Engine {
	r := engine.SetupEngine()
	var linksStorage storage.Storage

	if s != nil {
		linksStorage = s
	} else {
		linksStorage = storage.NewStorage()
	}
	return SetupRouter(r, linksStorage)
}

func setupTestConfig(config *config.AppConfig) {
	config.ServerAddress = "localhost:8080"
	config.BaseURL = "http://localhost:8080"
	config.FileStoragePath = ""
}
