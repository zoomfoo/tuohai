package models

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

func MsgReadInfo(cid, msgid, origin string) (int, map[string][]string, error) {
	// TODO 需要校验origin是否是当前用户，如果不是则返回空
	res := make(map[string][]string)
	c := rpool.Get()
	defer c.Close()

	// 获取消息未读数
	key := "msg:unread:cnt:" + cid + ":" + msgid + ":" + origin
	cnt, err := c.Do("GET", key)
	if err != nil {
		log.Println(err)
		return 0, nil, err
	}
	cc, _ := redis.Int(cnt, nil)
	// 获取消息未读人员列表
	key = "msg:unread:list:" + cid + ":" + msgid + ":" + origin
	unlist, err := c.Do("SMEMBERS", key)
	if err != nil {
		log.Println(err)
		return 0, nil, err
	}
	res["unread"], _ = redis.Strings(unlist, nil)
	// 获取消息已读人员列表
	key = "msg:read:list:" + cid + ":" + msgid + ":" + origin
	rlist, err := c.Do("SMEMBERS", key)
	if err != nil {
		log.Println(err)
		return 0, nil, err
	}
	res["read"], _ = redis.Strings(rlist, nil)
	return cc, res, nil
}

func SimpleUnread(userid, sid int) int {
	return SimpleUnreads(userid)[strconv.Itoa(sid)]
}

func SimpleUnreads(userid int) map[string]int {
	c := rpool.Get()
	defer c.Close()

	if _, err := c.Do("select", "1"); err != nil {
		log.Println(err)
		return nil
	}

	if val, err := redis.IntMap(c.Do("hgetall", "unread_"+strconv.Itoa(userid))); err != nil {
		log.Println(err)
		return map[string]int{}
	} else {
		return val
	}
}

func GroupUnread(uid, gid int) int {
	c := rpool.Get()
	defer c.Close()

	if _, err := c.Do("select", "1"); err != nil {
		log.Println(err)
		return 0
	}
	//获取群所有消息
	gkey := fmt.Sprintf("%d_im_group_msg", gid)
	v, err := redis.IntMap(c.Do("hgetall", gkey))
	if err != nil {
		log.Println(err)
		return 0
	}

	if count, ok := v["count"]; ok {
		//获取已读的消息
		ugkey := fmt.Sprintf("%d_%d_im_user_group", uid, gid)
		ug, err := redis.IntMap(c.Do("hgetall", ugkey))
		if err != nil {
			log.Println(err)
			return 0
		}

		if _, ok := ug["count"]; !ok {
			return 0
		}

		return count - ug["count"]
	} else {
		return 0
	}
}

//删除uid 的pid消息
func CleanSimpleAlRead(uid, pid int) (int, error) {
	c := rpool.Get()
	defer c.Close()

	if _, err := c.Do("select", "1"); err != nil {
		log.Println(err)
		return 0, err
	}

	if res, err := redis.Int(c.Do("hdel", "unread_"+strconv.Itoa(uid), pid)); err != nil {
		return 0, err
	} else {
		return res, nil
	}
}

func CleanGroupAlRead(uid, gid int) (bool, error) {
	c := rpool.Get()
	defer c.Close()

	if _, err := c.Do("select", "1"); err != nil {
		log.Println(err)
		return false, err
	}

	//获取群所有消息
	gkey := fmt.Sprintf("%d_im_group_msg", gid)
	v, err := redis.IntMap(c.Do("hgetall", gkey))
	if err != nil {
		log.Println(err)
		return false, err
	}

	if count, ok := v["count"]; ok {
		//获取已读的消息
		ugkey := fmt.Sprintf("%d_%d_im_user_group", uid, gid)
		ug, err := redis.IntMap(c.Do("hgetall", ugkey))
		if err != nil {
			log.Println(err)
			return false, err
		}

		if _, ok := ug["count"]; !ok {
			log.Println("ug[count] == nil")
			return false, nil
		}

		if ok, err := redis.String(c.Do("hmset", ugkey, "count", count)); err != nil {
			return false, err
		} else if ok == "OK" {
			return true, nil
		} else {
			return false, err
		}

	} else {
		return false, nil
	}
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
	if _, err := c.Do("select", "5"); err != nil {
		return []string{}, err
	}

	return redis.Strings(c.Do("hgetall", fmt.Sprintf("group:member:%s", gid)))
}

func QuitGroup(gid string, member []string) (bool, error) {
	c := rpool.Get()
	defer c.Close()

	if _, err := c.Do("select", "5"); err != nil {
		return false, err
	}

	var args = []interface{}{fmt.Sprintf("group:member:%s", gid)}
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

	if _, err := c.Do("select", "5"); err != nil {
		return false, err
	}

	res, err := redis.Int(c.Do("hexists", fmt.Sprintf("group:member:%s", gid), uid))
	if err != nil {
		return false, err
	}

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

	if _, err := c.Do("hmset", val...); err != nil {
		return err
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
	if _, err := c.Do("set", key, str); err != nil {
		return err
	}

	if err := Publish("cw:bot:add", str); err != nil {
		return err
	}
	return nil
}
