package render

import (
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
	"tuohai/internal/console"
)

func RenderJSON(ctx *gin.Context, json interface{}, err_status ...interface{}) {
	switch len(err_status) {
	case 0:
		ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": json})
		break
	case 1:
		ctx.JSON(http.StatusOK, gin.H{"code": err_status[0], "data": json})
		break
	case 2:
		ctx.JSON(http.StatusOK, gin.H{"code": err_status[0], "msg": err_status[1], "data": json})
		break
	}
}

func RenderSysError(ctx *gin.Context, err error) {
	if err != nil {
		console.StdLog.Error(err)
		RenderJSON(ctx, struct{}{}, 1, "远程服务器错误!")
	}
}
