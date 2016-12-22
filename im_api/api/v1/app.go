package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"gopkg.in/gin-gonic/gin.v1"
	"tuohai/im_api/models"
	"tuohai/internal/auth"
	"tuohai/internal/console"
	// "tuohai/internal/convert"
	"tuohai/im_api/options"
	httplib "tuohai/internal/http"
	"tuohai/internal/pb/IM_Message"
	"tuohai/internal/util"
)

//获取个人信息
//包括头像、昵称、个性签名、头像地址
func Profile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		main_user := ctx.MustGet("user").(*auth.MainUser)
		u, err := models.GetUserById(main_user.Uid)
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, struct{}{}, 0, "未找到数据")
			return
		}

		renderJSON(ctx, gin.H{
			"username": main_user.Username,
			"uuid":     main_user.Uid,
			"phone":    main_user.Phone,
			"avatar":   main_user.Avatar,
			"name":     main_user.Nickname,
			"email":    main_user.Username,
			"desc":     u.Desc,
		})

	}
}

//完善或者更改个人信息
//头像 上传不带host的url 客户端首先上传文件服务器 服务器返回url
//昵称 更改服务方名称
//个性签名 个性签名保存在im需要保存在本地im数据库中
func PutProfile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var result struct {
			Msg  string `json:"msg"`
			Data struct {
				Nickname string `json:"nickname"`
				Avatar   string `json:"avatar"`
			} `json:"data"`
			ErrorCode float64 `json:"error_code"`
		}

		//主站 /api/i/profile
		token := ctx.MustGet("token").(string)
		user := ctx.MustGet("user").(*auth.MainUser)
		nickname := ctx.PostForm("name")
		avatar := ctx.PostForm("avatar")
		desc := ctx.PostForm("desc")

		if desc != "" {
			fmt.Println("desc: ", desc)
			err := models.SaveUser(&models.User{Uuid: user.Uid, Desc: desc})
			if err != nil {
				console.StdLog.Error(err)
				renderJSON(ctx, struct{}{}, 1, "远程服务器错误2")
				return
			}
		}

		param := []string{fmt.Sprintf("user_id=%d", user.Id)}
		if nickname != "" {
			param = append(param, fmt.Sprintf("nickname=%s", nickname))
		}
		if avatar != "" {
			param = append(param, fmt.Sprintf("avatar=%s", avatar))
		}

		if len(param) < 2 {
			fmt.Println("未提供任何参数!")
			renderJSON(ctx, true)
			return
		}

		//更新主站用户信息
		auth_url := auth.GetUpdateUserInfoUrl(token, options.Opts.AuthHost, param)
		fmt.Println(auth_url)
		if err := httplib.Put(auth_url).ToJson(&result); err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, struct{}{}, 1, "远程服务器错误3")
			return
		}

		if result.ErrorCode == 0 {
			fmt.Println(result)
			renderJSON(ctx, true)
			return
		} else {
			console.StdLog.Error(fmt.Errorf("%s", result.Msg))
			renderJSON(ctx, struct{}{}, 1, result.Msg)
			return
		}
	}
}

//批量获取用户信息
//如果与操作者不是好友cid返回空
func Users() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("user").(*auth.MainUser)
		token := ctx.MustGet("token").(string)
		ids := ctx.Query("ids")
		u, err := auth.GetBatchUsers(token, options.Opts.AuthHost, []string{fmt.Sprintf("user_ids=%s", ids)})
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, struct{}{}, 1, "远程服务器错误")
			return
		}

		var list []gin.H
		for i, _ := range u {
			rel, _ := models.FriendSmallAndBig(user.Uid, u[i].Uuid)
			fmt.Println("获取好友rid", user.Uid, u[i].Uuid)
			list = append(list, gin.H{
				"uuid":   u[i].Uuid,
				"name":   u[i].Uname,
				"phone":  u[i].Phone,
				"email":  u[i].Email,
				"avatar": u[i].Avatar,
				"desc":   u[i].Desc,
				"cid":    rel.Rid,
			})
		}

		if len(list) == 0 {
			renderJSON(ctx, []int{})
			return
		}
		renderJSON(ctx, list)
	}
}

