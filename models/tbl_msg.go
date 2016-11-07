package models

import (
	"fmt"

	"tuohai/internal/convert"
)

type TblMsg struct {
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

func (t *TblMsg) TableName() string {
	return fmt.Sprintf("tbl_msg_%d", convert.RuneAccumulation(t.To)%16)
}

func GetTblMsgById(to string) ([]TblMsg, error) {
	var (
		msgs []TblMsg
	)

	err := db.Table((&TblMsg{To: to}).TableName()).Where("`to` = ?", to).Scan(&msgs).Error
	return msgs, err
}

func GetLastHistory(to string) (*TblMsg, error) {
	msg := &TblMsg{To: to}
	err := db.Order("`msg_id` desc").First(&msg, "`to` = ?", to).Error
	return msg, err
}
