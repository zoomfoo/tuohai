package v1

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gopkg.in/gin-gonic/gin.v1"
	// "tuohai/internal/console"
	// "tuohai/internal/util"
	"tuohai/models"
)

// 戳一下业务处理
func AddChuo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sender := ctx.Param("sender")
		cid := ctx.Param("cid")
		msg_id := ctx.Param("msg_id")
		tos := ctx.Param("to")
		to := strings.Split(tos, ",")
		if len(to) > 250 {
			renderJSON(ctx, struct{}{}, 0, "戳的人不能超过250")
			return
		}
		content := ctx.Param("content")
		ur := ctx.Param("urgent")
		urgent, err := strconv.Atoi(ur)
		if err != nil {
			renderJSON(ctx, struct{}{}, 0, "urgent 参数非法")
			return
		}
		now := time.Now().Unix()
		t := &models.TblChuoyixiaMeta{
			From:      sender,
			Cid:       cid,
			MsgId:     msg_id,
			Content:   content,
			TotalCnt:  int8(len(to)),
			Urgent:    int8(urgent),
			CreatedAt: int(now),
			UpdatedAt: int(now),
		}
		t.Chuoid = genChuoid(t)
	}
}

// 获取戳列表：我发出的和收到的
func GetChuoList() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

// 获取戳详情
func GetChuoInfo() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func genChuoid(t *models.TblChuoyixiaMeta) string {
	s, err := json.Marshal(t)
	if err != nil {
		return fmt.Sprintf("%s", md5.Sum([]byte(time.Now().String())))
	}
	return fmt.Sprintf("%x", md5.Sum(s))
}
