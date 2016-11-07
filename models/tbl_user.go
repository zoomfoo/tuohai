package models

import (
	"crypto/md5"
	"encoding/hex"
)

type TblUser struct {
	Id    int    `gorm:"column:id" json:"-"`
	Uuid  string `gorm:"column:uuid" json:"uuid"`
	Uname string `gorm:"column:uname" json:"name"`
	Token string `gorm:"column:token" json:"token"`
}

func (t *TblUser) TableName() string {
	return "tbl_user"
}

func GetTblUserById(uuid string) (*TblUser, error) {
	user := &TblUser{}
	err := db.Find(user, "uuid = ?", uuid).Error
	return user, err
}

func GetTblUserByIds(uuids []string) ([]TblUser, error) {
	var users []TblUser
	err := db.Find(&users, "uuid in (?)", uuids).Error
	return users, err
}

func Login(uname, pwd string) (*TblUser, error) {
	var user TblUser
	err := db.Find(&user, "uname = ? and passwd = ?", uname, pwd).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func generateToken(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}
