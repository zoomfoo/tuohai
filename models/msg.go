package models

import (
	"fmt"

	"tuohai/internal/convert"
)

type Message struct {
	Id        int    `gorm:"column:id" json:"-"`
	From      string `gorm:"column:from" json:"fr"`
	To        string `gorm:"column:to" json:"to"`
	Type      string `gorm:"column:type" json:"type"`
	Subtype   string `gorm:"column:subtype" json:"sub_type"`
	MsgId     uint64 `gorm:"column:msg_id" json:"mid"`
	MsgData   string `gorm:"column:msg_data" json:"data"`
	CreatedAt int    `gorm:"column:created_at" json:"time"`
	UpdatedAt int    `gorm:"column:updated_at" json:"-"`
}

type Msgrecord struct {
	FromId      uint32 `json:"user_id"`
	SessionType int8   `json:"session_type"`
	ToId        uint32 `json:"session_id"`
	MsgIdBegin  uint32 `json:"msg_id_begin"`
	MsgCnt      uint32 `json:"msg_cnt"`
}

func (t *Msg) TableName() string {
	return fmt.Sprintf("tbl_msg_%d", convert.RuneAccumulation(t.To)%16)
}

func GetMsgById(record *Msgrecord) ([]Message, error) {
	var (
		msgs []Message
	)

	err := db.Table((&Message{To: to}).TableName()).
		Where("`to` = ? and msg_id <=? order by created_at desc, id desc limit ?", record.ToId, record.MsgIdBegin, record.MsgCnt).
		Scan(&msgs).Error
	return msgs, err
}

func GetLastHistory(to string) (*Message, error) {
	msg := &Message{To: to}
	err := db.Order("`msg_id` desc").First(&msg, "`to` = ?", to).Error
	return msg, err
}
