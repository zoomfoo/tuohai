package main

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"tuohai/im_api/models"
)

type DB struct {
	DbUser,
	DbPwd,
	DbHost,
	DbName string

	RedisHost,
	RedisPwd string
}

func NewDb() *DB {
	return &DB{
		DbUser: "",
		DbPwd:  "",
		DbHost: "",
		DbName: "",

		RedisHost: "",
		RedisPwd:  "",
	}
}

func main() {
	db := NewDb()
	models.InitDB(fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		db.DbUser, db.DbPwd, db.DbHost, db.DbName))
	models.InitRedis(db.RedisHost, db.RedisPwd)
}
