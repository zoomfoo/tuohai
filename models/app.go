package models

import "time"

type AppState int

const (
	AppDisable AppState = iota
	AppNormal
	AppEscape
)

type App struct {
	Id             string    `gorm:"column:id" json:"id"`
	Idx            int       `gorm:"column:idx" json:"-"`
	Name           string    `gorm:"column:name" json:"name"`
	Icon           string    `gorm:"column:icon" json:"icon"`
	Desc           string    `gorm:"column:desc" json:"desc"`
	ConfigHelpInfo string    `gorm:"column:config_help_info" json:"config_help_info"`
	Url            string    `gorm:"column:url" json:"url"`
	HelpUrl        string    `gorm:"column:help_url" json:"help_url"`
	State          int       `gorm:"column:state" json:"state"`
	CreateTime     time.Time `gorm:"column:create_time" json:"create_time"`
	UpTime         time.Time `gorm:"column:up_time" json:"update_time"`
	AppURL         string    `gorm:"column:app_url" json:"app_url"`
}

func (app *App) TableName() string {
	return "app"
}

func Apps() ([]App, error) {
	var apps []App
	err := db.Find(&apps, "state = ?", AppNormal).Error
	return apps, err
}

func GetAppById(appid string) (*App, error) {
	a := &App{}
	err := db.Find(a, "id = ?", appid).Error
	return a, err
}
