package v1

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	// "time"

	"gopkg.in/gin-gonic/gin.v1"
	"tuohai/im_api/models"
	"tuohai/internal/auth"
	"tuohai/internal/console"
	httplib "tuohai/internal/http"
	"tuohai/internal/util"
)

//获取个人信息
//包括头像、昵称、个性签名、头像地址
func Profile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("user").(*auth.MainUser)
		u, err := models.GetUserById(user.Uid)
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, struct{}{}, 0, "未找到数据")
			return
		}
		renderJSON(ctx, u)
	}
}

//完善或者更改个人信息
//头像 上传不带host的url 客户端首先上传文件服务器 服务器返回url
//昵称 更改服务方名称
//个性签名 个性签名保存在im需要保存在本地im数据库中
func PutProfile(url string) gin.HandlerFunc {
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
		nickname := ctx.PostForm("nickname")
		avatar := ctx.PostForm("avatar")
		// desc := ctx.PostForm("desc")
		// u, err := models.GetUserById(user.Uid)
		// if err != nil {
		// 	console.StdLog.Error(err)
		// 	renderJSON(ctx, struct{}{}, 1, "远程服务器错误1")
		// }
		user_id := strconv.FormatInt(user.Id, 10)

		// if desc != "" {
		// 	fmt.Println("desc: ", desc)
		// 	err := models.UpdateUser(u)
		// 	if err != nil {
		// 		console.StdLog.Error(err)
		// 		renderJSON(ctx, struct{}{}, 1, "远程服务器错误2")
		// 		return
		// 	}
		// }
		//本地数据库更新个人信息

		//更新主站用户信息
		auth_url := auth.GetUpdateUserInfoUrl(token, url, []string{
			fmt.Sprintf("nickname=%s", nickname),
			fmt.Sprintf("avatar=%s", avatar),
			fmt.Sprintf("user_id=%s", user_id),
		})
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
		token := ctx.MustGet("user").(*auth.MainUser)
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
			user, err = models.GetUserById(token.Uid)
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
			list = append(list, gin.H{
				"sid": session.Sid,
				"cid": session.To,
				"msg": gin.H{
					"id":       history.MsgId,
					"data":     history.MsgData,
					"type":     history.Type,
					"sub_type": history.Subtype,
					"time":     history.CreatedAt,
				},
				"type": session.SType,
			})
		}
		renderJSON(ctx, list)
	}
}

func Messages() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// sid := ctx.Param("sid")
		record := &models.Msgrecord{}
		msg, err := models.GetMsgById(record)
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
		cid := ctx.Param("cid")
		msgid := ctx.Param("msgid")
		origin := ctx.Param("origin")
		fmt.Printf("get read info.cid:%s,msgid:%s,origin:%s\n", cid, msgid, origin)
		cnt, res, err := models.MsgReadInfo(cid, msgid, origin)
		if err != nil {
			renderJSON(ctx, map[string]string{}, 1, "查数据出现问题")
		} else {
			renderJSON(ctx, map[string]string{
				"cnt":    strconv.Itoa(cnt),
				"read":   strings.Join(res["read"], ","),
				"unread": strings.Join(res["unread"], ","),
			})
		}
		return
	}
}

func Unreads() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func CreateGroup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var group models.Group
		token := ctx.MustGet("token").(string)
		if err := ctx.Bind(&group); err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, []int{}, 0, err)
			return
		}
		if group.Gname == "" {
			renderJSON(ctx, []int{}, 1, "group_name is empty")
			return
		}

		group.Creator = token

		g, err := models.CreateGroup(&group)
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, []int{}, 1, "未找到数据")
			return
		}

		renderJSON(ctx, g)
	}
}

func GroupRename() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		gid := ctx.Param("gid")
		newname := ctx.Param("newname")
		token := ctx.MustGet("token").(string)
		if gid == "" || newname == "" {
			renderJSON(ctx, struct{}{}, 1, "无效的URL参数!")
			return
		}

		//判断操作这否有权限操作群
		if !models.Permit(gid, token, models.RENAME_GROUP).IsEditTitle() {
			//无权限
			renderJSON(ctx, struct{}{}, 1, "无权限更名!")
			return
		}

		if err := models.RenameGroup(gid, newname); err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, struct{}{}, 1, "修改群名称失败!")
			return
		}
		renderJSON(ctx, "ok")
		return
	}
}

//解散群组
func DismissGroup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		gid := ctx.Param("gid")
		token := ctx.MustGet("token").(string)
		if gid == "" {
			renderJSON(ctx, struct{}{}, 1, "无效的URL参数!")
			return
		}

		//判断操作这否有权限操作群
		if !models.Permit(gid, token, models.DISMISS_GROUP).IsDismissGroup() {
			//无权限
			renderJSON(ctx, struct{}{}, 1, "无权解散群!")
			return
		}

		if err := models.DismissGroup(gid); err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, struct{}{}, 1, "解散群失败!")
			return
		}

		renderJSON(ctx, "ok")
		return
	}
}

//添加群成员
func AddGroupMember() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ids := strings.Split(ctx.Request.FormValue("ids"), ",")
		gid := ctx.Param("gid")
		token := ctx.MustGet("token").(string)
		if len(ids) == 0 {
			renderJSON(ctx, struct{}{}, 1, "无效的参数")
			return
		}

		//判断操作这否有权限操作群
		if !models.Permit(gid, token, models.ADD_GROUP_MEMS).IsAddGroupMember() {
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
				renderJSON(ctx, struct{}{}, 1)
			}
			return
		}

		renderJSON(ctx, g)
		return
	}
}

//移除群成员
func RemoveGroupMember() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ids := strings.Split(ctx.Request.FormValue("ids"), ",")
		gid := ctx.Param("gid")
		token := ctx.MustGet("token").(string)
		if len(ids) == 0 {
			renderJSON(ctx, struct{}{}, 1, "无效的参数")
			return
		}

		//判断操作这否有权限操作群
		if !models.Permit(gid, token, models.ADD_GROUP_MEMS).IsRemoveGroupMember() {
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

		renderJSON(ctx, g)
		return
	}
}

//退出群成员
func QuitGroupMember() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.MustGet("token").(string)
		gid := ctx.Param("gid")

		g, err := models.DelGroupMember(gid, []string{token})
		if err != nil {
			if err == models.RecordNotFound {
				console.StdLog.Error(err)
				renderJSON(ctx, struct{}{}, 0)
			} else {
				renderJSON(ctx, struct{}{}, 1)
			}
			return
		}

		renderJSON(ctx, g)
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

			fuser, err := models.GetUserById(f_uuid)
			if err != nil {
				console.StdLog.Error(err)
			}
			list = append(list, gin.H{
				"f_name": fuser.Uname,
				"f_uuid": fuser.Uuid,
				"rid":    rel.Rid,
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

		fuser, err := models.GetUserById(f_uuid)
		if err != nil {
			console.StdLog.Error(err)
		}

		renderJSON(ctx, gin.H{
			"f_name": fuser.Uname,
			"f_uuid": fuser.Uuid,
			"rid":    rel.Rid,
		})
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
