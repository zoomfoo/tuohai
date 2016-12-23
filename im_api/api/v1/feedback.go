package v1

import (
	"encoding/json"
	"fmt"
	"gopkg.in/gin-gonic/gin.v1"

	"tuohai/im_api/models"
	"tuohai/im_api/options"
	"tuohai/internal/auth"
	"tuohai/internal/console"
	httplib "tuohai/internal/http"
	"tuohai/internal/pb/IM_Message"
)

func Feedback() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		main_user := ctx.MustGet("user").(*auth.MainUser)
		content := ctx.PostForm("content")
		if content == "" {
			renderJSON(ctx, struct{}{}, 1, "The content cannot be empty")
			return
		}
		if err := models.NewFeedback(main_user.Uid, content); err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, struct{}{}, 1, "远程服务器错误!")
			return
		}
		go func() {
			type sysmsg struct {
				Content string `json:"c"`
				Title   string `json:"title"`
				Cid     string `json:"cid"`
			}
			uid := main_user.Uid
			sm := &sysmsg{
				Content: "您提交的意见建议小云已收到，感谢您对云聊的支持与谅解，我们将尽快完善产品修正您提出的问题。",
				Title:   "意见反馈",
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
			m := &IM_Message.IMMsgData{
				Type:    "message",
				Subtype: "m_system",
				From:    options.Opts.SysUserYunliao,
				To:      srid,
				RcvId:   uid,
				MsgData: gs,
			}
			fmt.Printf("send friend refused event:%s", m)
			httplib.SendLogicMsg(options.Opts.RPCHost, m)
		}()
		renderJSON(ctx, true)
		return
	}
}
