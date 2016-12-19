package test

import (
	"testing"
	"tuohai/im_api/models"
)

//获取用户所在群列表
func TestGetGroupsByUid(t *testing.T) {
	t.Log(models.GetGroupsByUid("b56486af2d149c3b816d593bf0e4f1b5"))
}

// func TestCreateGroup(t *testing.T) {
// 	t.Log(models.CreateGroup(&models.Group{
// 		Gid:       "7072547",
// 		Gname:     "王阳",
// 		Creator:   "b56486af2d149c3b816d593bf0e4f1b5",
// 		IsPublic:  0,
// 		GroupMems: []string{"202cb962ac59075b964b07152d234b71"},
// 	}))
// }

func TestAddGroupMember(t *testing.T) {
	t.Log(models.AddGroupMember("7072547", []string{"1ecb3665e88cdb174dec77863438662a"}))
}

func TestDelGroupMember(t *testing.T) {
	t.Log(models.DelGroupMember("7072547", []string{"1ecb3665e88cdb174dec77863438662a", "202cb962ac59075b964b07152d234b71"}))
}

func TestSyncMysqlToRedis(t *testing.T) {
	models.SyncMysqlToRedis()
}
