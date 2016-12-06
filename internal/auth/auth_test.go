package auth

import (
	"testing"
)

func TestValidationToken(t *testing.T) {
	url := GetUserInfoUrl("W8cyc7VgjkF6ghJMhUPQMvPr5ed38cef2ba3bdea", "http://test.yunwoke.com")
	t.Log(ValidationToken(url))
}
