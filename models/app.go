package models

import "time"

type AppState int

const (
	AppDisable AppState = iota
	AppNormal
	AppEscape
)

type App struct {
	Id             string    `gorm:"column:id"`
	Idx            int       `gorm:"column:idx"`
	Name           string    `gorm:"column:name"`
	Icon           string    `gorm:"column:icon"`
	Desc           string    `gorm:"column:desc"`
	ConfigHelpInfo string    `gorm:"column:config_help_info"`
	Url            string    `gorm:"column:url"`
	HelpUrl        string    `gorm:"column:help_url"`
	State          int       `gorm:"column:state"`
	CreateTime     time.Time `gorm:"column:create_time"`
	UpTime         time.Time `gorm:"column:up_time"`
	AppURL         string    `gorm:"column:app_url"`
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
