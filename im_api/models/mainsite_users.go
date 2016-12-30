package models

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"tuohai/internal/convert"
	httplib "tuohai/internal/http"
)

//批量获取用户信息
func GetBatchUsersUrl(token, url string, params []string) string {
	return fmt.Sprintf("%s/api/v1.1/users/info?%s", url, SignStr(token, params...))
}

func GetBatchUsersFromMain(token, url string, params []string) ([]User, error) {
	var result struct {
		Msg      string `json:"msg"`
		MainUser []struct {
			Id     string `json:"id"`
			Name   string `json:"name"`
			Avatar string `json:"avatar"`
			Phone  string `json:"phone"`
			Email  string `json:"email"`
			Yltype int    `json:"yltype"`
		} `json:"users"`
		ErrorCode float64 `json:"error_code"`
	}
	//fmt.Println("URL: ", GetBatchUsersUrl(token, url, params))
	err := httplib.Get(GetBatchUsersUrl(token, url, params)).ToJson(&result)
	if err != nil {
		return nil, err
	}
	var users []User
	for i := 0; i < len(result.MainUser); i++ {
		u := User{
			Uuid:   result.MainUser[i].Id,
			Uname:  result.MainUser[i].Name,
			Avatar: result.MainUser[i].Avatar,
			Phone:  result.MainUser[i].Phone,
			Email:  result.MainUser[i].Email,
			Yltype: result.MainUser[i].Yltype,
		}
		users = append(users, u)
	}
	if len(users) == 0 {
		return nil, errors.New("none this user")
	}
	return users, nil
}

//生成签名参数
//params a=1212   b=2323   c=name
//url 后面作为参数
func SignStr(token string, params ...string) (session_token string) {
	var (
		par_str = strings.Join(params, "&")
		stamp   = convert.ToStr(time.Now().Unix())
	)
	if par_str != "" {
		par_str = "&" + par_str
	}
	params = append(
		params,
		fmt.Sprintf("session_token=%s", token),
		fmt.Sprintf("stamp=%s", stamp),
	)
	sort.Strings(params)
	sign_str := fmt.Sprintf("%scloudwork", strings.Join(params, ""))
	//fmt.Println("签名字符串: ", sign_str)
	session_token = fmt.Sprintf("session_token=%s&stamp=%s&sign=%s%s", token, stamp, getSign(sign_str), par_str)
	return
}

func getSign(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}
