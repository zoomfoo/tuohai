package models

import ()

type TblUser struct {
	Id    int    `gorm:"column:id" json:"-"`
	Uuid  string `gorm:"column:uuid" json:"uuid"`
	Uname string `gorm:"column:uname" json:"name"`
}

func (t *TblUser) TableName() string {
	return "tbl_user"
}

// AddTblUser insert a new TblUser into database and returns
// last inserted Id on success.
func AddTblUser(m *TblUser) (id int64, err error) {
	return 0, nil
}

// GetTblUserById retrieves TblUser by Id. Returns error if
// Id doesn't exist
func GetTblUserById(uuid string) (*TblUser, error) {
	user := &TblUser{}
	err := db.Find(user, "uuid = ?", uuid).Error
	return user, err
}

// GetAllTblUser retrieves all TblUser matches certain condition. Returns empty list if
// no records exist
func GetAllTblUser(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	return nil, nil
}

// UpdateTblUser updates TblUser by Id and returns error if
// the record to be updated doesn't exist
func UpdateTblUserById(m *TblUser) (err error) {
	return nil
}

// DeleteTblUser deletes TblUser by Id and returns error if
// the record to be deleted doesn't exist
func DeleteTblUser(id int) (err error) {
	return nil
}
