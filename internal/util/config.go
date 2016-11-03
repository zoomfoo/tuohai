package util

import (
	"github.com/noaway/config"
	"log"
	"path/filepath"
)

var ConfPath = "../../conf/app.conf"

func Config(group, key string) string {
	path, err := filepath.Abs(ConfPath)
	if err != nil {
		log.Fatalln(err)
		return ""
	}

	conf, err := config.ReadDefault(path)
	if err != nil {
		log.Fatalln(err)
		return ""
	}
	value, err := conf.String(group, key)
	if err != nil {
		log.Fatalln(err.Error())
		return ""
	}
	return value
}

func ConfigInt(group, key string) int {
	conf, err := config.ReadDefault(ConfPath)
	if err != nil {
		log.Fatalln(err.Error())
		return 0
	}
	value, err := conf.Int(group, key)
	if err != nil {
		log.Fatalln(err.Error())
		return 0
	}
	return value
}
