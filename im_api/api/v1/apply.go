package v1

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"gopkg.in/gin-gonic/gin.v1"
	"tuohai/im_api/models"
	"tuohai/im_api/options"
	"tuohai/internal/auth"
	"tuohai/internal/console"
	"tuohai/internal/convert"
	httplib "tuohai/internal/http"
	"tuohai/internal/pb/IM_Message"
)

func ApplyFriends() gin.HandlerFunc {
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

		list := RenderApplyFriends(apply, token)
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

func UnApplyFriends() gin.HandlerFunc {
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

		list := RenderApplyFriends(apply, token)
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

func RenderApplyFriends(apply []models.FriendApply, token string) []gin.H {
	var list []gin.H
	for i, _ := range apply {
		users, err := auth.GetBatchUsers(token, options.Opts.AuthHost, []string{fmt.Sprintf("user_ids=%s", apply[i].ApplyUid)})
		if err != nil {
			console.StdLog.Error(err)
			continue
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

func ProcessApplyFriend() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		main_user := ctx.MustGet("user").(*auth.MainUser)
		//好友申请表中的唯一标识
		id := ctx.PostForm("id")
		status := ctx.PostForm("status")
		if status == "" || id == "" {
			renderJSON(ctx, struct{}{}, 1, "id 或 status不能为空")
			return
		}
		status_int := convert.StrTo(status).MustInt()
		if status_int != 1 && status_int != 2 {
			renderJSON(ctx, struct{}{}, 1, "status invalid")
			return
		}
		fa, err := models.FriendApplyById(id, main_user.Uid)
		if err != nil {
			renderJSON(ctx, struct{}{}, 1, "处理的数据不存在")
			return
		}
		fa.Status = models.ApplyType(status_int)
		fa.ConfirmTime = time.Now().Unix()
		// TODO rework
		cid, err := models.ProcessFriendApply(fa)

		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, struct{}{}, 1, "数据处理有误")
			return
		}

		// 发送通知和消息
		if status_int == 1 {
			go FriendAddMsg(cid, fa.ApplyUid, main_user)
			go func() {
				// 更新通讯录匹配信息
				pm := &models.PersonMatched{
					From:      fa.ApplyUid,
					Partner:   fa.TargetUid,
					Status:    13,
					UpdatedAt: time.Now().Unix(),
				}
				err := models.UpdatePersonMatched(pm)
				if err != nil {
					fmt.Printf("update person matched error:%s", err)
				}
			}()
		} else if status_int == 2 {
			go FriendRefuseMsg(cid, fa.ApplyUid, main_user)
			go func() {
				// 更新通讯录匹配信息
				pm := &models.PersonMatched{
					From:      fa.ApplyUid,
					Partner:   fa.TargetUid,
					Status:    12,
					UpdatedAt: time.Now().Unix(),
				}
				err := models.UpdatePersonMatched(pm)
				if err != nil {
					fmt.Printf("update person matched error:%s", err)
				}
			}()
		}
		renderJSON(ctx, true)
	}
}

func FriendRefuseMsg(cid, uid string, user *auth.MainUser) {
	from := user.Uid
	m := &IM_Message.IMMsgData{
		Type:    "event",
		Subtype: "e_friend_refused",
		From:    from,
		RcvId:   uid,
		MsgData: []byte("{\"uid\":\"" + from + "\"}"),
	}
	fmt.Printf("send friend refused event:%s", m)
	httplib.SendLogicMsg(options.Opts.RPCHost, m)

	// 系统消息发送
	type sysmsg struct {
		Content string `json:"c"`
		Title   string `json:"title"`
		Cid     string `json:"cid"`
	}
	sm := &sysmsg{
		Content: fmt.Sprintf("用户【%s】拒绝了您的好友申请", user.Nickname),
		Title:   "好友申请拒绝",
	}
	srid := models.GetSysRid(options.Opts.SysUserYunliao, uid)
	if srid == "" {
		fmt.Printf("system relation no exist,uuid: %s\n", uid)
		return
	}
	gs, err := json.Marshal(sm)
	if err != nil {
		fmt.Printf("json marshal error,err:%s", err)
		return
	}
	m = &IM_Message.IMMsgData{
		Type:    "message",
		Subtype: "m_system",
		From:    options.Opts.SysUserYunliao,
		To:      srid,
		RcvId:   uid,
		MsgData: gs,
	}
	fmt.Printf("send friend refused event:%s", m)
	httplib.SendLogicMsg(options.Opts.RPCHost, m)
}

func FriendAddMsg(cid, uid string, user *auth.MainUser) {
	from := user.Uid
	m := &IM_Message.IMMsgData{
		Type:    "message",
		Subtype: "m_friend_added",
		From:    from,
		To:      cid,
		MsgData: []byte("{\"c\":\"我们已经成为好友了，开始聊天吧\"}"),
	}
	fmt.Printf("send friend added msg:%s", m)
	httplib.SendLogicMsg(options.Opts.RPCHost, m)

	m = &IM_Message.IMMsgData{
		Type:    "event",
		Subtype: "e_friend_added",
		From:    from,
		RcvId:   uid,
		MsgData: []byte("{\"uid\":\"" + from + "\"}"),
	}
	fmt.Printf("send friend added event:%s", m)
	httplib.SendLogicMsg(options.Opts.RPCHost, m)

	// 系统消息发送
	type sysmsg struct {
		Content string `json:"c"`
		Title   string `json:"title"`
		Cid     string `json:"cid"`
		Type    string `json:"type"`
	}
	sm := &sysmsg{
		Content: fmt.Sprintf("用户【%s】通过了您的好友申请，你们可以开始沟通啦", user.Nickname),
		Title:   "好友申请通过",
		Cid:     cid,
		Type:    "new_friend",
	}
	srid := models.GetSysRid(options.Opts.SysUserYunliao, uid)
	if srid == "" {
		fmt.Printf("system relation no exist,uuid: %s\n", uid)
		return
	}
	gs, err := json.Marshal(sm)
	if err != nil {
		fmt.Printf("json marshal error,err:%s", err)
		return
	}
	m = &IM_Message.IMMsgData{
		Type:    "message",
		Subtype: "m_system",
		From:    options.Opts.SysUserYunliao,
		To:      srid,
		RcvId:   uid,
		MsgData: gs,
	}
	fmt.Printf("send friend added event:%s", m)
	httplib.SendLogicMsg(options.Opts.RPCHost, m)
}
