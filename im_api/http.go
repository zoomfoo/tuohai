package im_api

import (
	// "net/http/pprof"
	// "time"
	"io/ioutil"
	"log"
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
	"tuohai/im_api/api/v1"
	"tuohai/internal/console"
)

func newHTTPServer() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(console.Logger())
	version1 := router.Group("v1", LoginAuth())
	{
		//列出IM常用的信息
		version1.GET("/im/profile", v1.Profile())

		//获取群组列表 √
		version1.GET("/groups", v1.Groups())
		//获取群组信息
		version1.GET("/groups/:gid", v1.Group())
		//更新群组
		version1.PUT("/groups/:gid", v1.UpdateGroup())
		//创建群组
		version1.POST("/groups", v1.CreateGroup())

		//获取session列表 √
		version1.GET("/sessions", v1.Sessions())
		//删除session
		version1.DELETE("/sessions/:sid", v1.RemoveSession())
		//获取消息历史记录 √
		version1.GET("/sessions/:sid/messages", v1.Messages())
		//消息已读确认 这个read 在restfull中为名词
		version1.PUT("/sessions/:sid/read", v1.MessageRead())

		//获取所有未读消息
		version1.GET("/unreads", v1.Unreads())
	}

	router.POST("/wangyang1", func(ctx *gin.Context) {
		data, err := ioutil.ReadAll(ctx.Request.Body)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(ctx.Request.Header.Get("bbb"))
		defer ctx.Request.Body.Close()
		ctx.String(200, "%s", data)
	})
	Debug(router)
	return router
}

func LoginAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if token := ctx.Request.Header.Get("session_token"); token == "" {
			ctx.Abort()
			ctx.JSON(http.StatusUnauthorized, gin.H{"err_code": 1, "data": "无权限访问"})
		} else {
			ctx.Set("token", token)
			ctx.Next()
		}
		return
	}
}

func Debug(router *gin.Engine) {
	// r := router.Group("/im/sys/debug")
	// r.GET("/pprof", pprof.Index)
	// r.GET("/pprof/cmdline", pprof.Cmdline)
	// r.GET("/pprof/symbol", pprof.Symbol)
	// r.POST("/pprof/symbol", pprof.Symbol)
	// r.GET("/pprof/profile", pprof.Profile)
	// r.GET("/pprof/heap", pprof.Handler("heap"))
	// r.GET("/pprof/goroutine", pprof.Handler("goroutine"))
	// r.GET("/pprof/block", pprof.Handler("block"))
	// r.GET("/pprof/threadcreate", pprof.Handler("threadcreate"))
}
