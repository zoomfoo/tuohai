package open_api

import (
	"gopkg.in/gin-gonic/gin.v1"
	"tuohai/internal/console"
	"tuohai/open_api/api/v1"
)

func newHTTPServer() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(AccessControlAllowOrigin())
	router.Use(console.Logger())

	version1 := router.Group("v1")
	{
		bot := version1.Group("bots", SessionAuth())
		{
			bot.GET("/", v1.BotList())
			bot.POST("/", v1.CreateBot())
			bot.PUT("/:botid", v1.UpdateBot())
			bot.DELETE("/:bot_id", v1.DeleteBot())
		}

		version1.GET("/apps", v1.Apps())

		//从第三方接到的webhook
		version1.POST("/hook/:bot_id", v1.PushHook())
		//接受消息推送
		version1.POST("/push_msg/:bot_id", v1.PushMsg())
	}

	return router
}

func AccessControlAllowOrigin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Add("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,OPTIONS")
		ctx.Writer.Header().Add("Access-Control-Allow-Headers", "Content-Type, Depth, User-Agent, X-File-Size, X-Requested-With, X-Requested-By, If-Modified-Since, X-File-Name, Cache-Control, Token, session_token")
		ctx.Next()
	}
}

func SessionAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// if ctx.Request.Header.Get("token") != "123456789" {
		// ctx.Abort()
		// } else {
		// ctx.Next()
		// }
	}
}
