package models

import (
	"fmt"

	"tuohai/internal/convert"
)

const (
	normal = iota
	deleted
)

type SessionType int8

const (
	SimpleSession SessionType = 1
	GroupSession  SessionType = 2
)

type Session struct {
	Id        int         `gorm:"column:id" json:"-"`
	Sid       string      `gorm:"column:sid" json:"sid"`
	From      string      `gorm:"column:from" json:"from"`
	To        string      `gorm:"column:to" json:"to"`
	Status    int8        `gorm:"column:status" json:"-"`
	SType     SessionType `gorm:"column:type" json:"type"`
	CreatedAt int         `gorm:"column:created_at" json:"-"`
	UpdatedAt int         `gorm:"column:updated_at" json:"-"`
}

func (t *Session) TableName() string {
	return fmt.Sprintf("tbl_session_%d", convert.RuneAccumulation(t.From)%4)
}

func GetSessionById(from string) ([]Session, error) {
	var sessions []Session
	sess := &Session{From: from}
	err := db.Table(sess.TableName()).Where("`from` = ? and status = 0", from).Scan(&sessions).Error
	return sessions, err
}
