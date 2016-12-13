package models

import (
// "github.com/garyburd/redigo/redis"
)

type ApplyType int8

const (
	UntreatedApply ApplyType = 0
	AgreedApply    ApplyType = 1
	RefusedApply   ApplyType = 2
)

type ApplyWay int8

const (
	FriendSeek ApplyWay = 0
	GroupSeek  ApplyWay = 1
)

type FriendApply struct {
	Id          int       `gorm:"column:id" json:"-"`
	Fid         string    `gorm:"column:fid" json:"id"`
	ApplyUid    string    `gorm:"column:apply_uid" json:"uuid"`
	TargetUid   string    `gorm:"column:target_uid" json:"-"`
	Way         ApplyWay  `gorm:"column:way" json:"way"`
	Attach      string    `gorm:"column:attach" json:"attach"`
	Status      ApplyType `gorm:"column:status" json:"status"`
	LaunchTime  int64     `gorm:"column:launch_time" json:"time"`
	ConfirmTime int64     `gorm:"column:confirm_time" json:"-"`
}

func (fa *FriendApply) TableName() string {
	return "friend_apply"
}

func (fa *FriendApply) ValidationField() string {
	if fa.Fid == "" {
		return "id 不能为空"
	}

	if fa.ApplyUid == "" {
		return "apply uid 不能为空"
	}

	if fa.TargetUid == "" {
		return "uuid 不能为空"
	}

	switch fa.Way {
	case FriendSeek, GroupSeek:
	default:
		return "未知的Way值 (添加途径(0: 通过账号查找，1: 通过群组添加))"
	}

	switch fa.Status {
	case UntreatedApply, AgreedApply, RefusedApply:
	default:
		return "未知的Status值 状态 (0， 未处理， 1 通过，2 拒绝)"
	}
	return ""
}

func FriendApplyById(id string) (*FriendApply, error) {
	apply := &FriendApply{}
	err := db.Find(apply, "id = ?", id).Error
	return apply, err
}

//获取自己的好友申请列表
//is=true 已经处理
//is=false 未处理
func FriendApplys(uid string, is bool, pageindex, pagesize int) ([]FriendApply, error) {
	var applys []FriendApply
	var status []int
	if is {
		status = []int{1, 2}
	} else {
		status = []int{0}
	}

	if pageindex != 0 && pagesize != 0 {
		pageindex = (pageindex - 1) * pagesize
	}

	err := db.Offset(pageindex).Limit(pagesize).Where("target_uid = ? and status in (?)", uid, status).Find(&applys).Error
	return applys, err
}

//
func SaveFriendApply(apply *FriendApply) error {
	tx := db.Begin()
	if err := tx.Table(apply.TableName()).Where("fid = ?", apply.Fid).Updates(apply).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := saveChennelToRedis(apply.Fid, []string{apply.ApplyUid, apply.TargetUid}); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

//创建
func CreateFriendApply(apply *FriendApply) error {
	return db.Create(apply).Error
}
