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
	// "tuohai/internal/convert"
	httplib "tuohai/internal/http"
	"tuohai/internal/util"
	"tuohai/internal/uuid"
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
		auth_url := auth.GetUpdateUserInfoUrl(token, url, param)
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

func Users(url string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.MustGet("token").(string)
		ids := ctx.Query("ids")
		u, err := auth.GetBatchUsers(token, url, []string{fmt.Sprintf("user_ids=%s", ids)})
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, struct{}{}, 1, "远程服务器错误")
			return
		}
		renderJSON(ctx, u)
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

func Group(url string) gin.HandlerFunc {
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
				u, err := auth.GetBatchUsers(token, url, []string{fmt.Sprintf("user_ids=%s", gm)})
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
				"member":  list,
			})

		} else {
			renderJSON(ctx, struct{}{}, 1, "当前用户不属于这个群")
		}
	}
}

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
			if group.GType == models.Project_Group {
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
		if len(list) == 0 {
			renderJSON(ctx, []int{})
			return
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
		user := ctx.MustGet("user").(*auth.MainUser)
		fmt.Printf("get read info.cid:%s,msgid:%s,origin:%s\n", cid, msgid, user.Uid)
		cnt, res, err := models.MsgReadInfo(cid, msgid, user.Uid)
		if err != nil {
			renderJSON(ctx, struct{}{}, 1, "查数据出现问题")
		} else {
			renderJSON(ctx, gin.H{
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

		g, err := models.CreateGroup(user.Uid, name, strings.Split(member, ","))
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, []int{}, 1, "远程服务器错误")
			return
		}

		renderJSON(ctx, g)
	}
}

func CreateProjectGroup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("user").(*auth.MainUser)
		name := ctx.PostForm("name")
		member := ctx.PostForm("member")

		if name == "" {
			renderJSON(ctx, []int{}, 1, "group_name is empty")
			return
		}

		g, err := models.CreateGroup(user.Uid, name, strings.Split(member, ","))
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, []int{}, 1, "未找到数据")
			return
		}

		//生成botid
		botid := uuid.NewV4().StringMd5()
		gid := g.Gid
		bot_access_token := uuid.NewV4().StringMd5()
		bot_name := "clouderwork"
		appid := "clouderwork"

		bot_info := gin.H{
			"bot_access_token": bot_access_token,
			"bot_id":           botid,
			"bot_name":         bot_name,
			"app_id":           appid,
			"cid":              gid,
		}

		if err := models.SaveBotInfo("bot:id:"+botid, bot_info); err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, []int{}, 1, "未找到数据")
			return
		}

		renderJSON(ctx, gin.H{
			"group": g,
			"bot":   bot_info,
		})

	}
}

func GroupRename() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		gid := ctx.Param("gid")
		newname := ctx.Param("newname")
		user := ctx.MustGet("user").(*auth.MainUser)
		if gid == "" || newname == "" {
			renderJSON(ctx, struct{}{}, 1, "无效的URL参数!")
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
		renderJSON(ctx, "ok")
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

		renderJSON(ctx, "ok")
		return
	}
}

//添加群成员
func AddGroupMember() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ids := strings.Split(ctx.Request.FormValue("ids"), ",")
		gid := ctx.Param("gid")
		user := ctx.MustGet("user").(*auth.MainUser)
		if len(ids) == 0 {
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
		user := ctx.MustGet("user").(*auth.MainUser)
		if len(ids) == 0 {
			renderJSON(ctx, struct{}{}, 1, "无效的参数")
			return
		}

		//判断操作这否有权限操作群
		if !models.Permit(gid, user.Uid, models.ADD_GROUP_MEMS).IsRemoveGroupMember() {
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
		user := ctx.MustGet("user").(*auth.MainUser)
		gid := ctx.Param("gid")

		g, err := models.DelGroupMember(gid, []string{user.Uid})
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
