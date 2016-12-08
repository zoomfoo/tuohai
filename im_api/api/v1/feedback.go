package v1

import (
	"gopkg.in/gin-gonic/gin.v1"
	"tuohai/im_api/models"
	"tuohai/internal/auth"
	"tuohai/internal/console"
)

func Feedback() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		main_user := ctx.MustGet("user").(*auth.MainUser)
		content := ctx.PostForm("content")
		if content == "" {
			renderJSON(ctx, struct{}{}, 1, "The content cannot be empty")
			return
		}
		if err := models.NewFeedback(main_user.Uid, content); err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, struct{}{}, 1, "远程服务器错误!")
			return
		}
		renderJSON(ctx, true)
		return
	}
}
