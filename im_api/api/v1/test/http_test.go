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
