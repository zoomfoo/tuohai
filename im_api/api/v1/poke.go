package v1

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gopkg.in/gin-gonic/gin.v1"
	"tuohai/im_api/models"
	"tuohai/im_api/options"
	"tuohai/internal/auth"
	"tuohai/internal/console"
	"tuohai/internal/convert"
	httplib "tuohai/internal/http"
	"tuohai/internal/pb/IM_Message"
)

func ConfirmChuo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		chuoid := ctx.Param("pid")
		user := ctx.MustGet("user").(*auth.MainUser)
		rcv := user.Uid
		chuo, err := models.GetChuoMeta(chuoid)
		if err != nil {
			renderJSON(ctx, struct{}{}, 0, "pid is invalid")
			return
		}
		if err := models.ConfirmChuo(chuoid, rcv); err != nil {
			renderJSON(ctx, struct{}{}, 0)
			return
		}
		type confirm struct {
			Chuoid    string `json:"chuoid"`
			Confirmer string `json:"confirmer"`
		}

		go func() {
			sc := &confirm{
				Chuoid:    chuoid,
				Confirmer: user.Uid,
			}
			b, err := json.Marshal(sc)
			if err != nil {
				// log here
				return
			}
			m := &IM_Message.IMMsgData{
				Type:    "event",
				Subtype: "e_chuo_confirmed",
				From:    user.Uid,
				RcvId:   chuo.Sender,
				MsgData: b,
			}
			httplib.SendLogicMsg(options.Opts.RPCHost, m)
		}()
		renderJSON(ctx, true)
		return
	}
}

// 戳一下业务处理
func AddChuo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// sender, _ := ctx.GetPostForm("sender") //戳一下本人
		token := ctx.MustGet("token").(string)
		user := ctx.MustGet("user").(*auth.MainUser)
		sender := user.Uid
		cid, _ := ctx.GetPostForm("cid")
		msg_id, _ := ctx.GetPostForm("mid")
		tos, ok := ctx.GetPostForm("cnee")
		if !ok {
			renderJSON(ctx, struct{}{}, 0, "必须提供to,以逗号分隔")
			return
		}
		to := strings.Split(tos, ",")
		if len(to) > 50 {
			renderJSON(ctx, struct{}{}, 0, "戳的人不能超过50")
			return
		}
		content, ok := ctx.GetPostForm("content")
		if !ok {
			renderJSON(ctx, struct{}{}, 0, "必须提供content")
			return
		}
		ur, ok := ctx.GetPostForm("urgent")
		if !ok {
			ur = "0"
		}
		urgent, err := strconv.Atoi(ur)
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, struct{}{}, 0, "urgent参数非法")
			return
		}

		us, err := auth.GetBatchUsers(token, options.Opts.AuthHost, []string{fmt.Sprintf("user_ids=%s", tos)})
		if err != nil {
			fmt.Println(err)
			renderJSON(ctx, struct{}{}, 0, "获取用户手机失败!")
			return
		}

		var str_phones string
		for i, _ := range us {
			str_phones += us[i].Phone + ","
		}

		now := time.Now().Unix()
		t := &models.TblChuoyixiaMeta{
			Sender:    sender,
			Cid:       cid,
			MsgId:     msg_id,
			Content:   content,
			TotalCnt:  int8(len(to)),
			Urgent:    int8(urgent),
			CreatedAt: int(now),
			UpdatedAt: int(now),
		}
		t.Chuoid = genChuoid(t)
		if err := models.AddChuo(t, to); err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, struct{}{}, 0, "处理错误")
			return
		}
		go sendChuoMsg(t, to)
		switch urgent {
		case 1:
			//发送短信
			go func() {
				cs := fmt.Sprintf("【云聊】%s戳了您一下：%s  请尽快登录云聊确认回复", user.Nickname, content)
                fmt.Println("sms send: ",cs)
				auth.SendSMS(options.Opts.AuthHost, token, []string{
					"phones=" + str_phones[:len(str_phones)-1],
					"content=" + cs,
					"site=yunliao",
					"user_id=" + convert.ToStr(user.Id),
				})
			}()
		case 2:
			//发送电话
		}
		renderJSON(ctx, true)
		return
	}
}

