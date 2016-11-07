package v1

import (
	// "bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gopkg.in/gin-gonic/gin.v1"
	"tuohai/internal/console"
	"tuohai/internal/convert"
	httplib "tuohai/internal/http"
	"tuohai/internal/pb/IM_Message"
	"tuohai/internal/uuid"
	"tuohai/models"
)

func BotList(api_host string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.MustGet("token").(string)
		u := api_host + "/v1/groups?session_token=" + token
		gs, err := httplib.Groups(u)
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, []int{}, 1, "远程服务器错误")
			return
		}

		var id []string
		for _, g := range gs {
			id = append(id, g.Gid)
		}
		bots, err := models.GetBots(id)
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, []int{}, 1, "远程服务器错误")
			return
		}
		renderJSON(ctx, bots)
	}
}

func Apps() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apps, err := models.Apps()
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, []int{}, 1, "远程服务器错误")
			return
		}
		renderJSON(ctx, apps)
	}
}

func CreateBot() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var bot models.Bot
		if err := ctx.Bind(&bot); err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, struct{}{}, 1, err)
			return
		}

		err := models.CreateBot(&models.Bot{
			Id:         uuid.NewV4().String(),
			Name:       bot.Name,
			Icon:       bot.Icon,
			CreatorId:  bot.CreatorId,
			ChannelId:  bot.ChannelId,
			AppId:      bot.AppId,
			State:      1,
			CreateTime: time.Now(),
			UpTime:     time.Now(),
			IsPub:      bot.IsPub,
		})
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, "no", 1)
		}
		renderJSON(ctx, "ok")
	}
}

func UpdateBot() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func DeleteBot() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func PushMsg() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			bot_id       = ctx.Param("bot_id")
			msg_type     = "message"
			msg_sub_type = "bot_msg"
		)

		bot, err := models.GetBotById(bot_id)
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, struct{}{}, 1, "没有找到数据")
			return
		}

		if bot.Id == "" {
			renderJSON(ctx, struct{}{}, 1, "没有找到数据")
			return
		}

		if bot.ChannelId == "" {
			renderJSON(ctx, struct{}{}, 1, "没有找到数据")
			return
		}

		data, err := ioutil.ReadAll(ctx.Request.Body)
		defer ctx.Request.Body.Close()
		if err != nil {
			renderJSON(ctx, struct{}{}, 1, "读取body失败")
			return
		}

		msg := &IM_Message.IMMsgData{
			Type:       msg_type,
			Subtype:    msg_sub_type,
			From:       bot.Id,
			To:         bot.ChannelId,
			MsgData:    data,
			CreateTime: convert.ToStr(time.Now().Unix()),
		}

		log.Println(*msg)

		httplib.SendLogicMsg("127.0.0.1:9003", msg)
		renderJSON(ctx, "ok")
	}
}

func PushHook() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			bot_id = ctx.Param("bot_id")
		)

		bot, err := models.GetBotById(bot_id)
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, struct{}{}, 1, "没有找到数据")
			return
		}

		app, err := models.GetAppById(bot.AppId)
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, struct{}{}, 1, "没有找到数据")
			return
		}

		log.Println(app.AppURL)
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, struct{}{}, 1, "远程服务器错误")
			return
		}

		bot_info, err := json.Marshal(bot)
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, struct{}{}, 1, "远程服务器错误")
			return
		}

		tarbody, err := ioutil.ReadAll(ctx.Request.Body)
		defer ctx.Request.Body.Close()
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, struct{}{}, 1, "远程服务器错误")
			return
		}

		val := url.Values{
			"bot_info": []string{string(bot_info)},
			"content":  []string{string(tarbody)},
		}

		payload := strings.NewReader(val.Encode())
		req, err := http.NewRequest("POST", app.AppURL, payload)
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, struct{}{}, 1, "远程服务器错误")
			return
		}

		req.Header.Add("content-type", "application/x-www-form-urlencoded")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, struct{}{}, 1, "远程服务器错误")
			return
		}
		defer res.Body.Close()

		if res.StatusCode == http.StatusOK {
			renderJSON(ctx, "ok")
		} else {
			renderJSON(ctx, struct{}{}, 1, res.Status)
		}
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
