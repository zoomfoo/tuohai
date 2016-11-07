package models

import (
	"fmt"

	"tuohai/internal/convert"
)

type TblSession struct {
	Id        int    `gorm:"column:id" json:"-"`
	Sid       string `gorm:"column:sid" json:"sid"`
	From      string `gorm:"column:from" json:"from"`
	To        string `gorm:"column:to" json:"to"`
	Status    int8   `gorm:"column:status" json:"-"`
	CreatedAt int    `gorm:"column:created_at" json:"-"`
	UpdatedAt int    `gorm:"column:updated_at" json:"-"`
}

func (t *TblSession) TableName() string {
	return fmt.Sprintf("tbl_session_%d", convert.RuneAccumulation(t.From)%4)
}

func GetTblSessionById(from string) ([]TblSession, error) {
	var sessions []TblSession
	sess := &TblSession{From: from}
	err := db.Table(sess.TableName()).Where("`from` = ? and status = 0", from).Scan(&sessions).Error
	return sessions, err
}
