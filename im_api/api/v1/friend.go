package v1

import (
	"fmt"
	"strings"
	"time"

	"gopkg.in/gin-gonic/gin.v1"
	"tuohai/im_api/models"
	"tuohai/im_api/options"
	"tuohai/internal/auth"
	"tuohai/internal/console"
	"tuohai/internal/convert"
	httplib "tuohai/internal/http"
	"tuohai/internal/pb/IM_Message"
	"tuohai/internal/util"
	"tuohai/internal/uuid"
)

func Friends() gin.HandlerFunc {
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

			u, err := auth.GetBatchUsers(token, options.Opts.AuthHost, []string{fmt.Sprintf("user_ids=%s", f_uuid)})
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

func Friend() gin.HandlerFunc {
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
		u, err := auth.GetBatchUsers(token, options.Opts.AuthHost, []string{fmt.Sprintf("user_ids=%s", id)})
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
			"uuid":     u[0].Uuid,
			"name":     u[0].Uname,
			"cid":      rel.Rid,
			"avatar":   u[0].Avatar,
			"phone":    u[0].Phone,
			"desc":     u[0].Desc,
			"relation": list,
		})

	}
}

func InviteFriend() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		renderJSON(ctx, struct{}{})
		return
	}
}

func AddFriend() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid := ctx.PostForm("uuid")
		attach := ctx.PostForm("attach")
		way := ctx.PostForm("way")
		num := ctx.PostForm("num")
		token := ctx.Query("session_token")

		if uid == "" && num == "" {
			renderJSON(ctx, struct{}{}, 1, "参数缺失")
			return
		}

		if attach == "" {
			renderJSON(ctx, struct{}{}, 1, "附言缺失")
			return
		}
		if way == "" {
			renderJSON(ctx, struct{}{}, 1, "来源方式缺失")
			return
		}

		user := ctx.MustGet("user").(*auth.MainUser)

		//通过uid添加好友
		if uid != "" {
			if user.Uid == uid {
				renderJSON(ctx, struct{}{}, 1, "不允许添加自己为好友")
				return
			}
			// 判断该uuid是否存在
			users, err := auth.GetBatchUsers(token, options.Opts.AuthHost, []string{"user_ids=" + uid})
			if err != nil {
				renderJSON(ctx, []int{}, 1, "查询有误")
				return
			}

			if len(users) == 0 {
				renderJSON(ctx, struct{}{}, 1, "您所添加的好友不是云沃客用户，无法添加好友")
				return
			}
			// 判断两人是否已经是好友
			rid := models.IsRelation(uid, user.Uid, 0)
			if rid != "" {
				renderJSON(ctx, struct{}{}, 1, "已经是好友了")
				return
			}

			res := addFriend(user, uid, way, attach, users[0].Uname)
			if res == "" {
				renderJSON(ctx, true)
				return
			}
			renderJSON(ctx, struct{}{}, 1, res)
			return
		}
		phone, email := "", ""

		//判断是否是手机号
		if util.ValidateMob(num) {
			phone = num
		} else if util.ValidateEmail(num) {
			//判断是否是邮箱
			email = num
		}

		//通过手机号添加好友
		if phone != "" {

			users, err := auth.GetBatchUsers(token, options.Opts.AuthHost, []string{"user_ids=" + phone, "t=phone"})
			if err != nil {
				renderJSON(ctx, []int{}, 1, "查询有误")
				return
			}

			if len(users) == 0 {
				renderJSON(ctx, struct{}{}, 1, "您所添加的好友不是云沃客用户，无法添加好友")
				return
			}
			// 不能添加自己为好友
			if user.Uid == users[0].Uuid {
				renderJSON(ctx, struct{}{}, 1, "不允许添加自己为好友")
				return
			}
			//判断是否已经是好友了
			rid := models.IsRelation(users[0].Uuid, user.Uid, 0)
			if rid != "" {
				renderJSON(ctx, struct{}{}, 1, "已经是好友了")
				return
			}

			res := addFriend(user, users[0].Uuid, way, attach, users[0].Uname)
			if res == "" {
				renderJSON(ctx, true)
				return
			}
			renderJSON(ctx, struct{}{}, 1, res)
			return
		}

		//通过邮箱添加好友
		if email != "" {
			users, err := auth.GetBatchUsers(token, options.Opts.AuthHost, []string{"user_ids=" + email, "t=email"})
			if err != nil {
				console.StdLog.Error(err)
				renderJSON(ctx, []int{}, 1, "查询有误")
				return
			}

			if len(users) == 0 {
				renderJSON(ctx, struct{}{}, 1, "您所添加的好友不是云沃客用户，无法添加好友")
				return
			}
			// 不能添加自己为好友
			if user.Uid == users[0].Uuid {
				renderJSON(ctx, struct{}{}, 1, "不允许添加自己为好友")
				return
			}
			//判断是否已经是好友了
			rid := models.IsRelation(users[0].Uuid, user.Uid, 0)
			if rid != "" {
				renderJSON(ctx, struct{}{}, 1, "已经是好友了")
				return
			}

			res := addFriend(user, users[0].Uuid, way, attach, users[0].Uname)
			if res == "" {
				renderJSON(ctx, true)
				return
			}
			renderJSON(ctx, struct{}{}, 1, res)
			return
		}

		renderJSON(ctx, struct{}{}, 1, "数据处理有误")
		return
	}
}

