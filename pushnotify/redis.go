package pushnotify

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

var rpool *redis.Pool

func InitRedis(RedisHost, RedisPwd string) {
	rpool = redis.NewPool(func() (redis.Conn, error) {
		return redis.Dial("tcp", RedisHost)
	}, 20)
}

func Subscribers(m chan redis.Message, key string) {
	c := rpool.Get()
	defer c.Close()

	psc := redis.PubSubConn{c}

	psc.Subscribe(key)
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			m <- v
		case redis.Subscription:
			fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
		case redis.PubSubConn:
			fmt.Println(v)
		case error:
			fmt.Println(v)
			return
		default:
		}
	}
}
