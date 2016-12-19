package test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

const (
	URI = "http://127.0.0.1:10011"
)

// /v1/groups?session_token=202cb962ac59075b964b07152d234b71
func TestGroups(t *testing.T) {
	u := fmt.Sprintf("%s/v1/groups?session_token=202cb962ac59075b964b07152d234b71", URI)
	res, err := http.Get(u)
	if err != nil {
		t.Error(err)
		return
	}
	buf := bytes.Buffer{}
	defer res.Body.Close()
	buf.ReadFrom(res.Body)
	t.Log(string(buf.Bytes()))
}

func TestRenameGroup(t *testing.T) {
	u := fmt.Sprintf("%s/v1/groups/g_xdfed12138/rename/zhoujielunnnnn?session_token=202cb962ac59075b964b07152d234b71", URI)
	req, err := http.NewRequest("PUT", u, nil)
	if err != nil {
		t.Error(err)
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
		return
	}
	defer res.Body.Close()

	buf := &bytes.Buffer{}
	buf.ReadFrom(res.Body)
	t.Log(string(buf.Bytes()))
}

//解散群组
func TestDismissGroup(t *testing.T) {

}

//添加好友
func TestAddFriend(t *testing.T) {
	api := URI + "/v1/friends"

	v := make(url.Values)
	v.Add("uuid", "94de7d8b6a2bf757")
	v.Add("attach", "asdfasdfas")
	v.Add("way", "0")

	payload := strings.NewReader(v.Encode())
	buf := &bytes.Buffer{}

	req, _ := http.NewRequest("POST", api, payload)
	req.Header.Add("session_token", "QJ2NvhcWj8N2Hek4FY5jg9G80cbc5c28796cebd7")
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	buf.ReadFrom(res.Body)

	t.Log(buf.String())
}

//戳一下

//获取用户申请列表
func TestApplyFriends(t *testing.T) {
	api := URI + "/v1/apply/friends/is/1/100"
	req, _ := http.NewRequest("GET", api, nil)

	req.Header.Add("session_token", "b9hvLB7RFDQ85vC9TcT77W5t635a3945b2241400")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	buf := &bytes.Buffer{}
	buf.ReadFrom(res.Body)

	t.Log(buf.String())
}

//同意好友申请
func TestAgreeApplyFriend(t *testing.T) {

	api := URI + "/v1/apply/friends"
	v := make(url.Values)
	v.Add("id", "2b1e8840618bb2ba0ef4b7425bd012aa")
	v.Add("status", "1")

	payload := strings.NewReader(v.Encode())

	req, _ := http.NewRequest("PUT", api, payload)

	req.Header.Add("session_token", "b9hvLB7RFDQ85vC9TcT77W5t635a3945b2241400")
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	buf := &bytes.Buffer{}
	buf.ReadFrom(res.Body)

	t.Log(buf.String())
}

func TestCreateProjectGroup(t *testing.T) {
	url := "http://127.0.0.1:10011/project/groups?session_token=KufHBfhc3Rnr7AhEXm7M8qLv5574dca8e85bc887"

	payload := strings.NewReader("-----011000010111000001101001\r\nContent-Disposition: form-data; name=\"creator\"\r\n\r\n9e3a7e7659f4e865\r\n-----011000010111000001101001\r\nContent-Disposition: form-data; name=\"name\"\r\n\r\n测试团队看看\r\n-----011000010111000001101001\r\nContent-Disposition: form-data; name=\"member\"\r\n\r\nb45b00a270b5ac6d,9e3a7e7659f4e865\r\n-----011000010111000001101001--")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "multipart/form-data; boundary=---011000010111000001101001")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	t.Log(res)
	t.Log(string(body))
}