//获取群组列表
func Groups() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("user").(*auth.MainUser)
		mems_groups, err := models.AssociationGroups(user.Uid)
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, []int{}, 0, "未找到数据")
			return
		}

		var groups []models.Group
		for _, mems := range mems_groups {
			group, err := models.GetGroupById(mems.GroupId)
			if err != nil {
				if err != models.RecordNotFound {
					console.StdLog.Error(err)
				}
				continue
			}
			//如果是项目群则过滤掉
			if group.GType == models.Project_Group {
				continue
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

//获取群组信息
func Group() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		gid := ctx.Param("gid")
		main_user := ctx.MustGet("user").(*auth.MainUser)
		token := ctx.MustGet("token").(string)
		wg := &util.WaitGroupWrapper{}

		var (
			group *models.Group
			user  *models.User
			err   error
		)
		wg.Add(2)
		wg.Wrap(func() {
			group, err = models.GetGroupById(gid)
			wg.Done()
		})

		wg.Wrap(func() {
			user, err = models.GetUserById(main_user.Uid)
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

			var list []gin.H
			for _, gm := range group.GroupMems {
				u, err := auth.GetBatchUsers(token, options.Opts.AuthHost, []string{fmt.Sprintf("user_ids=%s", gm)})
				if err != nil {
					fmt.Println(err)
					continue
				}

				name := ""
				avatar := ""
				phone := ""

				if len(u) != 0 {
					name = u[0].Uname
					avatar = u[0].Avatar
					phone = u[0].Phone
				}

				list = append(list, gin.H{
					"uuid":   gm,
					"name":   name,
					"avatar": avatar,
					"phone":  phone,
				})
			}

			renderJSON(ctx, gin.H{
				"gid":     group.Gid,
				"name":    group.Gname,
				"creator": group.Creator,
				"time":    group.CreatedTime,
				"type":    group.GType,
				"member":  list,
			})

		} else {
			renderJSON(ctx, struct{}{}, 1, "当前用户不属于这个群")
		}
	}
}

//获取团队列表
func Teams() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("user").(*auth.MainUser)
		mems_groups, err := models.AssociationGroups(user.Uid)
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, []int{}, 0, "未找到数据")
			return
		}

		var groups []models.Group
		for _, mems := range mems_groups {
			group, err := models.GetGroupById(mems.GroupId)
			if err != nil {
				if err != models.RecordNotFound {
					console.StdLog.Error(err)
				}
				continue
			}
			if group.GType == models.Team_Group {
				groups = append(groups, *group)
			}
		}

		if len(groups) == 0 {
			renderJSON(ctx, []int{})
			return
		}

		renderJSON(ctx, groups)
	}
}

//获取session列表
func Sessions() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("user").(*auth.MainUser)
		sessions, err := models.GetSessionById(user.Uid)
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
			// list = append(list, gin.H{
			// 	"sid": session.Sid,
			// 	"cid": session.To,
			// 	"msg": gin.H{
			// 		"msg_id":      history.MsgId,
			// 		"msg_data":    history.MsgData,
			// 		"type":        history.Type,
			// 		"sub_type":    history.Subtype,
			// 		"create_time": history.CreatedAt,
			// 	},
			// 	"type":       session.SType,
			// 	"unread_cnt": models.ChennelUnreadNum(session.To, user.Uid),
			// })
			list = append(list, gin.H{
				"sid":        session.Sid,
				"cid":        session.To,
				"msg":        history,
				"type":       session.SType,
				"unread_cnt": models.ChennelUnreadNum(session.To, user.Uid),
			})
		}
		if len(list) == 0 {
			renderJSON(ctx, []int{})
			return
		}
		renderJSON(ctx, list)
	}
}

func MsgHistory() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// start := ctx.Query("start")
		// end := ctx.Query("end")
		// cid := ctx.Query("cid")
		renderJSON(ctx, []int{})
		return
	}
}

func ForwardMsg() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		renderJSON(ctx, []int{})
		return
	}
}

