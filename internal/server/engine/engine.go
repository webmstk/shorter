package engine

import (
	"github.com/gin-gonic/gin"
	"github.com/webmstk/shorter/internal/server/middlewares"
)

func SetupEngine() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(middlewares.Compressor())
	return r
}
