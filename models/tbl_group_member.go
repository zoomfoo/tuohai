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

// AddTblGroupMember insert a new TblGroupMember into database and returns
// last inserted Id on success.
func AddTblGroupMember(m *TblGroupMember) (id int64, err error) {
	return 0, nil
}

// GetTblGroupMemberById retrieves TblGroupMember by Id. Returns error if
// Id doesn't exist
func GetTblGroupMemberById(id int) (*TblGroupMember, error) {
	mem := &TblGroupMember{}
	err := db.Find(mem, "id = ?", id).Error
	return mem, err
}

func AssociationGroups(uid string) ([]TblGroupMember, error) {
	var mems []TblGroupMember
	err := db.Table((&TblGroupMember{}).TableName()).Where("`member` = ?", uid).Scan(&mems).Error
	return mems, err
}

// GetAllTblGroupMember retrieves all TblGroupMember matches certain condition. Returns empty list if
// no records exist
func GetAllTblGroupMember(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	return nil, nil
}

// UpdateTblGroupMember updates TblGroupMember by Id and returns error if
// the record to be updated doesn't exist
func UpdateTblGroupMemberById(m *TblGroupMember) (err error) {
	return nil
}

// DeleteTblGroupMember deletes TblGroupMember by Id and returns error if
// the record to be deleted doesn't exist
func DeleteTblGroupMember(id int) (err error) {
	return nil
}
