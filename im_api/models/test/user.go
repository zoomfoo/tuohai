package test

import (
	"testing"
	"tuohai/im_api/models"
)

func TestSelectUsers(t *testing.T) {
	t.Log(models.SelectUsers(&models.User{Phone: "13301330131"}))
}
