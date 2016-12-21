package v1

import (
	"fmt"
	"strconv"

	"gopkg.in/gin-gonic/gin.v1"
	"tuohai/im_api/models"
	"tuohai/im_api/render"
	"tuohai/internal/auth"
	"tuohai/internal/console"
)

func GetMsgCollect(ctx *gin.Context) {
	main_user := ctx.MustGet("user").(*auth.MainUser)
	var (
		res   []models.MsgCollect
		total int
		err   error
	)
	limit, _ := strconv.Atoi(ctx.Request.FormValue("limit"))
	offset, _ := strconv.Atoi(ctx.Request.FormValue("offset"))
	pageindex, _ := strconv.Atoi(ctx.Request.FormValue("pageindex"))
	pagesize, _ := strconv.Atoi(ctx.Request.FormValue("pagesize"))
	if limit == 0 && offset == 0 {
		res, total, err = models.CollectsByPaging(main_user.Uid, pageindex, pagesize)
	}
	if pageindex == 0 && pagesize == 0 {
		res, total, err = models.CollectsByOffset(main_user.Uid, limit, offset)
	}

	fmt.Println(total)
	if err != nil {
		render.RenderJSON(ctx, struct{}{})
	} else {
		render.RenderJSON(ctx, gin.H{
			"total": total,
			"list":  res,
		})
	}
}

func AddMsgCollect(ctx *gin.Context) {
	cid := ctx.PostForm("cid")
	mid, _ := strconv.Atoi(ctx.PostForm("mid"))
	ctype := ctx.PostForm("type")
	main_user := ctx.MustGet("user").(*auth.MainUser)
	fmt.Println(cid, mid, ctype)

	err := models.AddMsgCollect(main_user.Uid, cid, ctype, uint64(mid))
	if err != nil {
		console.StdLog.Error(err)
		render.RenderJSON(ctx, false)
	} else {
		render.RenderJSON(ctx, true)
	}
	return
}

func DelMsgCollect(ctx *gin.Context) {
	cid := ctx.Request.FormValue("cid")
	mid, _ := strconv.Atoi(ctx.Request.FormValue("mid"))
	// ctype := ctx.PostForm("type")
	// main_user := ctx.MustGet("user").(*auth.MainUser)
	fmt.Println(ctx.Request.Form)
	if cid == "" {
		render.RenderJSON(ctx, struct{}{}, 1, "cid 不能为空")
		return
	}
	if mid == 0 {
		render.RenderJSON(ctx, struct{}{}, 1, "mid 不能为空")
		return
	}
	err := models.DelMsgCollect(cid, mid)
	if err != nil {
		console.StdLog.Error(err)
		render.RenderJSON(ctx, false)
	} else {
		render.RenderJSON(ctx, true)
	}
	return
}
