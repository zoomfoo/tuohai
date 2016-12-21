package models

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
	"tuohai/internal/console"
)

func MsgReadInfo(cid, msgid, origin string) (int, map[string][]string, error) {
	// TODO 需要校验origin是否是当前用户，如果不是则返回空
	res := make(map[string][]string)
	c := rpool.Get()
	defer c.Close()

	// 获取消息未读数
	key := "msg:unread:cnt:" + cid + ":" + msgid + ":" + origin

	cc, err := redis.Int(c.Do("GET", key))
	if err != nil {
		log.Println(err)
		return 0, nil, err
	}

	// 获取消息未读人员列表
	key = "msg:unread:list:" + cid + ":" + msgid + ":" + origin
	var unread_err error
	if res["unread"], unread_err = redis.Strings(c.Do("SMEMBERS", key)); unread_err != nil {
		log.Println(unread_err)
		return 0, nil, unread_err
	}
	fmt.Println("unread key: ", key)
	// 获取消息已读人员列表
	key = "msg:read:list:" + cid + ":" + msgid + ":" + origin
	var read_err error
	if res["read"], read_err = redis.Strings(c.Do("SMEMBERS", key)); read_err != nil {
		log.Println(read_err)
		return 0, nil, read_err
	}
	return cc, res, nil
}

func MsgUnreadCount(cid, msgid, origin string) int {
	c := rpool.Get()
	defer c.Close()

	// 获取消息未读数
	key := "msg:unread:cnt:" + cid + ":" + msgid + ":" + origin
	cc, err := redis.Int(c.Do("GET", key))
	if err != nil {
		log.Println(err)
		return 0
	}
	return cc
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

func GetGroupMem(gid string) ([]string, error) {
	c := rpool.Get()
	defer c.Close()
	return redis.Strings(c.Do("hgetall", fmt.Sprintf("channel:member:%s", gid)))
}

func QuitGroup(gid string, member []string) (bool, error) {
	c := rpool.Get()
	defer c.Close()

	var args = []interface{}{fmt.Sprintf("channel:member:%s", gid)}
	for _, m := range member {
		args = append(args, m)
	}

	if _, err := c.Do("hdel", args...); err != nil {
		return false, err
	}
	return true, nil
}

func IsGroupMember(gid, uid string) (bool, error) {
	c := rpool.Get()
	defer c.Close()

	res, err := redis.Int64(c.Do("hexists", fmt.Sprintf("channel:member:%s", gid), uid))
	if err != nil {
		return false, err
	}
	fmt.Println("当前[", uid, "]是否在群中", res)
	return res == 1, nil
}

func saveChennelToRedis(cid string, members []string) error {
	var (
		val = []interface{}{"channel:member:" + cid}
		now = time.Now().Unix()
	)

	c := rpool.Get()
	defer c.Close()

	for _, mem := range members {
		val = append(val, mem, now)
	}

	fmt.Println("val: ", val)
	if i, err := c.Do("hmset", val...); err != nil {
		return err
	} else {
		fmt.Println("写入redis:返回 ", i)
	}
	return nil
}

func Publish(args ...interface{}) error {
	c := rpool.Get()
	defer c.Close()
	if _, err := c.Do("publish", args...); err != nil {
		return err
	}
	return nil
}

func SaveBotInfo(key string, info map[string]interface{}) error {
	js, err := json.Marshal(info)
	if err != nil {
		return err
	}

	c := rpool.Get()
	defer c.Close()

	str := string(js)
	fmt.Println("发布并缓存bot信息: ", str)
	//缓存一份数据保证离线消息最终一致
	if _, err := c.Do("set", key, str); err != nil {
		return err
	}
	//发布一份数据
	if err := Publish("cw:bot:add", str); err != nil {
		return err
	}
	return nil
}

func ChennelUnreadNum(cid, uid string) int {
	c := rpool.Get()
	defer c.Close()
	res, err := redis.String(c.Do("hmget", "cnt:unread"+cid, uid))
	if err != nil {
		console.StdLog.Error(err)
		return 0
	}
	i, _ := strconv.Atoi(res)
	return i
}

func CleanSessionUnread(cid, uid string) bool {
	c := rpool.Get()
	defer c.Close()
	_, err := c.Do("hdel", "cnt:unread"+cid, uid)
	if err != nil {
		console.StdLog.Error(err)
		return false
	}
	return true
}
