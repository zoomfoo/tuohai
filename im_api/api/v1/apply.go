package v1

import (
	"fmt"
	"strconv"
	"time"

	"gopkg.in/gin-gonic/gin.v1"
	"tuohai/im_api/models"
	"tuohai/internal/auth"
	"tuohai/internal/console"
	"tuohai/internal/convert"
	httplib "tuohai/internal/http"
	"tuohai/internal/pb/IM_Message"
)

func ApplyFriends(url string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		main_user := ctx.MustGet("user").(*auth.MainUser)
		token := ctx.MustGet("token").(string)
		pageindex, _ := strconv.Atoi(ctx.Param("pageindex"))
		pagesize, _ := strconv.Atoi(ctx.Param("pagesize"))

		apply, err := models.FriendApplys(main_user.Uid, true, pageindex, pagesize)
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, struct{}{}, 1, "远程服务器错误")
			return
		}

		list := RenderApplyFriends(apply, token, url)
		if len(list) == 0 {
			renderJSON(ctx, []int{})
			return
		}

		renderJSON(ctx, gin.H{
			"list":  list,
			"total": models.FriendApplysCount(main_user.Uid, true),
		})
	}
}

func UnApplyFriends(url string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		main_user := ctx.MustGet("user").(*auth.MainUser)
		token := ctx.MustGet("token").(string)
		pageindex, _ := strconv.Atoi(ctx.Param("pageindex"))
		pagesize, _ := strconv.Atoi(ctx.Param("pagesize"))

		apply, err := models.FriendApplys(main_user.Uid, false, pageindex, pagesize)
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, struct{}{}, 1, "远程服务器错误")
			return
		}

		list := RenderApplyFriends(apply, token, url)
		if len(list) == 0 {
			renderJSON(ctx, []int{})
			return
		}

		renderJSON(ctx, gin.H{
			"list":  list,
			"total": models.FriendApplysCount(main_user.Uid, false),
		})
	}
}

func RenderApplyFriends(apply []models.FriendApply, token, url string) []gin.H {
	var list []gin.H
	for i, _ := range apply {
		users, err := auth.GetBatchUsers(token, url, []string{fmt.Sprintf("user_ids=%s", apply[i].ApplyUid)})
		if err != nil {
			console.StdLog.Error(err)
		}
		avatar, name := "", ""
		if len(users) != 0 {
			name = users[0].Uname
			avatar = users[0].Avatar
		}

		list = append(list, gin.H{
			"id":     apply[i].Fid,
			"uuid":   apply[i].ApplyUid,
			"way":    apply[i].Way,
			"attach": apply[i].Attach,
			"status": apply[i].Status,
			"avatar": avatar,
			"name":   name,
			"time":   apply[i].LaunchTime,
		})
	}
	return list
}

func AgreeApplyFriend(addr string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		main_user := ctx.MustGet("user").(*auth.MainUser)
		id := ctx.PostForm("id")
		status := ctx.PostForm("status")
		if status == "" || id == "" {
			renderJSON(ctx, struct{}{}, 1, "id 或 status不能为空")
			return
		}
		status_int := convert.StrTo(status).MustInt()
		cid, err := models.SaveFriendApply(&models.FriendApply{
			Fid: id, ApplyUid: id,
			TargetUid: main_user.Uid,
			Status:    models.ApplyType(status_int),
		})

		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, struct{}{}, 1, "远程服务器错误")
			return
		}

		//获取好友信息cid
		go FriendToLogic(addr, main_user.Uid, cid)
		renderJSON(ctx, true)
	}
}

func FriendToLogic(addr, from, cid string) {
	m := &IM_Message.IMMsgData{
		Type:       "message",
		Subtype:    "m_friend_added",
		From:       from,
		MsgData:    []byte("我们已经成为好友"),
		RcvId:      cid,
		CreateTime: strconv.Itoa(int(time.Now().Unix())),
	}
	httplib.SendLogicMsg(addr, m)
}
