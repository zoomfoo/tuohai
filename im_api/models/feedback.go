package models

import (
	"time"
)

type Feedback struct {
	Id         int    `gorm:"column:id"`
	UserId     string `gorm:"column:user_id"`
	Content    string `gorm:"column:content"`
	CreateTime int64  `gorm:"column:create_time"`
}

func (fb *Feedback) TableName() string {
	return "tbl_feedback"
}

func NewFeedback(uid, content string) error {
	fb := &Feedback{UserId: uid, Content: content, CreateTime: time.Now().Unix()}
	return db.Create(fb).Error
}