func Messages() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cid := ctx.Param("cid")
		size := ctx.Query("size")
		mid := ctx.Query("mid")
		if mid == "" {
			mid = "0"
		}
		if size == "" {
			renderJSON(ctx, []int{}, 0, "size 不能为空!")
			return
		}
		if size > "20" {
			size = "20"
		}
		msg, err := models.GetMsgById(cid, mid, size)
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, []int{}, 0, "未找到数据")
			return
		}

		var messages []gin.H
		user := ctx.MustGet("user").(*auth.MainUser)
		for i, _ := range msg {
			cnt := models.MsgUnreadCount(msg[i].To, strconv.Itoa(int(msg[i].MsgId)), msg[i].From)
			if user.Uid == msg[i].From {
				messages = append(messages, gin.H{
					"from":        msg[i].From,
					"cid":         msg[i].To,
					"type":        msg[i].Type,
					"subtype":     msg[i].Subtype,
					"msg_id":      msg[i].MsgId,
					"msg_data":    msg[i].MsgData,
					"create_time": msg[i].CreatedAt,
					"unread_cnt":  cnt,
				})
			} else {
				is_read := 0
				if cnt > 0 {
					rlist, err := models.MsgReadList(msg[i].To, strconv.Itoa(int(msg[i].MsgId)), msg[i].From)
					if err == nil {
						for i, _ := range rlist {
							if user.Uid == rlist[i] {
								is_read = 1
							}
						}
					}
				} else {
					is_read = 1
				}
				messages = append(messages, gin.H{
					"from":        msg[i].From,
					"cid":         msg[i].To,
					"type":        msg[i].Type,
					"subtype":     msg[i].Subtype,
					"msg_id":      msg[i].MsgId,
					"msg_data":    msg[i].MsgData,
					"create_time": msg[i].CreatedAt,
					"is_read":     is_read,
				})
			}
		}
		if len(messages) == 0 {
			renderJSON(ctx, messages)
			return
		}
		renderJSON(ctx, messages)
	}
}

func RemoveSession() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sid := ctx.Param("sid")
		user := ctx.MustGet("user").(*auth.MainUser)
		fmt.Println("移除sessionid: ", sid)
		//清除session
		if err := models.RemoveSession(sid, user.Uid); err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, false)
			return
		}

		//删除未读消息数
		models.CleanSessionUnread(sid, user.Uid)
		renderJSON(ctx, true)
	}
}

// 创建临时会话
func CreateTmpSession() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// user := ctx.MustGet("user").(*auth.MainUser)
		renderJSON(ctx, struct{}{})
	}
}

// 屏蔽或解除屏蔽临时会话
func ShieldSession(flag bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		renderJSON(ctx, struct{}{})
	}
}

func CleanSessionUnread() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("user").(*auth.MainUser)
		cid := ctx.PostForm("cid")
		uid := user.Uid
		if models.CleanSessionUnread(cid, uid) {
			renderJSON(ctx, true)
		} else {
			renderJSON(ctx, false)
		}
	}
}

func MessageRead() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cid := ctx.Param("cid")
		msgid := ctx.Param("msgid")
		user := ctx.MustGet("user").(*auth.MainUser)
		origin := ctx.Param("origin")
		fmt.Printf("get read info.cid:%s,msgid:%s,origin:%s\n", cid, msgid, user.Uid)
		cnt, res, err := models.MsgReadInfo(cid, msgid, origin)
		if err != nil {
			renderJSON(ctx, struct{}{}, 1, "查数据出现问题")
		} else {
			renderJSON(ctx, gin.H{
				"cnt":    strconv.Itoa(cnt),
				"read":   res["read"],
				"unread": res["unread"],
			})
		}
		return
	}
}

func Unreads() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cid := ctx.Param("cid")
		user := ctx.MustGet("user").(*auth.MainUser)
		if models.CleanSessionUnread(cid, user.Uid) {
			renderJSON(ctx, true)
		} else {
			renderJSON(ctx, false)
		}
		return
	}
}

type GroupChangeNotify struct {
	Uid  string `json:"uid"`
	Gid  string `json:"gid"`
	Type string `json:"type"`
	Tip  string `json:"tip"`
}

//创建普通群
//models.NORMAL_GROUP
func CreateGroup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("user").(*auth.MainUser)
		name := ctx.PostForm("name")
		member := ctx.PostForm("member")
		if name == "" {
			renderJSON(ctx, []int{}, 1, "name is empty")
			return
		}
		if member == "" {
			renderJSON(ctx, []int{}, 1, "群成员不能为空")
			return
		}

		g, err := models.CreateGroup(user.Uid, name, models.NORMAL_GROUP, strings.Split(member, ","))
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, []int{}, 1, "远程服务器错误")
			return
		}
		go func() {
			tip := "<@" + user.Uid + "> 邀请你参加群聊"
			gcn := &GroupChangeNotify{
				Uid:  user.Uid,
				Gid:  g.Gid,
				Type: "create",
				Tip:  tip,
			}
			gg, err := json.Marshal(gcn)
			if err != nil {
				return
			}
			httplib.SendLogicMsg(options.Opts.RPCHost, &IM_Message.IMMsgData{
				Type:    "message",
				Subtype: "m_group_changed",
				From:    user.Uid,
				To:      g.Gid,
				MsgData: gg,
			})
		}()
		renderJSON(ctx, g)
	}
}

