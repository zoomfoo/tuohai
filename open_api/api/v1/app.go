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

func BotList() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bots, err := models.GetBots()
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
		// var bot models.Bot
		// if err := ctx.Bind(&bot); err != nil {
		// 	console.StdLog.Error(err)
		// 	renderJSON(ctx, struct{}{}, 1, err)
		// 	return
		// }

		// bot.CreateTime = time.Now().Unix()
		// bot.UpTime = time.Now().Unix()

		log.Println(uuid.NewV4())
		renderJSON(ctx, "")
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
			content      = ctx.PostForm("content")
			msg_type     = "message"
			msg_sub_type = "bot_msg"
		)

		if content == "" {
			renderJSON(ctx, struct{}{}, 1, "未找到 content参数。")
			return
		}

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

		msg := &IM_Message.IMMsgData{
			Type:       msg_type,
			Subtype:    msg_sub_type,
			From:       bot.Id,
			To:         bot.ChannelId,
			MsgData:    []byte(content),
			CreateTime: convert.ToStr(time.Now().Unix()),
		}
		log.Println(*msg)
		httplib.SendLogicMsg("127.0.0.1:5004", msg)
		renderJSON(ctx, "ok")
	}
}

func PushHook() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			bot_id = ctx.Param("bot_id")
			o      = new(http.Request)
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
		targetURL, err := url.Parse(app.AppURL)
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

		val := url.Values{"bot_info": []string{string(bot_info)}}
		// buf := &bytes.Buffer{}
		*o = *ctx.Request
		o.URL = targetURL
		o.Method = "POST"

		tarbody, err := ioutil.ReadAll(o.Body)
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, struct{}{}, 1, "远程服务器错误")
			return
		}

		log.Println(ctx.ContentType())
		// switch ctx.ContentType() {
		// case "multipart/form-data":
		// case "application/x-www-form-urlencoded":
		// case "application/xml":
		// 	val.Add("content", string(tarbody))
		// 	o.Body = ioutil.NopCloser(strings.NewReader(val.Encode()))
		// 	log.Println(val.Encode())
		// default:

		// }
		log.Println(string(tarbody))
		val.Add("content", string(tarbody))
		o.Body = ioutil.NopCloser(strings.NewReader(val.Encode()))
		log.Println(val.Encode())

		o.Proto = "HTTP/1.1"
		o.ProtoMajor = 1
		o.ProtoMinor = 1
		o.Close = false
		o.ContentLength = 0
		o.Header.Add("Content-Type", "text/plain")
		transport := http.DefaultTransport
		res, err := transport.RoundTrip(o)
		if err != nil {
			console.StdLog.Errorf("http: proxy error: %v", err)
			renderJSON(ctx, struct{}{}, http.StatusInternalServerError, "回调服务器返回 StatusInternalServerError")
			return
		}

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
