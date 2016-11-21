package im_api

import (
	"context"
	"net/http"
	"net/http/pprof"
	"time"

	"gopkg.in/gin-gonic/gin.v1"
	"tuohai/im_api/api/v1"
	"tuohai/internal/console"
)

func newHTTPServer() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(console.Logger())
	router.Use(AccessControlAllowOrigin())

	version1 := router.Group("v1", LoginAuth())
	{
		ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
		//列出IM常用的信息
		version1.GET("/im/profile", v1.Profile(ctx))

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
		version1.DELETE("/sessions/:sid", v1.RemoveSession()) //no
		//获取消息历史记录 √
		version1.GET("/sessions/:sid/messages", v1.Messages())
		//消息已读确认 这个read 在restfull中为名词
		version1.PUT("/sessions/:sid/read", v1.MessageRead())

		//获取所有未读消息
		version1.GET("/unreads", v1.Unreads())

		//获取用户信息
		version1.GET("/user/:uid", v1.UserInfo())

		//获取好友列表
		version1.GET("/friends", v1.Friends())
		version1.GET("/friends/:f_uuid", v1.Friend())

		//文件上传
		version1.POST("/files", v1.UploadFile())
		version1.GET("/files", v1.Files())
	}

	//登录
	router.POST("/login", v1.Login())

	Debug(router)
	return router
}

func LoginAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if token := ctx.Query("session_token"); token == "" {
			ctx.Abort()
			ctx.JSON(http.StatusUnauthorized, gin.H{"err_code": 1, "data": "无权限访问"})
		} else {
			ctx.Set("token", token)
			ctx.Next()
		}
		return
	}
}

func AccessControlAllowOrigin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Add("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,OPTIONS")
		ctx.Writer.Header().Add("Access-Control-Allow-Headers", "Content-Type, Depth, User-Agent, X-File-Size, X-Requested-With, X-Requested-By, If-Modified-Since, X-File-Name, Cache-Control, Token, session_token")
		ctx.Next()
	}
}

func Debug(router *gin.Engine) {

	r := router.Group("/im/sys/debug")
	r.GET("/pprof", func(ctx *gin.Context) {
		pprof.Index(ctx.Writer, ctx.Request)
	})
	r.GET("/pprof/cmdline", func(ctx *gin.Context) {
		pprof.Cmdline(ctx.Writer, ctx.Request)
	})
	r.GET("/pprof/symbol", func(ctx *gin.Context) {
		pprof.Symbol(ctx.Writer, ctx.Request)
	})
	r.POST("/pprof/symbol", func(ctx *gin.Context) {
		pprof.Symbol(ctx.Writer, ctx.Request)
	})
	r.GET("/pprof/profile", func(ctx *gin.Context) {
		pprof.Profile(ctx.Writer, ctx.Request)
	})
	r.GET("/pprof/heap", func(ctx *gin.Context) {
		pprof.Handler("heap")
	})

	r.GET("/goroutine", func(ctx *gin.Context) {
		pprof.Handler("goroutine").ServeHTTP(ctx.Writer, ctx.Request)
	})
	r.GET("/block", func(ctx *gin.Context) {
		pprof.Handler("block").ServeHTTP(ctx.Writer, ctx.Request)
	})
	r.GET("/heap", func(ctx *gin.Context) {
		pprof.Handler("heap").ServeHTTP(ctx.Writer, ctx.Request)
	})
	r.GET("/threadcreate", func(ctx *gin.Context) {
		pprof.Handler("threadcreate").ServeHTTP(ctx.Writer, ctx.Request)
	})
}
