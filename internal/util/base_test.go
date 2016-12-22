package util

import (
	"testing"
)

func TestValidateEmail(t *testing.T) {
	t.Log(ValidateEmail("7072547@qq.com"))
	t.Log(ValidateEmail("tosingular@gmail.com"))
	t.Log(ValidateEmail("7072547@qq.co"))
	t.Log(ValidateEmail("7072547@qq1111.com"))
	t.Log(ValidateEmail("ajljlk"))
	t.Log(ValidateEmail("7072547ajljlk.com"))
	t.Log(ValidateEmail("7072547ajljlkqq.com"))
	t.Log(ValidateEmail("7072547ajljlkm"))
	t.Log(ValidateEmail("70725ajljlkm"))
	t.Log(ValidateEmail("7072ajljlkcom"))
	t.Log(ValidateEmail("7ajljlk.com"))
}
