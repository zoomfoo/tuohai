package file_api

import (
	"fmt"

	"gopkg.in/gin-gonic/gin.v1"
	"tuohai/file_api/api/v1"
	"tuohai/internal/auth"
	"tuohai/internal/console"
)

func newHTTPServer() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(AccessControlAllowOrigin())
	router.Use(console.Logger())

	version1 := router.Group("v1", auth.LoginAuth(Opts.AuthHost, 0))
	{
		version1.POST("/upload", v1.Upload())
		version1.GET("/files", v1.Files())
		version1.POST("/upload/avatar", v1.UploadAvatar())
	}
	return router
}

func AccessControlAllowOrigin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Add("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,OPTIONS")
		ctx.Writer.Header().Add("Access-Control-Allow-Headers", "Content-Type, Depth, User-Agent, X-File-Size, X-Requested-With, X-Requested-By, If-Modified-Since, X-File-Name, Cache-Control, Token, session_token")
		fmt.Println(ctx.Request.Method)
		if ctx.Request.Method == "OPTIONS" {
			ctx.JSON(200, "")
			ctx.Abort()
		}
		ctx.Next()
	}
}
