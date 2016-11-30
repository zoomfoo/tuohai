package test

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
)

var url = "http://127.0.0.1:10011"

// /v1/groups?session_token=202cb962ac59075b964b07152d234b71
func TestGroups(t *testing.T) {
	u := fmt.Sprintf("%s/v1/groups?session_token=202cb962ac59075b964b07152d234b71", url)
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
	u := fmt.Sprintf("%s/v1/groups/g_xdfed12138/rename/zhoujielunnnnn?session_token=202cb962ac59075b964b07152d234b71", url)
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
