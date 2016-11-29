package test

import (
	api "tuohai/im_api"
	"tuohai/models"
)

func init() {
	Opts := api.NewOptions()
	models.InitDB(Opts.MysqlOptions())
	models.InitRedis(Opts.RedisHost, api.Opts.RedisPwd)
}
