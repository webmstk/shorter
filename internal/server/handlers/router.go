package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/webmstk/shorter/internal/storage"
)

func SetupRouter(linksStorage storage.Storage) *gin.Engine {
	r := gin.Default()
	r.POST("/", HandlerShorten(linksStorage))
	r.GET("/:shortURL", HandlerExpand(linksStorage))
	return r
}
