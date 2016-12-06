package models

import (
	"crypto/md5"
	"encoding/hex"
)

type User struct {
	Id    int    `gorm:"column:id" json:"-"`
	Uuid  string `gorm:"column:uuid" json:"uuid"`
	Uname string `gorm:"column:uname" json:"name"`
	Desc  string `gorm:"column:description" json:"desc"` //个性签名
	Token string `gorm:"column:token" json:"token"`
}

func (t *User) TableName() string {
	return "tbl_user"
}

func GetUserById(uuid string) (*User, error) {
	user := &User{}
	err := db.Find(user, "uuid = ?", uuid).Error
	return user, err
}

func GetUserByIds(uuids []string) ([]User, error) {
	var users []User
	err := db.Find(&users, "uuid in (?)", uuids).Error
	return users, err
}

func Login(uname, pwd string) (*User, error) {
	var user User
	err := db.Find(&user, "uname = ? and passwd = ?", uname, pwd).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateOrCreateUser(u *User) error {
	// user,err:=GetUserById(u.Uuid)
	return nil
}

func generateToken(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}
