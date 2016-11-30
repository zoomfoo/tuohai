package v1

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gopkg.in/gin-gonic/gin.v1"
	msgsender "tuohai/internal/http"
	"tuohai/internal/pb/IM_Message"
	// "tuohai/internal/console"
	// "tuohai/internal/util"
	"tuohai/models"
)

const addr = "127.0.0.1:5004"

func ConfirmChuo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		chuoid, ok := ctx.GetPostForm("chuoid")
		if !ok {
			renderJSON(ctx, struct{}{}, 0, "chuoid缺失")
			return
		}
		rcv, _ := ctx.GetPostForm("rcv")
		fmt.Println(chuoid, rcv)
		if err := models.ConfirmChuo(chuoid, rcv); err != nil {
			renderJSON(ctx, struct{}{}, 0, "数据有误")
			return
		}
		renderJSON(ctx, struct{}{})
		return
	}
}

// 戳一下业务处理
func AddChuo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sender, _ := ctx.GetPostForm("sender")
		cid, _ := ctx.GetPostForm("cid")
		msg_id, _ := ctx.GetPostForm("msg_id")
		tos, ok := ctx.GetPostForm("to")
		if !ok {
			renderJSON(ctx, struct{}{}, 0, "必须提供to,以逗号分隔")
			return
		}
		to := strings.Split(tos, ",")
		if len(to) > 250 {
			renderJSON(ctx, struct{}{}, 0, "戳的人不能超过250")
			return
		}
		content, ok := ctx.GetPostForm("content")
		if !ok {
			renderJSON(ctx, struct{}{}, 0, "必须提供content")
			return
		}
		ur, ok := ctx.GetPostForm("urgent")
		if !ok {
			ur = "1"
		}
		urgent, err := strconv.Atoi(ur)
		if err != nil {
			renderJSON(ctx, struct{}{}, 0, "urgent参数非法")
			return
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
			renderJSON(ctx, struct{}{}, 0, "处理错误")
			return
		}
		go sendChuoEvent(t, to)
		renderJSON(ctx, struct{}{})
		return
	}
}

func sendChuoEvent(t *models.TblChuoyixiaMeta, tos []string) error {
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
		Type:       "event",
		Subtype:    "e_chuo_rcv",
		From:       t.Sender,
		MsgData:    b,
		CreateTime: strconv.Itoa(t.CreatedAt),
	}
	for _, to := range tos {
		m.RcvId = to
		// 可以改为异步
		msgsender.SendLogicMsg(addr, m)
	}
	return nil
}

// 获取戳列表：我发出的
func GetChuoListFrom() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid := ctx.Param("uid")
		clist, err := models.GetChuoFrom(uid)
		if err != nil {
			renderJSON(ctx, struct{}{}, "query error")
			return
		}

		renderJSON(ctx, clist)
	}
}

// 获取戳列表：我收到的
func GetChuoListRcv() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid := ctx.Param("uid")
		clist, err := models.GetChuoRcv(uid)
		if err != nil {
			renderJSON(ctx, struct{}{}, "query error")
			return
		}

		renderJSON(ctx, clist)
	}
}

// 获取戳详情
func GetChuoInfo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		chuoid := ctx.Param("chuoid")
		clist, err := models.GetChuo(chuoid)
		if err != nil {
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
