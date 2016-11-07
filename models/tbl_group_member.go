package models

import ()

type TblGroupMember struct {
	Id        int    `gorm:"column:id"`
	GroupId   string `gorm:"column:gid"`
	Member    string `gorm:"column:member"`
	Role      uint8  `gorm:"column:role"`
	Status    uint8  `gorm:"column:status"`
	CreatedAt uint   `gorm:"column:created_at"`
	UpdatedAt uint   `gorm:"column:updated_at"`
}

func (t *TblGroupMember) TableName() string {
	return "tbl_group_member"
}

func GetTblGroupMemberById(id int) (*TblGroupMember, error) {
	mem := &TblGroupMember{}
	err := db.Find(mem, "id = ?", id).Error
	return mem, err
}

func GroupMemsId(gid string) ([]TblGroupMember, error) {
	var (
		mems []TblGroupMember
	)
	err := db.Find(&mems, "gid = ? and status = 0", gid).Error
	return mems, err
}

func AssociationGroups(uid string) ([]TblGroupMember, error) {
	var mems []TblGroupMember
	err := db.Table((&TblGroupMember{}).TableName()).Where("`member` = ?", uid).Scan(&mems).Error
	return mems, err
}
