package auth

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"gopkg.in/gin-gonic/gin.v1"
	"tuohai/im_api/models"
	"tuohai/im_api/options"
	"tuohai/internal/convert"
	httplib "tuohai/internal/http"
)

func LoginAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ""
		token = ctx.Request.Header.Get("session_token")
		if token == "" {
			token = ctx.Query("session_token")
		}

		url := GetUserInfoUrl(token, options.Opts.AuthHost)
		if user, err := ValidationToken(url); err != nil || user == nil {
			ctx.Abort()
			fmt.Println(err)
			ctx.JSON(http.StatusUnauthorized, gin.H{"err_code": 1, "data": "无权限访问"})
		} else {
			//每一个新的连接会创建一个goroutine
			//每一个新的请求都会新创建Context
			//所以Context类下面的key是不需要加锁的
			//这里user可以放心使用
			fmt.Println("user id: ", user.Uid)
			models.ValidAndCreate(&models.User{
				Uuid:         user.Uid,
				Uname:        user.Nickname,
				Phone:        user.Phone,
				Email:        user.Email,
				Token:        token,
				IsFirstlogin: 1,
			})
			ctx.Set("user", user)
			ctx.Set("token", token)
			ctx.Next()
		}
	}
}

//主站用户体系
type MainUser struct {
	Id       int64  `json:"id"`
	Uid      string `json:"uuid"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
}

func ValidationToken(url string) (*MainUser, error) {
	var result struct {
		Msg       string    `json:"msg"`
		MainUser  *MainUser `json:"user"`
		ErrorCode float64   `json:"error_code"`
	}
	fmt.Println("auth URL: ", url)
	if url == "" {
		return nil, fmt.Errorf("%v", "url is empty")
	}
	err := httplib.Get(url).ToJson(&result)
	if err != nil {
		return nil, err
	}

	if result.ErrorCode == 0 {
		return result.MainUser, nil
	}
	return nil, fmt.Errorf("session check ret:%s", result.Msg)
}

//获取个人profile URL
func GetUserInfoUrl(token, url string) string {
	return fmt.Sprintf("%s/api/i/session?%s", url, SignStr(token))
}

//获取更新用户信息URL
func GetUpdateUserInfoUrl(token, url string, params []string) string {
	return fmt.Sprintf("%s/api/i/profile?%s", url, SignStr(token, params...))
}

//获取主站好友列表URL
func GetFriendsUrl(token, url string, params []string) string {
	return fmt.Sprintf("%s/api/i/friends?%s", url, SignStr(token, params...))
}

//批量获取用户信息
func GetBatchUsersUrl(token, url string, params []string) string {
	return fmt.Sprintf("%s/api/v1.1/users/info?%s", url, SignStr(token, params...))
}

//获取邮箱服务URL
func GetEmailServiceURL(token, url string, params []string) string {
	return fmt.Sprintf("%s/api/i/email?%s", url, SignStr(token, params...))
}

//获取短信服务URL
func GetSMSServiceURL(token, url string, params []string) string {
	return fmt.Sprintf("%s/api/i/sms?%s", url, SignStr(token, params...))
}

func GetBatchUsers(token, url string, params []string) ([]models.User, error) {
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
	fmt.Println("URL: ", GetBatchUsersUrl(token, url, params))
	err := httplib.Get(GetBatchUsersUrl(token, url, params)).ToJson(&result)
	if err != nil {
		return nil, err
	}
	var users []models.User
	for i := 0; i < len(result.MainUser); i++ {
		localuser, _ := models.GetUserById(result.MainUser[i].Id)
		u := models.User{
			Uuid:   result.MainUser[i].Id,
			Uname:  result.MainUser[i].Name,
			Avatar: result.MainUser[i].Avatar,
			Phone:  result.MainUser[i].Phone,
			Email:  result.MainUser[i].Email,
			Desc:   localuser.Desc,
			Yltype: result.MainUser[i].Yltype,
		}
		users = append(users, u)
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

//发短信
func SendSMS(url, token string, param []string) (bool, error) {
	smsurl := url + "/api/i/sms"
	fmt.Println("sms url: ", smsurl)
	req := httplib.Post(smsurl)
	req.Header("content-type", "application/x-www-form-urlencoded")
	req.Body(SignStr(token, param...))
	js := make(map[string]interface{})
	if err := req.ToJson(&js); err != nil {
		return false, err
	}
	fmt.Println("sms ret: ", js)
	if error_code, ok := js["error_code"].(float64); ok {
		return error_code == 0, nil
	} else {
		return false, fmt.Errorf("%s", "服务器错误参数")
	}
}

//发邮件
func SendEmail(url, token string, param []string) (bool, error) {
	email := GetEmailServiceURL(token, url, param)
	fmt.Println("email: ", email)
	data := httplib.Post(email)
	js := make(map[string]interface{})
	if err := data.ToJson(&js); err != nil {
		return false, err
	}

	if error_code, ok := js["error_code"].(float64); ok {
		return error_code == 0, nil
	} else {
		return false, fmt.Errorf("%s", "服务器错误参数")
	}
}
