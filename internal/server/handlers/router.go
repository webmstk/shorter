package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/webmstk/shorter/internal/storage"
)

func SetupRouter(r *gin.Engine, linksStorage storage.Storage) *gin.Engine {
	r.POST("/", HandlerShorten(linksStorage))
	r.GET("/:shortURL", HandlerExpand(linksStorage))

	apiGroup := r.Group("/api")
	{
		apiGroup.POST("/shorten", HandlerAPIShorten(linksStorage))
	}

	return r
}
