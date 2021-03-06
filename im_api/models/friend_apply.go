package models

import (
	"fmt"
	"time"

	"tuohai/internal/console"
	"tuohai/internal/convert"
)

type ApplyType int8

const (
	UntreatedApply ApplyType = 0
	AgreedApply    ApplyType = 1
	RefusedApply   ApplyType = 2
)

type ApplyWay int8

const (
	UuidWay        ApplyWay = 0
	PhoneWay       ApplyWay = 1
	EmailWay       ApplyWay = 2
	TeamWay        ApplyWay = 3
	AddressbookWay ApplyWay = 4
	ClouderworkWay ApplyWay = 5
	OthersWay      ApplyWay = 6
)

type FriendApply struct {
	Id          int       `gorm:"column:id" json:"-"`
	Fid         string    `gorm:"column:fid" json:"id"`
	ApplyUid    string    `gorm:"column:apply_uid" json:"uuid"`
	TargetUid   string    `gorm:"column:target_uid" json:"target"`
	Way         ApplyWay  `gorm:"column:way" json:"way"`
	Attach      string    `gorm:"column:attach" json:"attach"`
	Status      ApplyType `gorm:"column:status" json:"status"`
	LaunchTime  int64     `gorm:"column:launch_time" json:"time"`
	ConfirmTime int64     `gorm:"column:confirm_time" json:"confirm_time"`
	Note        string    `gorm:"column:note" json:"note"`
}

func (fa *FriendApply) TableName() string {
	return "tbl_friend_apply"
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

	if int(fa.Way) > 6 || int(fa.Way) < 0 {

		return "未知的Way值 (添加途径(0:uuid,1: 手机号,2: 邮箱,3,团队,4:通讯录,5:云沃客,6:其他))"
	}

	switch fa.Status {
	case UntreatedApply, AgreedApply, RefusedApply:
	default:
		return "未知的Status值 状态 (0， 未处理， 1 通过，2 拒绝)"
	}
	return ""
}

func FriendApplyById(fid, uid string) (*FriendApply, error) {
	apply := &FriendApply{}
	err := db.Find(apply, "fid = ?  and target_uid = ? and status= 0", fid, uid).Error
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
	} else {
		err := db.Where("target_uid = ? and status in (?)", uid, status).Order("launch_time desc").Find(&applys).Error
		return applys, err
	}

	err := db.Offset(pageindex).Limit(pagesize).Where("target_uid = ? and status in (?)", uid, status).Order("launch_time desc").Find(&applys).Error
	return applys, err
}

//is=true 已经处理
//is=false 未处理
func FriendApplysCount(uid string, is bool) int {
	var total int
	var applys []FriendApply
	var status []int
	if is {
		status = []int{1, 2}
	} else {
		status = []int{0}
	}
	err := db.Where("target_uid = ? and status in (?)", uid, status).Find(&applys).Count(&total).Error
	if err != nil {
		console.StdLog.Error(err)
		return 0
	}
	return total
}

//
func ProcessFriendApply(apply *FriendApply) (string, error) {
	apply.ConfirmTime = time.Now().Unix()

	tx := db.Begin()
	if err := tx.Table(apply.TableName()).Where("fid = ?", apply.Fid).Updates(apply).Error; err != nil {
		tx.Rollback()
		return "", err
	}
	if apply.Status == AgreedApply {
		small, big := convert.StringSortByRune(apply.ApplyUid, apply.TargetUid)
		fmt.Println(small, big)

		if cid := IsRelation(small, big, 0); cid != "" {
			return cid, tx.Commit().Error
		}

		if cid, err := CreateRelation(small, big, 0, int(apply.Way), apply.Note); err != nil {
			tx.Rollback()
			return "", err
		} else {
			return cid, tx.Commit().Error
		}
	}
	return "", nil
}

//创建
func CreateFriendApply(apply *FriendApply) error {
	ns := &FriendApply{}
	err := db.Find(ns, "target_uid = ? and apply_uid = ? and status= 0", apply.TargetUid, apply.ApplyUid).Error
	if err != nil {
		fmt.Println("friend apply create")
		return db.Create(apply).Error
	} else {
		fmt.Println("friend apply update")
		apply.Id = ns.Id
		return db.Save(apply).Error
	}
}

func GetFriendApplyes(uuid string) ([]FriendApply, error) {
	var fa []FriendApply
	err := db.Order("confirm_time desc").Find(&fa, "target_uid = ?", uuid).Error
	return fa, err
}
