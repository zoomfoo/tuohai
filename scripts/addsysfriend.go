package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"

	"tuohai/im_api/models"
	"tuohai/im_api/options"
)

var sysusers = []string{"84558b0cf90a4166", "e4b1b4018c147b1c"}

func main() {
	opts := options.NewOptions()
	models.InitDB(opts.MysqlOptions())
	models.InitRedis(opts.RedisHost, opts.RedisPwd)
	fmt.Printf("system users:%s\n", sysusers)
	allusers, err := models.GetAllUsers()
	if err != nil {
		panic("error")
	}

	for _, user := range allusers {
		u := user.Uuid
		fmt.Printf("处理 user:%s\n", u)
		for _, sys := range sysusers {
			if u == sys {
				continue
			}
			_, err := models.GetSysRid(sys, u)
			if err != nil {
				fmt.Printf("error:user[%s],sys[%s]", u, sys)
			}
			fmt.Printf("系统好友添加成功\n")
		}
	}

}
