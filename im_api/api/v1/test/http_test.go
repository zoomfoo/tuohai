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
	api := URI + "/v1/friends?session_token=QJ2NvhcWj8N2Hek4FY5jg9G80cbc5c28796cebd7"

	v := make(url.Values)
	v.Add("uuid", "94de7d8b6a2bf757")
	v.Add("attach", "asdfasdfas")
	v.Add("way", "0")

	payload := strings.NewReader(v.Encode())
	buf := &bytes.Buffer{}

	req, _ := http.NewRequest("POST", api, payload)

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	buf.ReadFrom(res.Body)

	t.Log(buf.String())
}
