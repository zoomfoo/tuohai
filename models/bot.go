package models

import "time"

type Bot struct {
	Id         string    `gorm:"column:id"`
	Idx        int       `gorm:"column:idx"`
	Name       string    `gorm:"column:name"`
	Icon       string    `gorm:"column:icon"`
	CreatorId  string    `gorm:"column:creator_id"`
	ChannelId  string    `gorm:"column:channel_id"`
	AppId      string    `gorm:"column:app_id"`
	State      int       `gorm:"column:state"`
	CreateTime time.Time `gorm:"column:create_time"`
	UpTime     time.Time `gorm:"column:up_time"`
	IsPub      int       `gorm:"column:is_pub"`
}

func (b *Bot) TableName() string {
	return "bot"
}

func GetBotById(bot_id string) (*Bot, error) {
	var b Bot
	err := db.Find(&b, "id = ? and state = 1", bot_id).Error
	return &b, err
}
