package v1

import (
	"gopkg.in/gin-gonic/gin.v1"
	"tuohai/im_api/models"
	"tuohai/internal/auth"
	"tuohai/internal/console"
	"tuohai/internal/convert"
)

func ApplyFriends() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		main_user := ctx.MustGet("user").(*auth.MainUser)
		apply, err := models.FriendApplys(main_user.Uid)
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, struct{}{}, 1, "远程服务器错误")
			return
		}

		renderJSON(ctx, apply)
	}
}

func ConfirmApplyFriend() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// main_user := ctx.MustGet("user").(*auth.MainUser)
		id := ctx.PostForm("id")
		status := ctx.PostForm("status")
		if status == "" || id == "" {
			renderJSON(ctx, struct{}{}, 1, "id 或 status不能为空")
			return
		}
		status_int := convert.StrTo(status).MustInt()
		err := models.SaveFriendApply(&models.FriendApply{Id: id, Status: models.ApplyType(status_int)})
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, struct{}{}, 1, "远程服务器错误")
			return
		}

		renderJSON(ctx, true)
	}
}
