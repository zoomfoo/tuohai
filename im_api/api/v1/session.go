package v1

import (
	"gopkg.in/gin-gonic/gin.v1"

	"tuohai/im_api/models"
	"tuohai/internal/auth"
)

// 屏蔽或解除屏蔽临时会话
func ShieldSession(flag bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("user").(*auth.MainUser)
		cid := ctx.Query("cid")
		if len(cid) == 0 || len(cid) > 64 {
			renderJSON(ctx, struct{}{}, 1, "参数非法")
			return
		}
		r, err := models.Friend(user.Uid, cid)
		if err != nil || r.Rtype != 1 {
			renderJSON(ctx, struct{}{}, 1, "临时关系不存在")
			return
		}
		if flag {
			// 屏蔽
			err = models.ShieldProcess(cid, user.Uid, 1)
		} else {
			// 解除屏蔽
			err = models.ShieldProcess(cid, user.Uid, 0)
		}
		if err != nil {
			renderJSON(ctx, struct{}{}, 1, "处理有误")
			return
		}
		renderJSON(ctx, struct{}{})
	}
}
