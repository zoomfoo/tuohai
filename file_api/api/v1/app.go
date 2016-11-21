package v1

import (
	"gopkg.in/gin-gonic/gin.v1"
)

func Upload() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		f, h, err := ctx.Request.FormFile("file")

	}
}
