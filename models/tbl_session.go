package models

import (
	"fmt"

	"tuohai/internal/convert"
)

type TblSession struct {
	Id        int    `gorm:"column:id"`
	Sid       string `gorm:"column:sid"`
	From      string `gorm:"column:from"`
	To        string `gorm:"column:to"`
	Type      uint8  `gorm:"column:type"`
	Status    int8   `gorm:"column:status"`
	CreatedAt int    `gorm:"column:created_at"`
	UpdatedAt int    `gorm:"column:updated_at"`
}

func (t *TblSession) TableName() string {
	return fmt.Sprintf("tbl_session_%d", convert.RuneAccumulation(t.From)%4)
}

// AddTblSession0 insert a new TblSession0 into database and returns
// last inserted Id on success.
func AddTblSession(m *TblSession) (id int64, err error) {
	return 0, nil
}

// GetTblSession0ById retrieves TblSession0 by Id. Returns error if
// Id doesn't exist
func GetTblSessionById(from string) ([]TblSession, error) {
	var sessions []TblSession
	sess := &TblSession{From: from}
	fmt.Println(sess.TableName())
	err := db.Table(sess.TableName()).Where("`from` = ? and status = 0", from).Scan(&sessions).Error
	return sessions, err
}

// GetAllTblSession0 retrieves all TblSession0 matches certain condition. Returns empty list if
// no records exist
func GetAllTblSession(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	return nil, nil
}

// UpdateTblSession0 updates TblSession0 by Id and returns error if
// the record to be updated doesn't exist
func UpdateTblSessionById(m *TblSession) (err error) {
	return nil
}

// DeleteTblSession0 deletes TblSession0 by Id and returns error if
// the record to be deleted doesn't exist
func DeleteTblSession(id int) (err error) {
	return nil
}