func GroupRename() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		gid := ctx.Param("gid")
		newname := ctx.PostForm("name")
		user := ctx.MustGet("user").(*auth.MainUser)
		if gid == "" || newname == "" {
			renderJSON(ctx, struct{}{}, 1, "name 不能为空")
			return
		}

		//判断操作这否有权限操作群
		if !models.Permit(gid, user.Uid, models.RENAME_GROUP).IsEditTitle() {
			//无权限
			renderJSON(ctx, struct{}{}, 1, "无权限更名!")
			return
		}

		if err := models.RenameGroup(gid, newname); err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, struct{}{}, 1, "修改群名称失败!")
			return
		}

		go func() {
			tip := "<@" + user.Uid + "> 将群名修改为：" + newname
			gcn := &GroupChangeNotify{
				Uid:  user.Uid,
				Gid:  gid,
				Type: "rename",
				Tip:  tip,
			}
			gg, err := json.Marshal(gcn)
			if err != nil {
				return
			}
			//RPC通知IM
			httplib.SendLogicMsg(options.Opts.RPCHost, &IM_Message.IMMsgData{
				Type:    "message",
				Subtype: "m_group_changed",
				From:    user.Uid,
				To:      gid,
				MsgData: gg,
			})
		}()
		renderJSON(ctx, true)
		return
	}
}

//解散群组
func DismissGroup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		gid := ctx.Param("gid")
		user := ctx.MustGet("user").(*auth.MainUser)
		if gid == "" {
			renderJSON(ctx, struct{}{}, 1, "无效的URL参数!")
			return
		}

		//判断操作这否有权限操作群
		if !models.Permit(gid, user.Uid, models.DISMISS_GROUP).IsDismissGroup() {
			//无权限
			renderJSON(ctx, struct{}{}, 1, "无权解散群!")
			return
		}

		if err := models.DismissGroup(gid); err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, struct{}{}, 1, "解散群失败!")
			return
		}

		go func() {
			tip := "<@" + user.Uid + "> 已将本群解散"
			gcn := &GroupChangeNotify{
				Uid:  user.Uid,
				Gid:  gid,
				Type: "dismiss",
				Tip:  tip,
			}
			gg, err := json.Marshal(gcn)
			if err != nil {
				return
			}
			//RPC通知IM
			httplib.SendLogicMsg(options.Opts.RPCHost, &IM_Message.IMMsgData{
				Type:    "message",
				Subtype: "m_group_changed",
				From:    user.Uid,
				To:      gid,
				MsgData: gg,
			})
		}()
		renderJSON(ctx, "ok")
		return
	}
}

//添加群成员
func AddGroupMember() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ids := strings.Split(ctx.Request.FormValue("member"), ",")
		gid := ctx.Param("gid")
		user := ctx.MustGet("user").(*auth.MainUser)
		if len(ids) == 0 || len(ids) > 50 {
			renderJSON(ctx, struct{}{}, 1, "无效的参数")
			return
		}

		//判断操作这否有权限操作群
		if !models.Permit(gid, user.Uid, models.ADD_GROUP_MEMS).IsAddGroupMember() {
			//无权限
			renderJSON(ctx, struct{}{}, 1, "无权添加群!")
			return
		}

		g, err := models.AddGroupMember(gid, ids)
		if err != nil {
			if err == models.RecordNotFound {
				console.StdLog.Error(err)
				renderJSON(ctx, struct{}{}, 0)
			} else {
				renderJSON(ctx, struct{}{}, 1, err)
			}
			return
		}

		//RPC通知IM
		go func() {
			ns := " "
			for _, id := range ids {
				ns += ("<@" + id + "> ")
			}
			tip := "<@" + user.Uid + "> 邀请" + ns + "加入群聊"
			gcn := &GroupChangeNotify{
				Uid:  user.Uid,
				Gid:  gid,
				Type: "add",
				Tip:  tip,
			}
			gg, err := json.Marshal(gcn)
			if err != nil {
				return
			}
			httplib.SendLogicMsg(options.Opts.RPCHost, &IM_Message.IMMsgData{
				Type:    "message",
				Subtype: "m_group_changed",
				From:    user.Uid,
				To:      gid,
				MsgData: gg,
			})
		}()
		renderJSON(ctx, g)
		return
	}
}

