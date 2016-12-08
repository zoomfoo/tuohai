package test

import (
	"testing"
	"tuohai/im_api/models"
)

func TestMainSite(t *testing.T) {
	models.SyncFriends()
}
