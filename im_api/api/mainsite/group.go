package mainsite

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"strings"

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
		sign := ctx.PostForm("sign")
		ts := ctx.PostForm("ts")
		if !CheckSign(ts, sign) {
			render.RenderJSON(ctx, struct{}{}, 1, "无效的参数")
			return
		}
		console.StdLog.Debug("creator: ", creator, "name:", name, "member: ", member)
		if name == "" {
			render.RenderJSON(ctx, []int{}, 1, "name is empty")
			return
		}

		ProTeamGroup(ctx, creator, name, models.Project_Group, strings.Split(member, ","))
	}
}

func CreateTeamGroup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		creator := ctx.PostForm("creator")
		name := ctx.PostForm("name")
		member := ctx.PostForm("member")
		if len(creator) == 0 || len(name) == 0 || len(member) == 0 {
			render.RenderJSON(ctx, struct{}{}, 1, "缺少参数")
			return
		}
		sign := ctx.PostForm("sign")
		ts := ctx.PostForm("ts")
		if !CheckSign(ts, sign) {
			render.RenderJSON(ctx, struct{}{}, 1, "无效的参数")
			return
		}
		ProTeamGroup(ctx, creator, name, models.Team_Group, strings.Split(member, ","))
	}
}

func ProTeamGroup(ctx *gin.Context, creator, name string, gtype models.GroupType, member []string) {
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
	appid := uuid.NewV4().StringMd5()

	bot_info := gin.H{
		"bot_access_token": bot_access_token,
		"bot_id":           botid,
		"name":             bot_name,
		"app_id":           appid,
		"channel_id":       gid,
		"creator_id":       creator,
	}
	console.StdLog.Debug("boti_info: ", bot_info)
	if err := models.SaveBotInfo("bot:id:"+botid, bot_info); err != nil {
		console.StdLog.Error(err)
		render.RenderJSON(ctx, []int{}, 1, "未找到数据")
		return
	}
	go func() {
		ns := " "
		for _, id := range member {
			ns += ("<@" + id + "> ")
		}
		tip := "<@" + creator + "> 邀请" + ns + "加入群聊"
		gcn := &GroupChangeNotify{
			Uid:  creator,
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
			From:    creator,
			To:      gid,
			MsgData: gg,
		})
	}()
	render.RenderJSON(ctx, gin.H{
		"web_hook":         fmt.Sprintf("%s/v1/hook/%s", options.Opts.WebHookHost, botid),
		"group_id":         gid,
		"bot_access_token": bot_access_token,
		"bot_id":           botid,
	})
}

//退出群成员
func QuitGroupMember(ctx *gin.Context) {
	uid := ctx.PostForm("uid")
	gid := ctx.Param("gid")
	sign := ctx.PostForm("sign")
	ts := ctx.PostForm("ts")
	if !CheckSign(ts, sign) {
		render.RenderJSON(ctx, struct{}{}, 1, "无效的参数")
		return
	}
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
		tip := "<@" + uid + "> 已退出群聊"
		gcn := &GroupChangeNotify{
			Uid:  uid,
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
			From:    uid,
			To:      gid,
			MsgData: gg,
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
	if len(ids) == 0 || len(gid) == 0 || len(uid) == 0 {
		render.RenderJSON(ctx, struct{}{}, 1, "无效的参数")
		return
	}
	sign := ctx.PostForm("sign")
	ts := ctx.PostForm("ts")
	if !CheckSign(ts, sign) {
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
		ns := " "
		for _, id := range ids {
			ns += ("<@" + id + "> ")
		}
		tip := "<@" + uid + "> 邀请" + ns + "加入群聊"
		gcn := &GroupChangeNotify{
			Uid:  uid,
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
			From:    uid,
			To:      g.Gid,
			MsgData: gg,
		})
	}()
	render.RenderJSON(ctx, true)
	return
}

func SendSystemMsg(ctx *gin.Context) {
	from := ctx.PostForm("from")
	to := ctx.PostForm("to")
	msg := ctx.PostForm("msg")
	sign := ctx.PostForm("sign")
	ts := ctx.PostForm("ts")
	fmt.Printf("rcv sysmsg data:from:%s,to:%s,msg:%s\n", from, to, msg)
	if len(from) == 0 || len(to) == 0 || !CheckSign(ts, sign) {
		render.RenderJSON(ctx, struct{}{}, 1, "无效的参数")
		return
	}
	if len(msg) == 0 || len(msg) > 1024 {
		render.RenderJSON(ctx, struct{}{}, 1, "好长的消息")
		return
	}
	rid := models.GetSysRid(to, from)
	if rid == "" {
		render.RenderJSON(ctx, struct{}{}, 1, "无效的好友参数")
		return
	}

	//RPC通知IM
	go func() {
		m := &IM_Message.IMMsgData{
			Type:    "message",
			Subtype: "m_system",
			From:    from,
			To:      rid,
			MsgData: []byte("{\"c\":" + msg + "}"),
		}
		fmt.Printf("send sysmsg %s:", m)
		httplib.SendLogicMsg(options.Opts.RPCHost, m)
	}()
	render.RenderJSON(ctx, true)
	return
}

func CheckSign(ts, sign string) bool {
	ns := fmt.Sprintf("%x", md5.Sum([]byte("clouderworkgots="+ts)))
	return ns == sign
}

type GroupChangeNotify struct {
	Uid  string `json:"uuid"`
	Gid  string `json:"gid"`
	Type string `json:"type"`
	Tip  string `json:"tip"`
}
