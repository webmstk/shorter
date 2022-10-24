package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/webmstk/shorter/internal/storage"
)

func SetupRouter(linksStorage storage.Storage) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.POST("/", HandlerShorten(linksStorage))
	r.GET("/:shortURL", HandlerExpand(linksStorage))
	apiGroup := r.Group("/api")
	{
		apiGroup.POST("/shorten", HandlerAPIShorten(linksStorage))
	}
	return r
}
