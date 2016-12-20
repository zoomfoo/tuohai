package v1

import (
	"gopkg.in/gin-gonic/gin.v1"
)

func AddMsgCollect(ctx *gin.Context) {
	// 	cid        房间id(chennel id)
	// mid        消息id
	// type      只有两种选择 1 个人消息  2 群主消息
	ctx.PostForm("cid")
	ctx.PostForm("mid")
	ctx.PostForm("type")
}