//移除群成员
func RemoveGroupMember() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ids := strings.Split(ctx.Request.FormValue("member"), ",")
		gid := ctx.Param("gid")
		user := ctx.MustGet("user").(*auth.MainUser)
		if len(ids) == 0 {
			renderJSON(ctx, struct{}{}, 1, "无效的参数")
			return
		}

		//判断操作这否有权限操作群
		if !models.Permit(gid, user.Uid, models.DEL_GROUP_MEMS).IsRemoveGroupMember() {
			//无权限
			renderJSON(ctx, struct{}{}, 1, "无权删除成员!")
			return
		}

		g, err := models.DelGroupMember(gid, ids)
		if err != nil {
			if err == models.RecordNotFound {
				console.StdLog.Error(err)
				renderJSON(ctx, struct{}{}, 0)
			} else {
				renderJSON(ctx, struct{}{}, 1)
			}
			return
		}

		go func() {
			ns := " "
			for _, id := range ids {
				ns += ("<@" + id + "> ")
			}
			tip := "<@" + user.Uid + "> 已将" + ns + "移出本群"
			gcn := &GroupChangeNotify{
				Uid:  user.Uid,
				Gid:  gid,
				Type: "remove",
				Tip:  tip,
			}
			gg, err := json.Marshal(gcn)
			if err != nil {
				return
			}
			//RPC通知IM
			httplib.SendLogicMsg(options.Opts.RPCHost, &IM_Message.IMMsgData{
				Type:    "message",
				Subtype: "m_group_changed",
				From:    user.Uid,
				To:      g.Gid,
				MsgData: gg,
			})
		}()
		renderJSON(ctx, g)
		return
	}
}

//退出群成员
func QuitGroupMember() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("user").(*auth.MainUser)
		gid := ctx.Param("gid")

		_, err := models.DelGroupMember(gid, []string{user.Uid})
		if err != nil {
			if err == models.RecordNotFound {
				console.StdLog.Error(err)
				renderJSON(ctx, struct{}{}, 0)
			} else {
				renderJSON(ctx, struct{}{}, 1)
			}
			return
		}

		go func() {
			tip := "<@" + user.Uid + "> 已退出群聊"
			gcn := &GroupChangeNotify{
				Uid:  user.Uid,
				Gid:  gid,
				Type: "quit",
				Tip:  tip,
			}
			gg, err := json.Marshal(gcn)
			if err != nil {
				return
			}
			//RPC通知IM
			httplib.SendLogicMsg(options.Opts.RPCHost, &IM_Message.IMMsgData{
				Type:    "message",
				Subtype: "m_group_changed",
				From:    user.Uid,
				To:      gid,
				MsgData: gg,
			})
		}()
		renderJSON(ctx, true)
		return
	}
}

func UserInfo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uids := ctx.Param("uid")

		uid := strings.Split(uids, ",")

		user, err := models.GetUserByIds(uid)
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, []int{}, 1, "未找到数据")
			return
		}
		renderJSON(ctx, user)
		return
	}
}

func Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uname := ctx.PostForm("username")
		pwd := ctx.PostForm("password")
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

func GetIds() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			id         = ""
			rels, gros []string
			rerr, gerr error
		)
		wg := &util.WaitGroupWrapper{}
		wg.Add(2)
		wg.Wrap(func() {
			rels, rerr = models.GetMyRelationId(id)
			wg.Done()
		})
		wg.Wrap(func() {
			gros, gerr = models.GetMyGroupId(id)
			wg.Done()
		})
		wg.Wait()

		if rerr != nil || gerr != nil {

		}
		rels = append(rels, gros...)

		renderJSON(ctx, gin.H{
			"ids": rels,
		})
	}
}

func renderJSON(ctx *gin.Context, json interface{}, err_status ...interface{}) {
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

func renderSysError(ctx *gin.Context, err error) {
	if err != nil {
		console.StdLog.Error(err)
		renderJSON(ctx, struct{}{}, 1, "远程服务器错误!")
	}
}
