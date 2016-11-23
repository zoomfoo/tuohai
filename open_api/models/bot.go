package models

import "time"

type Bot struct {
	Id         string    `gorm:"column:id" json:"id" form:"-"`
	Idx        int       `gorm:"column:idx" json:"-" form:"-"`
	Name       string    `gorm:"column:name" json:"name" form:"name" binding:"required"`
	Icon       string    `gorm:"column:icon" json:"icon" form:"icon" binding:"required"`
	CreatorId  string    `gorm:"column:creator_id" json:"creator_id" form:"creator_id" binding:"required"`
	ChannelId  string    `gorm:"column:channel_id" json:"channel_id" form:"channel_id" binding:"required"`
	AppId      string    `gorm:"column:app_id" json:"app_id" form:"app_id" binding:"required"`
	State      int       `gorm:"column:state" json:"state" form:"state"`
	CreateTime time.Time `gorm:"column:create_time" json:"-" form:"-"`
	UpTime     time.Time `gorm:"column:up_time" json:"-" form:"-"`
	IsPub      int       `gorm:"column:is_pub" json:"is_pub" form:"-"`
}

func (b *Bot) TableName() string {
	return "bot"
}

func GetBotById(bot_id string) (*Bot, error) {
	var b Bot
	err := db.Find(&b, "id = ? and state = 1", bot_id).Error
	return &b, err
}

func GetBots(id []string) ([]Bot, error) {
	var b []Bot
	err := db.Find(&b, "state = 1 and channel_id in (?)", id).Error
	return b, err
}

func CreateBot(b *Bot) error {
	return db.Create(b).Error
}
