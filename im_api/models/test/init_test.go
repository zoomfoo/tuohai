package test

import (
	"tuohai/im_api/models"
	api "tuohai/im_api/options"
)

func init() {
	Opts := api.NewOptions()
	models.InitDB(Opts.MysqlOptions())
	models.InitRedis(Opts.RedisHost, api.Opts.RedisPwd)

	//初始化主站数据库
	models.InitMainSiteDB(Opts.MainSiteMysql())
}
