package im_api

import (
	"gopkg.in/gin-gonic/gin.v1"
	"tuohai/im_api/api/mainsite"
	"tuohai/im_api/api/v1"
	"tuohai/internal/auth"
	"tuohai/internal/console"
)

func newHTTPServer() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(console.Logger())
	router.Use(AccessControlAllowOrigin())

	version1 := router.Group("v1", auth.LoginAuth())
	{
		//列出IM常用的信息
		version1.GET("/profile", v1.Profile())
		//更新个人信息
		version1.PUT("/profile", v1.PutProfile())
		//批量获取用户信息
		version1.GET("/users", v1.Users())

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
			//groups.DELETE("/:gid/dismiss", v1.DismissGroup())
			//退出群
			groups.DELETE("/:gid/quit", v1.QuitGroupMember())
			//添加群成员
			groups.POST("/:gid/add", v1.AddGroupMember())
			//移除群成员
			groups.POST("/:gid/remove", v1.RemoveGroupMember())
		}

		//获取团队群
		version1.GET("/teams", v1.Teams())
		//反馈
		version1.POST("/feedback", v1.Feedback())

		//session
		sessions := version1.Group("sessions")
		{
			//获取session列表 √
			sessions.GET("", v1.Sessions())
			//删除session no
			sessions.DELETE("/:sid", v1.RemoveSession())
			sessions.DELETE("/:sid/unread", v1.CleanSessionUnread())
		}

		session := version1.Group("session")
		{
			session.GET("/tmp/shield", v1.ShieldSession(true))
			session.GET("/tmp/unshield", v1.ShieldSession(false))
		}

		//消息
		messages := version1.Group("messages")
		{
			//获取消息历史记录 √
			messages.GET("/:cid", v1.Messages())
			//获取消息未读详情信息
			messages.GET("/:cid/readinfo/:msgid/:origin", v1.MessageRead())
		}

		message := version1.Group("message")
		{
			// 获取历史消息(单独窗口中显示时)
			message.GET("/history", v1.MsgHistory())
			message.POST("/forward", v1.ForwardMsg())
			message.POST("/collects", v1.AddMsgCollect())
			message.DELETE("/collects", v1.DelMsgCollect())
			message.GET("/collects", v1.GetMsgCollect())
			message.GET("/info", v1.MessageRead())
		}

		//戳一下
		poke := version1.Group("pokes")
		{
			// 戳一下
			poke.POST("", v1.AddChuo())
			// 确认收到
			poke.POST("/:pid/confirm", v1.ConfirmChuo())
			// 获取戳的详情
			poke.GET("/:pid", v1.GetChuoInfo())

			//获取戳列表
			poke.GET("", v1.GetChuoListFrom())
			// 获取我发出
			version1.GET("/poke/send", v1.GetChuoListFrom())
			// 获取我收到的戳
			version1.GET("/poke/recv", v1.GetChuoListRcv())
		}

		//好友
		friends := version1.Group("friends")
		{
			friends.GET("", v1.Friends())
			friends.GET("/:fid", v1.Friend())
			//添加好友
			friends.POST("", v1.AddFriend())
			friends.DELETE("/:cid", v1.DelFriend())
		}
		// 邀请好友
		version1.GET("/friend/invite", v1.InviteFriend())
		version1.GET("/friend/tmp", v1.CreateTmpFriend())

		//好友申请
		apply := version1.Group("apply")
		{
			//获得好友申请列表
			apply.GET("", v1.ApplyFriends())
			apply.GET("/friends/is/:pageindex/:pagesize", v1.ApplyFriends())
			apply.GET("/friends/not/:pageindex/:pagesize", v1.UnApplyFriends())
			apply.POST("/friends", v1.ProcessApplyFriend())
			apply.PUT("/friends", v1.ProcessApplyFriend())
		}

		//chennel未读确认
		version1.GET("/unreads/:cid", v1.Unreads())

		//获取用户信息
		version1.GET("/user/:uid", v1.UserInfo())
		//获取好友列表

		version1.GET("/get_toids", v1.GetIds())
	}

	//登录
	// router.POST("/login", v1.Login())

	//为主站提供接口
	mainsite.NewMainSiteRouter(router.Group("/m"))

	return router
}

func AccessControlAllowOrigin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Add("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,OPTIONS")
		ctx.Writer.Header().Add("Access-Control-Allow-Headers", "Content-Type, Depth, User-Agent, X-File-Size, X-Requested-With, X-Requested-By, If-Modified-Since, X-File-Name, Cache-Control, Token, session_token")

		if ctx.Request.Method == "OPTIONS" {
			ctx.JSON(200, "")
			ctx.Abort()
		}
		ctx.Next()
	}
}
