package models

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"strconv"
)

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

func QuitGroup(gid uint32, uids []uint32) (bool, error) {
	c := rpool.Get()
	defer c.Close()

	if _, err := c.Do("select", "5"); err != nil {
		return false, err
	}

	var args = []interface{}{fmt.Sprintf("group_member_%d", gid)}
	for _, uid := range uids {
		args = append(args, uid)
	}

	if b, err := c.Do("hdel", args...); err != nil {
		return false, err
	} else {
		log.Println(b)
	}

	return true, nil
}

func IsGroupMember(gid, uid int) (bool, error) {
	c := rpool.Get()
	defer c.Close()

	if _, err := c.Do("select", "5"); err != nil {
		return false, err
	}

	res, err := redis.Int(c.Do("hexists", fmt.Sprintf("group:member:%d", gid), uid))
	if err != nil {
		return false, err
	}

	return res == 1, nil
}
