package options

import (
	"fmt"
	"log"
	"os"

	"tuohai/internal/console"
)

var Opts *Options

type Options struct {
	ID          int
	HTTPAddress string
	LogFilePath string
	AuthHost    string
	RPCHost     string
	WebHookHost string

	DbHost string
	DbUser string
	DbPwd  string
	DbName string

	RedisHost string
	RedisPwd  string

	SysUserYunliao     string
	SysUserClouderwork string

	Logger *console.Console
}

func NewOptions() *Options {
	Opts = &Options{
		ID:          os.Getpid(),
		HTTPAddress: "0.0.0.0:10011",
		AuthHost:    "http://test.yunwoke.com",
		RPCHost:     "127.0.0.1:9003",
		WebHookHost: "test.yunwoke.com:8880",

		DbHost: "120.27.45.244:3306",
		DbUser: "root",
		DbPwd:  "yzjmysql",
		DbName: "newim",

		RedisHost: "120.27.45.244:6379",
		RedisPwd:  "",

		SysUserYunliao:     "84558b0cf90a4166",
		SysUserClouderwork: "e4b1b4018c147b1c",

		Logger: console.New(*log.New(os.Stderr, "im_api-", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)),
	}
	return Opts
}

func (opts *Options) MysqlOptions() string {
	if opts.DbUser == "" || opts.DbPwd == "" || opts.DbHost == "" || opts.DbName == "" {
		log.Println("Database related field cannot be empty")
		os.Exit(1)
	}
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		opts.DbUser, opts.DbPwd, opts.DbHost, opts.DbName)
}

func (opts *Options) MainSiteMysql() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		"root", "yzjmysql", "120.27.45.244:3306", "cloudwork")
}
