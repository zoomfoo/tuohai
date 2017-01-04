package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"tuohai/im_api/options"
)

func SyncFriends() error {
	fmt.Println("sync main site friends")
	sy, _ := GetSyncTime()
	if sy == 0 {
		sy = 1262275200 // 2010.01.01
	}
	last := strconv.Itoa(sy)
	now := strconv.Itoa(int(time.Now().Unix()) - 60)

	res, err := GetFrindFromMainSite(options.Opts.AuthHost, last)
	if err != nil || res == nil {
		fmt.Println("sync main site friends error: ", err)
		return err
	}
	var v map[string]interface{}
	err = json.Unmarshal(res, &v)
	if v["error_code"].(float64) == 0 {
		for _, e := range v["friends"].([]interface{}) {
			a := e.([]interface{})[0].(map[string]interface{})["id"].(string)
			b := e.([]interface{})[1].(map[string]interface{})["id"].(string)
			fmt.Printf("get two persons:%s::%s\n", a, b)
			rid := IsRelation(a, b, 0)
			if rid == "" {
				fmt.Printf("add new friend:%s::%s\n", a, b)
				rid, err := CreateRelation(a, b, 0, 0, "")
				if err != nil {
					fmt.Printf("add new friends fails:%s\n", err)
				}
				fmt.Printf("new friend rid:%s\n", rid)
			}
		}
	}
	SetSyncTime(now)
	//fmt.Println("unmarshal ", v)
	return nil
}