func addFriend(user *auth.MainUser, uid, way, attach, uname string) string {
	fa := &models.FriendApply{
		Fid:         uuid.NewV4().StringMd5(),
		ApplyUid:    user.Uid,
		TargetUid:   uid,
		Way:         models.ApplyWay(convert.StrTo(way).MustInt()),
		Attach:      attach,
		Status:      models.UntreatedApply,
		LaunchTime:  time.Now().Unix(),
		ConfirmTime: time.Now().Unix(),
	}
	res := fa.ValidationField()
	if res != "" {
		return res
	}

	err := models.CreateFriendApply(fa)

	if err != nil {
		return "远程服务器错误"
	}
	go func() {
		m := &IM_Message.IMMsgData{
			Type:    "event",
			Subtype: "e_friend_apply",
			From:    user.Uid,
			RcvId:   uid,
			MsgData: []byte("{\"name\":\"" + uname + "\"}"),
		}
		fmt.Printf("send friend apply event:%s", m)
		httplib.SendLogicMsg(options.Opts.RPCHost, m)
	}()
	return ""
}

func DelFriend() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("user").(*auth.MainUser)
		cid := ctx.Param("cid")
		if cid == "" {
			renderJSON(ctx, struct{}{}, 1, "cid 不能为空")
			return
		}
		// 是否是好友
		r, err := models.Friend(user.Uid, cid)
		if err != nil {
			renderJSON(ctx, struct{}{}, 1, "cid 非法")
			return
		}
		var f string
		if r.BigId == user.Uid {
			f = r.SmallId
		} else if r.SmallId == user.Uid {
			f = r.BigId
		} else {
			renderJSON(ctx, struct{}{}, 1, "cid 不能为空")
			return
		}
		if f == "" {
			renderJSON(ctx, struct{}{}, 1, "cid 不能为空")
			return
		}

		err = models.DelRelation(cid)
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, struct{}{}, 1)
			return
		}

		//删除session
		if err := models.RemoveSessionByCidAndUid(cid, user.Uid); err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, false)
			return
		}
		go func() {
			m := &IM_Message.IMMsgData{
				Type:    "event",
				Subtype: "e_friend_removed",
				From:    user.Uid,
				RcvId:   f,
				MsgData: []byte("{\"uid\":" + user.Uid + "}"),
			}
			fmt.Printf("send friend apply event:%s", m)
			httplib.SendLogicMsg(options.Opts.RPCHost, m)
		}()
		//聊天记录标记
		renderJSON(ctx, models.CleanSessionUnread(cid, user.Uid))
	}
}

// 创建临时好友
func CreateTmpFriend() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("user").(*auth.MainUser)
		partner := ctx.Query("partner")
		if len(partner) == 0 {
			renderJSON(ctx, struct{}{}, 1, "参数为空")
			return
		}
		u, err := auth.GetBatchUsers(ctx.Query("session_token"), options.Opts.AuthHost, []string{"user_ids=" + partner})
		if err != nil || len(u) == 0 {
			renderJSON(ctx, struct{}{}, 1, "该用户不存在")
			return
		}
		rid, err := models.CreateRelation(user.Uid, partner, 1)
		if err != nil {
			renderJSON(ctx, struct{}{}, 1, "临时好友创建错误")
			return
		}
		r, err := models.Friend(partner, rid)
		if err != nil {
			renderJSON(ctx, struct{}{}, 1, "临时好友创建失败")
			return
		}
		renderJSON(ctx, r)
	}
}

func MatchFriend() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("user").(*auth.MainUser)
		phones := ctx.PostForm("phones")
		if len(phones) == 0 || len(phones) > 2500 {
			renderJSON(ctx, struct{}{}, 1, "手机号参数非法")
			return
		}
		ps := strings.Split(phones, ",")
		if len(ps) > 200 {
			renderJSON(ctx, struct{}{}, 1, "手机号超过了200个")
			return
		}
		params := []string{"t=phone", "user_ids=" + phones}
		users, err := auth.GetBatchUsers(ctx.Query("session_token"), options.Opts.AuthHost, params)
		if err != nil {
			renderJSON(ctx, struct{}{}, 1, "手机号查询失败")
			return
		}
		ret, err := models.MatchFriends(user.Uid, users)
		if err != nil {
			renderJSON(ctx, struct{}{}, 1, "通讯录匹配查询错误")
			return
		}
		renderJSON(ctx, ret)
	}
}

func NewPersons() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("user").(*auth.MainUser)
		st := ctx.Query("session_token")
		ret, err := models.NewPersons(user.Uid, st)
		if err != nil {
			renderJSON(ctx, struct{}{}, "内部处理有误")
			return
		}
		renderJSON(ctx, ret)
	}
}
