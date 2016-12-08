package im_api

import (
	"net/http/pprof"

	"gopkg.in/gin-gonic/gin.v1"
	"tuohai/im_api/api/v1"
	"tuohai/internal/auth"
	"tuohai/internal/console"
)

func newHTTPServer() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(console.Logger())
	router.Use(AccessControlAllowOrigin())

	router.OPTIONS("/", func(ctx *gin.Context) {})
	version1 := router.Group("v1", auth.LoginAuth(Opts.AuthHost))
	{
		//列出IM常用的信息
		version1.GET("/profile", v1.Profile())
		version1.PUT("/profile", v1.PutProfile(Opts.AuthHost))

		//群组创建 更新
		//成员管理
		//管理员管理
		groups := version1.Group("groups")
		{
			//获取群组列表 √
			groups.GET("", v1.Groups())
			//获取群组信息 √
			groups.GET("/:gid", v1.Group())
			//创建群组 √
			groups.POST("", v1.CreateGroup())

			//群管理
			//群重命名 √
			groups.PUT("/:gid/name", v1.GroupRename())
			//解散群 √
			groups.DELETE("/:gid/dismiss", v1.DismissGroup())
			//退出群
			groups.DELETE("/:gid/quit", v1.QuitGroupMember())
			//添加群成员
			groups.POST("/:gid/add", v1.AddGroupMember())
			//移除群成员
			groups.DELETE("/:gid/remove", v1.RemoveGroupMember())
		}
		//创建项目群
		version1.POST("/project/groups", v1.CreateProjectGroup())
		version1.GET("/teams", v1.Teams())
		version1.POST("/feedback", v1.Feedback())

		sessions := version1.Group("sessions")
		{
			//获取session列表 √
			sessions.GET("", v1.Sessions())
			//删除session no
			sessions.DELETE("/:sid", v1.RemoveSession())
			//获取消息历史记录 √
			//消息已读确认 这个read 在restfull中为名词
			sessions.PUT("/:sid/read", v1.MessageRead())
		}

		messages := version1.Group("messages")
		{
			messages.GET("/:cid", v1.Messages())
		}

		poke := version1.Group("pokes")
		{
			// 戳一下
			poke.POST("", v1.AddChuo())
			// 确认收到
			poke.POST("/:pid/confirm", v1.ConfirmChuo())
			// 获取戳的详情
			poke.GET("/:pid", v1.GetChuoInfo())

			poke.GET("", v1.GetChuoListFrom())
			// 获取我发出
			version1.GET("/poke/send", v1.GetChuoListFrom())
			// 获取我收到的戳
			version1.GET("/poke/recv", v1.GetChuoListRcv())
		}

		friends := version1.Group("friends")
		{
			friends.GET("", v1.Friends(Opts.AuthHost))
			friends.GET("/:fid", v1.Friend())
			friends.POST("", v1.AddFriend())
			friends.DELETE("", v1.DelFriend())
		}

		apply := version1.Group("apply")
		{
			apply.GET("/friends", v1.ApplyFriends())
			apply.PUT("/friends", v1.AgreeApplyFriend())
		}

		//获取所有未读消息
		version1.GET("/unreads", v1.Unreads())
		//获取用户信息
		version1.GET("/user/:uid", v1.UserInfo())
		//获取好友列表

		//获取消息未读详情信息
		version1.GET("/readinfo/:cid/:msgid/:origin", v1.MessageRead())

		version1.GET("/get_toids", v1.GetIds())
	}

	//登录
	router.POST("/login", v1.Login())

	Debug(router)
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
