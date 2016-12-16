package options

import (
	"fmt"
	"log"
	"os"

	"tuohai/internal/console"
)

var Opts *Options

type Options struct {
	ID int
	HTTPAddress,
	LogFilePath string
	AuthHost string

	DbHost,
	DbUser,
	DbPwd,
	DbName string

	RedisHost,
	RedisPwd string

	//oss
	AccessKeyId,
	AccessKeySecret,
	OSSHost,
	AvatarBucket string

	FileOSSHost,
	FileHost,
	FileBucket string

	Logger *console.Console
}

func NewOptions() *Options {
	Opts = &Options{
		ID:          os.Getpid(),
		HTTPAddress: "0.0.0.0:10012",
		AuthHost:    "http://test.yunwoke.com",

		DbHost: "120.27.45.244:3306",
		DbUser: "root",
		DbPwd:  "yzjmysql",
		DbName: "fileadmin",

		RedisHost: "127.0.0.1:6379",
		RedisPwd:  "",

		AccessKeyId:     "LTAIACG7yVKjx0OY",
		AccessKeySecret: "RAYoo1YJM2Ah93pHp5qwLrCpnS24IA",
		OSSHost:         "oss-cn-beijing.aliyuncs.com",
		AvatarBucket:    "cloudwork",

		FileOSSHost: "oss-cn-qingdao.aliyuncs.com",
		FileBucket:  "clouderwork",
		FileHost:    "https://imfile.clouderwork.com",

		Logger: console.New(*log.New(os.Stderr, "file_api-", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)),
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