func sendChuoMsg(t *models.TblChuoyixiaMeta, tos []string) error {
	type chuo struct {
		Sender  string `json:"sender"`
		Cid     string `json:"cid"`
		MsgId   string `json:"msg_id"`
		Content string `json:"content"`
		Chuoid  string `json:"chuoid"`
	}
	sc := &chuo{
		Sender:  t.Sender,
		Cid:     t.Cid,
		MsgId:   t.MsgId,
		Content: t.Content,
		Chuoid:  t.Chuoid,
	}
	b, err := json.Marshal(sc)
	if err != nil {
		// log here
		return err
	}
	m := &IM_Message.IMMsgData{
		Type:    "event",
		Subtype: "e_chuo_rcv",
		From:    t.Sender,
		MsgData: b,
	}
	for _, to := range tos {
		m.RcvId = to
		// 可以改为异步
		httplib.SendLogicMsg(options.Opts.RPCHost, m)
	}
	return nil
}

// 获取戳列表：我发出的
func GetChuoListFrom() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("user").(*auth.MainUser)
		token := ctx.MustGet("token").(string)
		uid := user.Uid

		clist, err := models.GetChuoFrom(uid)
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, struct{}{}, "query error")
			return
		}

		var list []gin.H
		for i, _ := range clist {
			u, err := auth.GetBatchUsers(token, options.Opts.AuthHost, []string{fmt.Sprintf("user_ids=%s", clist[i].Sender)})
			name := ""
			if err == nil && len(u) > 0 {
				name = u[0].Uname
			}
			list = append(list, gin.H{
				"poke_id": clist[i].Chuoid,
				"sender":  clist[i].Sender,
				"urgent":  clist[i].Urgent,
				"total":   clist[i].TotalCnt,
				"remain":  clist[i].ConfirmedCnt,
				"cid":     clist[i].Cid,
				"mid":     clist[i].MsgId,
				"content": clist[i].Content,
				"time":    clist[i].CreatedAt,
				"name":    name,
			})
		}

		if len(list) == 0 {
			renderJSON(ctx, []int{})
			return
		}
		renderJSON(ctx, list)
	}
}

// 获取戳列表：我收到的
func GetChuoListRcv() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// uid := ctx.Param("uid")
		user := ctx.MustGet("user").(*auth.MainUser)
		token := ctx.MustGet("token").(string)
		uid := user.Uid
		clist, err := models.GetChuoRcv(uid)
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, struct{}{}, "query error")
			return
		}

		var list []gin.H
		for i, _ := range clist {
			u, err := auth.GetBatchUsers(token, options.Opts.AuthHost, []string{fmt.Sprintf("user_ids=%s", clist[i].Sender)})
			name := ""
			if err == nil && len(u) > 0 {
				name = u[0].Uname
			}
			t := models.GetChuoByUidAndPid(clist[i].Chuoid, user.Uid)

			list = append(list, gin.H{
				"poke_id":   clist[i].Chuoid,
				"sender":    clist[i].Sender,
				"urgent":    clist[i].Urgent,
				"total":     clist[i].TotalCnt,
				"remain":    clist[i].ConfirmedCnt,
				"cid":       clist[i].Cid,
				"mid":       clist[i].MsgId,
				"content":   clist[i].Content,
				"time":      clist[i].CreatedAt,
				"name":      name,
				"confirmed": t.IsConfirmed,
			})
		}

		if len(list) == 0 {
			renderJSON(ctx, []int{})
			return
		}
		renderJSON(ctx, list)
	}
}

// 获取戳详情
func GetChuoInfo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		chuoid := ctx.Param("pid")
		clist, err := models.GetChuo(chuoid)
		if err != nil {
			console.StdLog.Error(err)
			renderJSON(ctx, struct{}{}, "query error")
			return
		}
		con := []string{}
		uncon := []string{}
		for _, x := range clist {
			if x.IsConfirmed == 0 {
				uncon = append(uncon, x.Rcv)
			} else {
				con = append(con, x.Rcv)
			}
		}
		xd := map[string][]string{
			"con":   con,
			"uncon": uncon,
		}
		renderJSON(ctx, xd)
	}
}

func genChuoid(t *models.TblChuoyixiaMeta) string {
	s, err := json.Marshal(t)
	if err != nil {
		return fmt.Sprintf("%s", md5.Sum([]byte(time.Now().String())))
	}
	return fmt.Sprintf("%x", md5.Sum(s))
}
