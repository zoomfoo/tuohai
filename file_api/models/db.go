package models

import (
	"log"

	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var (
	db    *gorm.DB
	rpool *redis.Pool
)

type DbModel struct {
	gorm *gorm.DB
}

func InitDB(MysqlOptions string) (m *DbModel, err error) {
	m = &DbModel{}
	m.gorm, err = gorm.Open("mysql", MysqlOptions)
	if err != nil {
		return nil, err
	}
	db = m.gorm
	db.DB().SetMaxOpenConns(100)
	db.DB().SetMaxIdleConns(10)

	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	db.LogMode(true)
	return
}

func InitRedis(RedisHost, RedisPwd string) {
	rpool = redis.NewPool(func() (redis.Conn, error) {
		return redis.Dial("tcp", RedisHost)
	}, 20)
}

func (db *DbModel) Close() {
	if db.gorm != nil {
		db.gorm.Close()
	}
}
