package mainsite

import (
	"fmt"
	"strings"

	"gopkg.in/gin-gonic/gin.v1"
	"tuohai/im_api/models"
	"tuohai/im_api/options"
	"tuohai/im_api/render"
	"tuohai/internal/console"
	"tuohai/internal/uuid"
)

func CreateProjectGroup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		creator := ctx.PostForm("creator")
		name := ctx.PostForm("name")
		member := ctx.PostForm("member")

		if name == "" {
			render.RenderJSON(ctx, []int{}, 1, "group_name is empty")
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
		"bot_name":         bot_name,
		"app_id":           appid,
		"cid":              gid,
	}

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
