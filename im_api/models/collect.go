package models

import (
	"fmt"
	"time"
)

type MsgCollect struct {
	Id        int    `json:"-"`
	Collector string `json:"collector"`
	To        string `json:"cid"`
	Type      string `json:"type"`
	MsgId     uint64 `json:"mid"`
	Created   int64  `json:"time"`
	Updated   int64  `json:"-"`

	HistoryMsg
}

type HistoryMsg struct {
	From      string `gorm:"-" json:"from_id"`
	MsgData   string `gorm:"-" json:"msg_data"`
	Subtype   string `gorm:"-" json:"subtype"`
	GroupName string `gorm:"-" json:"group_name"`
}

func (c *MsgCollect) TableName() string {
	return "tbl_msg_collect"
}

func (c *MsgCollect) validationField() string {
	if c.Collector == "" {
		return "collector 不能为空"
	}
	if c.To == "" {
		return "to 不能为空"
	}
	if c.Type == "" {
		return "type 不能为空"
	}
	if c.MsgId == 0 {
		return "msgid 不能为空"
	}
	return ""
}

func (c *MsgCollect) getHistoryMsg() error {
	msg := &Message{}
	msg.To = c.To
	msg.MsgId = c.MsgId
	GetMessage(msg)

	c.From = msg.From
	c.MsgData = msg.MsgData
	c.Subtype = msg.Subtype
	return nil
}
func getCollect(collector, cid string, mid uint64) (*MsgCollect, error) {
	mc := &MsgCollect{}
	err := db.Find(mc, "collector = ? and `to` = ? and msg_id = ?", collector, cid, mid).Error
	fmt.Println(mc, err)
	return mc, err
}

func (c *MsgCollect) msgCollects(collector string, limit, offset int) ([]MsgCollect, int, error) {
	var (
		mc    []MsgCollect
		total int
	)

	err := db.Limit(limit).Offset(offset).Find(&mc, "`collector` = ?", collector).Error
	db.Table(c.TableName()).Where("`collector` = ?", collector).Count(&total)
	for i, _ := range mc {
		mc[i].getHistoryMsg()
	}

	return mc, total, err
}

func (c *MsgCollect) addMsgCollect() error {
	if err_msg := c.validationField(); err_msg != "" {
		return fmt.Errorf("%s", err_msg)
	}
	_, err := getCollect(c.Collector, c.To, c.MsgId)
	fmt.Println(err, RecordNotFound, err == RecordNotFound)
	if err != nil {
		if err.Error() == RecordNotFound.Error() {
			c.Created = time.Now().Unix()
			c.Updated = time.Now().Unix()
			return db.Create(c).Error
		}
		return err
	}
	return nil
}

func (c *MsgCollect) delMsgCollect() error {
	return db.Delete(c, "id > 0 and `to` = ? and msg_id = ?", c.To, c.MsgId).Error
}

func DelMsgCollect(collector, cid string, mid int) error {
	c, err := getCollect(collector, cid, uint64(mid))
	if err != nil {
		return err
	}
	return c.delMsgCollect()
}

//collector 收藏者
//cid 管道id
//mid消息唯一标识
//ctype 1单聊 2群聊
func AddMsgCollect(collector, cid, ctype string, mid uint64) error {
	mc := &MsgCollect{Collector: collector, To: cid, MsgId: mid, Type: ctype}
	return mc.addMsgCollect()
}

func CollectsByPaging(collector string, pageindex, pagesize int) ([]MsgCollect, int, error) {
	var mc MsgCollect
	if pageindex == 0 || pagesize == 0 {
		pageindex = 0
		pagesize = 20
	} else {
		pageindex = (pageindex - 1) * pagesize
	}
	fmt.Println("pageindex", pageindex, "pagesize", pagesize)
	return mc.msgCollects(collector, pagesize, pageindex)
}

func CollectsByOffset(collector string, limit, offset int) ([]MsgCollect, int, error) {
	var mc MsgCollect
	return mc.msgCollects(collector, limit, offset)
}
