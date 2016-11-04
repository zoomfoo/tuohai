package open_api

import (
	"gopkg.in/gin-gonic/gin.v1"
	"tuohai/internal/console"
	"tuohai/open_api/api/v1"
)

func newHTTPServer() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
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

		app := version1.Group("apps")
		{
			app.GET("/", v1.Apps())
		}

		//从第三方接到的webhook
		version1.POST("/hook/:bot_id", v1.PushHook())
		//接受消息推送
		version1.POST("/push_msg/:bot_id", v1.PushMsg())
	}

	return router
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
