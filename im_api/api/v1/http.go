package v1

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"gopkg.in/gin-gonic/gin.v1"
	"tuohai/internal/console"
	"tuohai/internal/util"
	"tuohai/models"
)

func Profile(c context.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// go func() {
		// 	for {
		// 		select {
		// 		case <-c.Done():
		// 			return
		// 		default:
		// 		}
		// 	}
		// }()
		token := ctx.MustGet("token").(string)
		user, err := models.GetTblUserById(token)
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, struct{}{}, 0, "未找到数据")
			return
		}
		renderJSON(ctx, user)
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
		wg := &util.WaitGroupWrapper{}

		var (
			group *models.TblGroup
			user  *models.TblUser
			err   error
		)
		wg.Add(2)
		wg.Wrap(func() {
			group, err = models.GetTblGroupById(gid)
			wg.Done()
		})

		wg.Wrap(func() {
			user, err = models.GetTblUserById(token)
			wg.Done()
		})

		wg.Wait()

		if err != nil {
			console.StdLog.Error(err)
			fmt.Println(err.Error(), err == models.RecordNotFound, models.RecordNotFound.Error())
			if err.Error() == models.RecordNotFound.Error() {
				renderJSON(ctx, struct{}{})
				return
			}
			renderJSON(ctx, struct{}{}, 1, "远程服务器错误!")
			return
		}

		if group == nil {
			renderJSON(ctx, struct{}{})
			return
		}

		ig, err := models.IsGroupMember(group.Gid, user.Uuid)
		if err != nil {
			console.StdLog.Error(err)
			if err == models.RecordNotFound {
				renderJSON(ctx, struct{}{})
			}
			renderJSON(ctx, struct{}{}, 1, "远程服务器错误!")
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
		}

		if sessions == nil {
			renderJSON(ctx, []int{}, 1, "未找到数据")
			return
		}

		var list []gin.H
		for _, session := range sessions {
			history, err := models.GetLastHistory(session.To)
			if err != nil {
				console.StdLog.Error(err)
			}
			list = append(list, gin.H{
				"sid":  session.Sid,
				"fr":   session.From,
				"to":   session.To,
				"mid":  history.MsgId,
				"data": history.MsgData,
				"time": history.CreatedAt,
			})
		}
		renderJSON(ctx, list)
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
		uids := ctx.Param("uid")

		uid := strings.Split(uids, ",")

		user, err := models.GetTblUserByIds(uid)
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, []int{}, 1, "未找到数据")
			return
		}
		renderJSON(ctx, user)
		return
	}
}

func Friends() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.MustGet("token").(string)
		r, err := models.Friends(token)
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, []int{}, 1, "未找到数据")
			return
		}

		var list []interface{}
		for _, rel := range r {
			f_uuid := ""
			switch token {
			case rel.SmallId:
				f_uuid = rel.BigId
			case rel.BigId:
				f_uuid = rel.SmallId
			}

			fuser, err := models.GetTblUserById(f_uuid)
			if err != nil {
				console.StdLog.Error(err)
			}
			list = append(list, gin.H{
				"f_name": fuser.Uname,
				"f_uuid": fuser.Uuid,
			})
		}

		if len(list) == 0 {
			renderJSON(ctx, []int{})
			return
		}
		renderJSON(ctx, list)
	}
}

func Friend() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		f_uuid := ctx.Param("f_uuid")
		token := ctx.MustGet("token").(string)
		rel, err := models.Friend(token, f_uuid)
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, []int{}, 1, "未找到数据")
			return
		}

		switch token {
		case rel.SmallId:
			f_uuid = rel.BigId
		case rel.BigId:
			f_uuid = rel.SmallId
		}

		fuser, err := models.GetTblUserById(f_uuid)
		if err != nil {
			console.StdLog.Error(err)
		}

		renderJSON(ctx, gin.H{
			"f_name": fuser.Uname,
			"f_uuid": fuser.Uuid,
		})
	}
}

func Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uname := ctx.PostForm("user_name")
		pwd := ctx.PostForm("user_pwd")
		user, err := models.Login(uname, pwd)
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, struct{}{}, 1, "用户名或密码错误")
			return
		}
		if user.Id > 0 {
			renderJSON(ctx, user)
			return
		} else {
			renderJSON(ctx, struct{}{}, 0, "用户名或密码错误")
			return
		}
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
