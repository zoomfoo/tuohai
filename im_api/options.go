package im_api

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
	LogFilePath,
	AuthHost string
	RPCHost string

	DbHost,
	DbUser,
	DbPwd,
	DbName string

	RedisHost,
	RedisPwd string

	Logger *console.Console
}

func NewOptions() *Options {
	Opts = &Options{
		ID:          os.Getpid(),
		HTTPAddress: "0.0.0.0:10011",
		AuthHost:    "http://test.yunwoke.com",
		RPCHost:     "127.0.0.1:5004",

		DbHost: "120.27.45.244:3306",
		DbUser: "root",
		DbPwd:  "yzjmysql",
		DbName: "newim",

		RedisHost: "120.27.45.244:6379",
		RedisPwd:  "",

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
