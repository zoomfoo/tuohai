package models

import (
	"fmt"

	"tuohai/internal/convert"
)

type Message struct {
	Id        int    `gorm:"column:id" json:"-"`
	From      string `gorm:"column:from" json:"from"`
	To        string `gorm:"column:to" json:"cid"`
	Type      string `gorm:"column:type" json:"type"`
	Subtype   string `gorm:"column:subtype" json:"sub_type"`
	MsgId     uint64 `gorm:"column:msg_id" json:"msg_id"`
	MsgData   string `gorm:"column:msg_data" json:"msg_data"`
	CreatedAt int    `gorm:"column:created_at" json:"create_time"`
	UpdatedAt int    `gorm:"column:updated_at" json:"-"`
}

func (t *Message) TableName() string {
	return fmt.Sprintf("tbl_msg_%d", convert.RuneAccumulation(t.To)%16)
}

func GetMsgById(cid, mid, size string) ([]Message, error) {
	var (
		msgs []Message
	)

	fmt.Println("消息数据库表名称: ", fmt.Sprintf("tbl_msg_%d", convert.RuneAccumulation(cid)%16))
	if mid != "" && mid != "0" {
		err := db.Table((&Message{To: cid}).TableName()).
			Where("`to` = ? and `msg_id` <= ?", cid, mid).Limit(size).Order("msg_id desc").
			Scan(&msgs).Error
		return msgs, err
	} else {
		err := db.Table((&Message{To: cid}).TableName()).
			Where("`to` = ?", cid).Limit(size).Order("msg_id desc").
			Scan(&msgs).Error
		return msgs, err
	}
}

func GetLastHistory(to string) (*Message, error) {
	msg := &Message{To: to}
	err := db.Order("`msg_id` desc").First(&msg, "`to` = ?", to).Error
	return msg, err
}
