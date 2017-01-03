package models

import (
	"crypto/md5"
	"errors"
	"fmt"
	"time"

	"tuohai/internal/convert"
)

const (
	normal = iota
	deleted
)

type SessionType int8

const (
	SimpleSession    SessionType = 1
	GroupSession     SessionType = 2
	SingleTmpSession SessionType = 3
)

type Session struct {
	Id        int         `gorm:"column:id" json:"-"`
	Sid       string      `gorm:"column:sid" json:"sid"`
	From      string      `gorm:"column:from" json:"from"`
	To        string      `gorm:"column:to" json:"to"`
	Status    int8        `gorm:"column:status" json:"-"`
	SType     SessionType `gorm:"column:type" json:"type"`
	CreatedAt int64       `gorm:"column:created_at" json:"-"`
	UpdatedAt int64       `gorm:"column:updated_at" json:"-"`
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

func RemoveSession(sid, uid string) error {
	s := &Session{From: uid}
	return db.Table(s.TableName()).Where("sid = ?", sid).Updates(map[string]interface{}{"status": deleted, "updated_at": time.Now().Unix()}).Error
}

func GetSessionBySid(sid, uid string) (*Session, error) {
	s := &Session{From: uid}
	err := db.Table(s.TableName()).Where("sid = ?", sid).Scan(s).Error
	return s, err
}

func RemoveSessionByCidAndUid(cid, uid string) error {
	s := &Session{From: uid}
	return db.Table(s.TableName()).Where("`from` = ? and `to` = ?", uid, cid).Updates(map[string]interface{}{"status": deleted, "updated_at": time.Now().Unix()}).Error
}

func CreateTmpSession(from, to string) (*Session, error) {
	sid := fmt.Sprintf("%x", md5.Sum([]byte(from+fmt.Sprintf("%d", time.Now().UnixNano()))))
	s := &Session{
		Sid:       sid,
		From:      from,
		To:        to,
		Status:    0,
		SType:     SingleTmpSession,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	// 查询是否已有临时会话
	tms := &Session{From: from}
	err := db.Table(s.TableName()).Where("`from` = ? and `to` = ? and type = 3 and status = 0", from, to).Scan(tms).Error
	if err == nil {
		return tms, errors.New("temp session is already existed")
	}
	// 创建新记录
	return s, db.Table(s.TableName()).Create(s).Error

}
