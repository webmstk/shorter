package engine

import (
	"github.com/gin-gonic/gin"
	"github.com/webmstk/shorter/internal/server/middlewares"
	"github.com/webmstk/shorter/internal/storage"
)

func SetupEngine(storage storage.Storage) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(middlewares.Compressor())
	r.Use(middlewares.UserCookie(storage))
	return r
}
