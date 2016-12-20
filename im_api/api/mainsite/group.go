package mainsite

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"gopkg.in/gin-gonic/gin.v1"
	"tuohai/im_api/models"
	"tuohai/im_api/options"
	"tuohai/im_api/render"
	"tuohai/internal/console"
	httplib "tuohai/internal/http"
	"tuohai/internal/pb/IM_Message"
	"tuohai/internal/uuid"
)

func CreateProjectGroup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		creator := ctx.PostForm("creator")
		name := ctx.PostForm("name")
		member := ctx.PostForm("member")
		console.StdLog.Debug("creator: ", creator, "name:", name, "member: ", member)
		if name == "" {
			render.RenderJSON(ctx, []int{}, 1, "name is empty")
			return
		}

		ProTeaGroup(ctx, creator, name, models.Project_Group, strings.Split(member, ","))
	}
}

func CreateTeamGroup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		creator := ctx.PostForm("creator")
		name := ctx.PostForm("name")
		member := ctx.PostForm("member")
		ProTeaGroup(ctx, creator, name, models.Team_Group, strings.Split(member, ","))
	}
}

func ProTeaGroup(ctx *gin.Context, creator, name string, gtype models.GroupType, member []string) {
	g, err := models.CreateGroup(creator, name, gtype, member)
	if err != nil {
		console.StdLog.Error(err)
		render.RenderJSON(ctx, []int{}, 1, "未找到数据")
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
		"name":             bot_name,
		"app_id":           appid,
		"channel_id":       gid,
	}
	console.StdLog.Debug("boti_info: ", bot_info)
	if err := models.SaveBotInfo("bot:id:"+botid, bot_info); err != nil {
		console.StdLog.Error(err)
		render.RenderJSON(ctx, []int{}, 1, "未找到数据")
		return
	}

	render.RenderJSON(ctx, gin.H{
		"web_hook":         fmt.Sprintf("%s/hook/%s", options.Opts.WebHookHost, botid),
		"group_id":         gid,
		"bot_access_token": bot_access_token,
		"bot_id":           botid,
	})
}

//退出群成员
func QuitGroupMember(ctx *gin.Context) {
	uid := ctx.PostForm("uid")
	gid := ctx.Param("gid")

	_, err := models.DelGroupMember(gid, []string{uid})
	if err != nil {
		if err == models.RecordNotFound {
			console.StdLog.Error(err)
			render.RenderJSON(ctx, struct{}{}, 0)
		} else {
			render.RenderJSON(ctx, struct{}{}, 1)
		}
		return
	}

	go func() {
		//RPC通知IM
		httplib.SendLogicMsg(options.Opts.RPCHost, &IM_Message.IMMsgData{
			Type:       "message",
			Subtype:    "m_group_changed",
			From:       uid,
			To:         gid,
			MsgData:    []byte("quitmember"),
			CreateTime: strconv.Itoa(int(time.Now().Unix())),
		})
	}()
	render.RenderJSON(ctx, true)
	return
}

//添加群成员
func AddGroupMember(ctx *gin.Context) {
	ids := strings.Split(ctx.Request.FormValue("member"), ",")
	gid := ctx.Param("gid")
	uid := ctx.PostForm("uid")
	if len(ids) == 0 {
		render.RenderJSON(ctx, struct{}{}, 1, "无效的参数")
		return
	}

	g, err := models.AddGroupMember(gid, ids)
	if err != nil {
		if err == models.RecordNotFound {
			console.StdLog.Error(err)
		} else {
			console.StdLog.Error(err)
		}
		render.RenderJSON(ctx, true)
		return
	}

	//RPC通知IM
	go func() {
		httplib.SendLogicMsg(options.Opts.RPCHost, &IM_Message.IMMsgData{
			Type:       "message",
			Subtype:    "m_group_changed",
			From:       uid,
			To:         g.Gid,
			MsgData:    []byte("addmember"),
			CreateTime: strconv.Itoa(int(time.Now().Unix())),
		})
	}()
	render.RenderJSON(ctx, true)
	return
}
