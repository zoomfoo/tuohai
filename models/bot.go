package models

import "time"

type Bot struct {
	Id         string    `gorm:"column:id" json:"id"`
	Idx        int       `gorm:"column:idx" json:"-"`
	Name       string    `gorm:"column:name" json:"name"`
	Icon       string    `gorm:"column:icon" json:"icon"`
	CreatorId  string    `gorm:"column:creator_id" json:"creator_id"`
	ChannelId  string    `gorm:"column:channel_id" json:"channel_id"`
	AppId      string    `gorm:"column:app_id" json:"app_id"`
	State      int       `gorm:"column:state" json:"state"`
	CreateTime time.Time `gorm:"column:create_time" json:"-"`
	UpTime     time.Time `gorm:"column:up_time" json:"-"`
	IsPub      int       `gorm:"column:is_pub" json:"is_pub"`
}

func (b *Bot) TableName() string {
	return "bot"
}

func GetBotById(bot_id string) (*Bot, error) {
	var b Bot
	err := db.Find(&b, "id = ? and state = 1", bot_id).Error
	return &b, err
}
