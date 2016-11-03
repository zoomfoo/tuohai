package models

import (
	"fmt"

	"tuohai/internal/convert"
)

type TblMsg struct {
	Id        int    `gorm:"column:id"`
	From      string `gorm:"column:from"`
	To        string `gorm:"column:to"`
	Type      string `gorm:"column:type"`
	Subtype   string `gorm:"column:subtype"`
	MsgId     uint64 `gorm:"column:msg_id"`
	MsgData   string `gorm:"column:msg_data"`
	CreatedAt int    `gorm:"column:created_at"`
	UpdatedAt int    `gorm:"column:updated_at"`
}

func (t *TblMsg) TableName() string {
	return fmt.Sprintf("tbl_msg_%d", convert.RuneAccumulation(t.To)%16)
}

// AddTblMsg0 insert a new TblMsg0 into database and returns
// last inserted Id on success.
func AddTblMsg(m *TblMsg) (id int64, err error) {
	return 0, nil
}

// GetTblMsg0ById retrieves TblMsg0 by Id. Returns error if
// Id doesn't exist
func GetTblMsgById(to string) ([]TblMsg, error) {
	var (
		msgs []TblMsg
	)

	err := db.Table((&TblMsg{To: to}).TableName()).Where("`to` = ?", to).Scan(&msgs).Error
	return msgs, err
}

// GetAllTblMsg0 retrieves all TblMsg0 matches certain condition. Returns empty list if
// no records exist
func GetAllTblMsg(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	return nil, nil
}

// UpdateTblMsg0 updates TblMsg0 by Id and returns error if
// the record to be updated doesn't exist
func UpdateTblMsg0ById(m *TblMsg) (err error) {
	return nil
}

// DeleteTblMsg0 deletes TblMsg0 by Id and returns error if
// the record to be deleted doesn't exist
func DeleteTblMsg0(id int) (err error) {
	return nil
}
