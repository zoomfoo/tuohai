package v1

import (
	"io"
	"log"
	"net/http"
	"net/url"

	"gopkg.in/gin-gonic/gin.v1"
	"time"
	"tuohai/internal/console"
	"tuohai/internal/convert"
	httplib "tuohai/internal/http"
	"tuohai/internal/pb/IM_Message"
	"tuohai/models"
)

func AppList() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apps, err := models.Apps()
		if err != nil {
			console.StdLog.Error(err)
			ctx.JSON(http.StatusOK, map[string]interface{}{"no": "no"})
			return
		}
		ctx.JSON(http.StatusOK, apps)
	}
}

func BotList() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func CreateBot() gin.HandlerFunc {
	return func(ctx *gin.Context) {

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
			log.Println(err)
			return
		}

		*o = *ctx.Request
		o.URL = targetURL
		o.Method = "POST"

		log.Println(*o)
		if q := o.URL.RawQuery; q != "" {
			o.URL.RawPath = o.URL.Path + "?" + q
		} else {
			o.URL.RawPath = o.URL.Path
		}

		o.URL.RawQuery = targetURL.RawQuery

		o.Proto = "HTTP/1.1"
		o.ProtoMajor = 1
		o.ProtoMinor = 1
		o.Close = false

		transport := http.DefaultTransport
		res, err := transport.RoundTrip(o)

		if err != nil {
			log.Printf("http: proxy error: %v", err)
			ctx.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		hdr := ctx.Writer.Header()

		for k, vv := range res.Header {
			for _, v := range vv {
				hdr.Add(k, v)
			}
		}

		ctx.Writer.WriteHeader(res.StatusCode)

		if res.Body != nil {
			io.Copy(ctx.Writer, res.Body)
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
