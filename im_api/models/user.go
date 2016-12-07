package models

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

type User struct {
	Id    int64  `gorm:"primary_key" json:"-"`
	Uuid  string `gorm:"column:uuid" json:"uuid"`
	Uname string `gorm:"column:uname" json:"name"`
	Desc  string `gorm:"column:description" json:"desc"` //个性签名
	Token string `gorm:"column:token" json:"token"`
}

func (t *User) TableName() string {
	return "tbl_user"
}

func (u *User) IsUserNotExist() bool {
	return u.Uuid == ""
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

func ValidAndCreate(u *User) error {
	user, _ := GetUserById(u.Uuid)
	if user == nil || user.Uuid == "" {
		//用户不存在
		fmt.Println("创建用户")
		return db.Create(u).Error
	}
	return nil
}

func SaveUser(u *User) error {
	if u.Uuid == "" {
		return fmt.Errorf("%s", "uuid is empty")
	}
	fmt.Println("Save User: ", *u)
	return db.Table(u.TableName()).Where("uuid = ?", u.Uuid).
		Updates(u).Error
}

func generateToken(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}
