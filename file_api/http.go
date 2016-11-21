package file_api

import (
	"gopkg.in/gin-gonic/gin.v1"
	"tuohai/file_api/api/v1"
	"tuohai/internal/console"
)

func newHTTPServer() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(console.Logger())
	version1 := router.Group("v1")
	{
		version1.POST("/upload", v1.Upload())
	}
	return router
}
