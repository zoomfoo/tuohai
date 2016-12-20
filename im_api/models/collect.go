package models

// import (
// 	"fmt"
// )

// type MsgCollect struct {
// 	Id        int    `gorm:"column:id"`
// 	Collector string `gorm:"column:collector"`
// 	To        string `gorm:"column:to"`
// 	Type      string `gorm:"column:"type"`
// 	MsgId     int    `gorm:"column:msg_id"`
// 	Created   int64  `gorm:"column:created"`
// 	Updated   int64  `gorm:"column:updated"`

// 	HistoryMsg
// }

// type HistoryMsg struct {
// 	From      string `gorm:"-" json:"from_id"`
// 	MsgData   string `gorm:"-" json:"msg_data"`
// 	Subtype   string `gorm:"-" json:"subtype"`
// 	GroupName string `gorm:"-" json:"group_name"`
// }

// func (c *MsgCollect) TableName() string {
// 	return "tbl_msg_collect"
// }

// func (c *MsgCollect) validationField() string {
// 	if c.Collector == "" {
// 		return "collector 不能为空"
// 	}
// 	if c.To == "" {
// 		return "to 不能为空"
// 	}
// 	if c.Type == "" {
// 		return "type 不能为空"
// 	}
// 	if c.MsgId == 0 {
// 		return "msgid 不能为空"
// 	}
// 	return ""
// }

// func (c *MsgCollect) getHistoryMsg() error {
// 	msg := &Message{To: c.To, MsgId, c.MsgId}
// 	res := GetMessage(msg)

// 	c.From = msg.From
// 	c.MsgData = msg.MsgData
// 	c.Subtype = msg.Subtype
// }
// func (c *MsgCollect) getCollect() {
// 	// db.Find(c, ...)
// }

// func (c *MsgCollect) msgCollects(collector, limit, offset string) ([]MsgCollect, error) {
// 	var mc []MsgCollect
// 	db.Limit(limit).Offset(offset).Order("created desc").Find(&mc, "`collector` = ?", collector)

// 	for i, _ := range mc {
// 		mc[i].GetHistoryMsg()
// 	}

// 	return mc, nil
// }

// func (c *MsgCollect) addMsgCollect() error {
// 	if err_msg := c.ValidationField(); err_msg != "" {
// 		return fmt.Errorf("%s", err_msg)
// 	}

// }

// func (c *MsgCollect) DelMsgCollect() {

// }

// //collector 收藏者
// //cid 管道id
// //mid消息唯一标识
// //ctype 1单聊 2群聊
// func AddMsgCollect(collector, cid, mid, ctype string) error {
// 	mc := &MsgCollect{Collector: collector, To: cid, MsgId: mid, Type: ctype}
// 	return mc.AddMsgCollect()
// }
