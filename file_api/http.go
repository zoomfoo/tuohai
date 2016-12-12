package file_api

import (
	"gopkg.in/gin-gonic/gin.v1"
	"tuohai/file_api/api/v1"
	"tuohai/internal/auth"
	"tuohai/internal/console"
)

func newHTTPServer() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(console.Logger())

	version1 := router.Group("v1", auth.LoginAuth(Opts.AuthHost, 0))
	{
		version1.POST("/upload", v1.Upload())
		version1.GET("/files", v1.Files())
	}
	return router
}
