package v1

import (
	// "fmt"
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
	"tuohai/internal/console"
	"tuohai/models"
)

func Profile() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		renderJSON(ctx, 0, nil)
	}
}

func Groups() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.MustGet("token").(string)
		mems_groups, err := models.AssociationGroups(token)
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, []int{}, 0, "未找到数据")
			return
		}

		var groups []models.TblGroup
		for _, mems := range mems_groups {

			group, err := models.GetTblGroupById(mems.GroupId)
			if err != nil {
				console.StdLog.Error(err)
			}
			groups = append(groups, *group)
		}

		if len(groups) == 0 {
			renderJSON(ctx, []int{})
			return
		}

		renderJSON(ctx, groups)
	}
}

func Group() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		gid := ctx.Param("gid")
		token := ctx.MustGet("token").(string)
		group, err := models.GetTblGroupById(gid)
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, struct{}{}, 0, "未找到数据")
			return
		}

		user, err := models.GetTblUserById(token)
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, struct{}{}, 0, "未找到数据")
			return
		}

		ig, err := models.IsGroupMember(group.Id, user.Id)
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, struct{}{}, 0, "未找到数据")
			return
		}

		if ig {
			renderJSON(ctx, group, 0)
		} else {
			renderJSON(ctx, struct{}{}, 1, "当前用户不属于这个群")
		}
	}
}

func UpdateGroup() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func Sessions() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessions, err := models.GetTblSessionById(ctx.MustGet("token").(string))
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, []int{}, 0, "未找到数据")
			return
		}
		renderJSON(ctx, sessions)
	}
}

func Messages() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sid := ctx.Param("sid")
		msg, err := models.GetTblMsgById(sid)
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, []int{}, 0, "未找到数据")
			return
		}

		renderJSON(ctx, msg)
	}
}

func RemoveSession() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func MessageRead() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func Unreads() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func CreateGroup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var group models.TblGroup
		if err := ctx.Bind(&group); err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, []int{}, 0, err)
			return
		}

		g, err := models.CreateGroup(&group)
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, []int{}, 0, "未找到数据")
			return
		}

		renderJSON(ctx, g)
	}
}

func UserInfo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid := ctx.Param("uid")
		user, err := models.GetTblUserById(uid)
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, []int{}, 0, "未找到数据")
			return
		}

		renderJSON(ctx, user)
	}
}

func renderJSON(ctx *gin.Context, json interface{}, err_status ...interface{}) {
	switch len(err_status) {
	case 0:
		ctx.JSON(http.StatusOK, gin.H{"err_code": 0, "data": json})
		break
	case 1:
		ctx.JSON(http.StatusOK, gin.H{"err_code": err_status[0], "data": json})
		break
	case 2:
		ctx.JSON(http.StatusOK, gin.H{"err_code": err_status[0], "err_msg": err_status[1], "data": json})
		break
	}
}

func renderSysError(ctx *gin.Context, err error) {
	if err != nil {
		console.StdLog.Error(err)
		renderJSON(ctx, struct{}{}, 1, "远程服务器错误!")
	}
}
