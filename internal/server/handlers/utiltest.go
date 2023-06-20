package handlers

import (
	"github.com/gin-gonic/gin"
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
	var linksStorage storage.Storage

	if s != nil {
		linksStorage = s
	} else {
		linksStorage, _ = storage.NewStorage()
	}

	r := engine.SetupEngine(linksStorage)
	return SetupRouter(r, linksStorage)
}
