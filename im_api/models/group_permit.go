package models

import (
	"fmt"
)

type PermitType uint8

const (
	PERMIT_NO PermitType = 0
	PERMIT_OK PermitType = 1
)

type GroupPermit struct {
	GType       int8          `gorm:"column:group_type" json:"-"`
	RoleInfo    OperationVerb `gorm:"column:role_info" json:"-"`
	Permit      PermitType    `gorm:"column:permit" json:"-"`
	Role        uint8         `gorm:"column:role" json:"-"`
	Status      int8          `gorm:"column:status" json:"-"`
	Create_time int64         `gorm:"column:create_time" json:"-"`
	Update_time int64         `gorm:"column:update_time" json:"-"`
}

func (gp *GroupPermit) TableName() string {
	return "tbl_group_permit"
}

func (gp *GroupPermit) IsEditTitle() bool {
	return gp.RoleInfo == RENAME_GROUP && gp.Permit == PERMIT_OK
}

func (gp *GroupPermit) IsAddGroupMember() bool {
	return gp.RoleInfo == ADD_GROUP_MEMS && gp.Permit == PERMIT_OK
}

func (gp *GroupPermit) IsRemoveGroupMember() bool {
	return gp.RoleInfo == DEL_GROUP_MEMS && gp.Permit == PERMIT_OK
}

func (gp *GroupPermit) IsTransferGroupLeader() bool {
	return gp.RoleInfo == TRANSFER_GROUP && gp.Permit == PERMIT_OK
}

func (gp *GroupPermit) IsExitGroup() bool {
	return gp.RoleInfo == QUIT_GROUP && gp.Permit == PERMIT_OK
}

func (gp *GroupPermit) IsDismissGroup() bool {
	return gp.RoleInfo == DISMISS_GROUP && gp.Permit == PERMIT_OK
}

func Permit(gid, uid string, role_info OperationVerb) *GroupPermit {
	g := NewGroup(gid)
	ri := RoleInfo(gid, uid)
	//gtype 群组类型 role_info 群组权限动作 ri 为role 2是创建者
	fmt.Println("群组权限: ", g.GType, role_info, ri)
	gp := &GroupPermit{}
	err := db.Find(gp, "group_type = ? and role_info = ? and role = ? and status = 0", g.GType, role_info, ri).Error
	if err != nil {
		fmt.Println(err)
	}
	return gp
}
