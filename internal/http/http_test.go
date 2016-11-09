package http

import (
	"testing"
)

func TestUsers(t *testing.T) {
	url := "http://test.yunwoke.com:10011/v1/user/202cb962ac59075b964b07152d234b71?session_token=202cb962ac59075b964b07152d234b71"
	t.Log(Users(url))
}
