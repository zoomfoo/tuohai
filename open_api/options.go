package open_api

import (
	"fmt"
	"log"
	"os"

	"tuohai/internal/console"
	"tuohai/internal/util"
)

var Opts *Options

type Options struct {
	ID int
	HTTPAddress,
	LogFilePath string

	DbHost,
	DbUser,
	DbPwd,
	DbName string

	ConnLogicRPCAddress string

	IMAPI_HOST  string
	WebHookHOST string

	Logger *console.Console
}

func NewOptions() *Options {
	Opts = &Options{
		ID:          os.Getpid(),
		HTTPAddress: util.Config("open_api", "HTTPAddress"),
		LogFilePath: "",

		DbHost: util.Config("open_api", "DbHost"),
		DbUser: util.Config("open_api", "DbUser"),
		DbPwd:  util.Config("open_api", "DbPwd"),
		DbName: util.Config("open_api", "DbName"),

		ConnLogicRPCAddress: util.Config("open_api", "ConnLogicRPCAddress"),

		IMAPI_HOST:  util.Config("open_api", "IMAPI_HOST"),
		WebHookHOST: util.Config("open_api", "WebHookHOST"),

		Logger: console.New(*log.New(os.Stderr, "open_api-", log.Ldate|log.Ltime|log.Lmicroseconds)),
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

func GenerateConfFile() {

}

var conf = `[im_api]



[open_api]
ConnLogicRPCAddress: 127.0.0.1:5004
HTTPAddress: :8089
DbHost: 120.27.45.244:3306
DbUser: root
DbPwd:  yzjmysql
DbName: im_bot_app

IMAPI_HOST: http://127.0.0.1:10011`
