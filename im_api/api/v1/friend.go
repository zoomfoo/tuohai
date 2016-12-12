package v1

import (
	"fmt"
	// "strings"
	"time"

	"gopkg.in/gin-gonic/gin.v1"
	"tuohai/im_api/models"
	"tuohai/internal/auth"
	"tuohai/internal/console"
	"tuohai/internal/convert"
	"tuohai/internal/uuid"
)

func Friends(url string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("user").(*auth.MainUser)
		token := ctx.MustGet("token").(string)
		r, err := models.Friends(user.Uid)
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, []int{}, 1, "未找到数据")
			return
		}
		fmt.Println("friend: ", r, "uid: ", user.Uid)
		var list []gin.H
		for _, rel := range r {
			f_uuid := ""
			switch user.Uid {
			case rel.SmallId:
				f_uuid = rel.BigId
			case rel.BigId:
				f_uuid = rel.SmallId
			}

			u, err := auth.GetBatchUsers(token, url, []string{fmt.Sprintf("user_ids=%s", f_uuid)})
			if err != nil {
				fmt.Println(err)
				continue
			}
			if len(u) == 0 {
				continue
			}

			list = append(list, gin.H{
				"uuid":   u[0].Uuid,
				"name":   u[0].Uname,
				"cid":    rel.Rid,
				"avatar": u[0].Avatar,
				"phone":  u[0].Phone,
				"desc":   u[0].Desc,
			})
		}

		renderJSON(ctx, list)
	}
}

func Friend(url string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		f_uuid := ctx.Param("fid")
		user := ctx.MustGet("user").(*auth.MainUser)
		fmt.Println("接受 fuuid", f_uuid, " uid:", user.Uid)
		rel, err := models.Friend(user.Uid, f_uuid)
		token := ctx.MustGet("token").(string)
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, []int{}, 1, "未找到数据")
			return
		}

		id := ""
		switch user.Uid {
		case rel.SmallId:
			id = rel.BigId
		case rel.BigId:
			id = rel.SmallId
		}
		u, err := auth.GetBatchUsers(token, url, []string{fmt.Sprintf("user_ids=%s", id)})
		if err != nil {
			fmt.Println(err)
		}
		if len(u) == 0 {
			renderJSON(ctx, struct{}{}, 0)
			return
		}

		//获取与我相关群组
		gm, err := models.GroupGhosting(user.Uid, id)
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, struct{}{}, 1, "远程服务器错误")
			return
		}

		var list []gin.H
		for i, _ := range gm {
			g, err := models.GetGroupById(gm[i].GroupId)
			if err != nil {
				continue
			}
			list = append(list, gin.H{
				"gid":  gm[i].GroupId,
				"name": g.Gname,
			})
		}

		renderJSON(ctx, gin.H{
			"uid":      u[0].Uuid,
			"name":     u[0].Uname,
			"cid":      rel.Rid,
			"avatar":   u[0].Avatar,
			"phone":    u[0].Phone,
			"desc":     u[0].Desc,
			"relation": list,
		})

	}
}

func AddFriend() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid := ctx.PostForm("uuid")
		attach := ctx.PostForm("attach")
		way := ctx.PostForm("way")
		phone := ctx.PostForm("phone")
		email := ctx.PostForm("email")
		user := ctx.MustGet("user").(*auth.MainUser)
		if user.Uid == uid {
			renderJSON(ctx, struct{}{}, 1, "不允许添加自己为好友")
			return
		}

		//通过uid添加好友
		if uid != "" {
			res := addFriend(user, uid, way, attach)
			if res == "" {
				renderJSON(ctx, true)
				return
			}
			renderJSON(ctx, struct{}{}, 1, res)
			return
		}

		//通过手机号添加好友
		if phone != "" {
			fmt.Println("添加好友 手机号: ", phone)
			users, err := models.SelectUsers(&models.User{Phone: phone})
			if err != nil {
				console.StdLog.Error(err)
				renderJSON(ctx, []int{}, 1, "远程服务器错误")
				return
			}

			if len(users) == 0 {
				renderJSON(ctx, struct{}{}, 1, "未找到好友")
				return
			}

			res := addFriend(user, users[0].Uuid, way, attach)
			if res == "" {
				renderJSON(ctx, true)
				return
			}
			renderJSON(ctx, struct{}{}, 1, res)
			return
		}

		//通过邮箱添加好友
		if email != "" {
			users, err := models.SelectUsers(&models.User{Email: email})
			if err != nil {
				console.StdLog.Error(err)
				renderJSON(ctx, []int{}, 1, "远程服务器错误")
				return
			}

			if len(users) == 0 {
				renderJSON(ctx, struct{}{}, 1, "未找到好友")
				return
			}

			res := addFriend(user, users[0].Uuid, way, attach)
			if res == "" {
				renderJSON(ctx, true)
				return
			}
			renderJSON(ctx, struct{}{}, 1, res)
			return
		}
	}
}

func addFriend(user *auth.MainUser, uid, way, attach string) string {
	fa := &models.FriendApply{
		Fid:        uuid.NewV4().StringMd5(),
		ApplyUid:   user.Uid,
		TargetUid:  uid,
		Way:        models.ApplyWay(convert.StrTo(way).MustInt()),
		Attach:     attach,
		Status:     models.UntreatedApply,
		LaunchTime: time.Now().Unix(),
	}
	res := fa.ValidationField()
	if res != "" {
		return res
	}

	err := models.CreateFriendApply(fa)

	if err != nil {
		console.StdLog.Error(err)
		return "远程服务器错误"
	}

	return ""
}

func DelFriend() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("user").(*auth.MainUser)
		uid := ctx.PostForm("uuid")
		if uid == "" {
			renderJSON(ctx, struct{}{}, 1, "uid 不能为空")
			return
		}

		err := models.DelRelation(convert.StringSort(user.Uid, uid))
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, struct{}{}, 1, "远程服务器错误")
			return
		}

		renderJSON(ctx, true)
	}
}
