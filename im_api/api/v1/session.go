package v1

import (
	"gopkg.in/gin-gonic/gin.v1"

	_ "tuohai/im_api/models"
	_ "tuohai/internal/auth"
)

// 屏蔽或解除屏蔽临时会话
func ShieldSession(flag bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		renderJSON(ctx, struct{}{})
	}
}
