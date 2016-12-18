package auth

import (
	"testing"
)

func TestValidationToken(t *testing.T) {
	url := GetUserInfoUrl("W8cyc7VgjkF6ghJMhUPQMvPr5ed38cef2ba3bdea", "http://test.yunwoke.com")
	t.Log(ValidationToken(url))
}

func TestSendSMS(t *testing.T) {
	t.Log(SendSMS("378X6fgBDsSF2Vt3TP3Y25M4d6b2d0f878b81c5d", "http://test.yunwoke.com",
		[]string{"phones=15040565139", "content=我测试看见了不用回复", "site=yunliao", "user_id=21"},
	))
}

func TestSendEmail(t *testing.T) {
	t.Log(SendEmail("378X6fgBDsSF2Vt3TP3Y25M4d6b2d0f878b81c5d", "http://test.yunwoke.com",
		[]string{"email=7072547@qq.com",
			"title=asdfas",
			"content=232323"},
	))
}
