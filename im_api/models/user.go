package models

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"tuohai/im_api/options"
)

type User struct {
	Id           int64  `gorm:"column:id" json:"-"`
	Uuid         string `gorm:"column:uuid" json:"uuid"`
	Uname        string `gorm:"column:uname" json:"name"`
	Phone        string `gorm:"column:phone" json:"phone"`
	Email        string `gorm:"column:email" json:"email"`
	Avatar       string `gorm:"-" json:"avatar"`
	Desc         string `gorm:"column:description" json:"desc"` //个性签名
	Token        string `gorm:"column:token" json:"-"`
	IsFirstlogin int    `gorm:"column:is_fristlogin" json:"is_firstlogin"`
	Yltype       int    `gorm:"-" json:"yltype"`
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
	//fmt.Println(*user)
	if user == nil || user.Uuid == "" {
		//用户不存在
		fmt.Println("create new user in yunliao")
		err := db.Create(u).Error
		// 自动添加系统好友
		go func() {
			_, err := createRelation(options.Opts.SysUserYunliao, user.Uuid, 0, 2, 0, "")
			if err != nil {
				fmt.Println("add system relation fails")
			}
			_, err = createRelation(options.Opts.SysUserClouderwork, user.Uuid, 0, 2, 0, "")
			if err != nil {
				fmt.Println("add system relation fails")
			}
		}()
		return err
	} else {
		go func() {
			GetSysRid(options.Opts.SysUserYunliao, user.Uuid)
			GetSysRid(options.Opts.SysUserClouderwork, user.Uuid)
		}()
	}
	return nil
}

func SaveUser(u *User) error {
	if u.Uuid == "" {
		return fmt.Errorf("%s", "uuid is empty")
	}
	fmt.Println("Save User: ", *u)
	return db.Table(u.TableName()).Where("uuid = ?", u.Uuid).
		Updates(map[string]interface{}{"description": u.Desc, "is_firstlogin": u.IsFirstlogin, "uname": u.Uname}).Error
}

func CreateUser(u *User) error {
	return db.Create(u).Error
}

func generateToken(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}

func GetBatchUsers(uids []string) ([]User, error) {
	var us []User
	err := db.Find(&us, "uuid in (?)", uids).Error
	return us, err
}

func SelectUsers(u *User) ([]User, error) {
	var users []User
	err := db.Table(u.TableName()).Where(u).Scan(&users).Error
	return users, err
}

func GetAllUsers() ([]User, error) {
	var users []User
	err := db.Find(&users).Error
	return users, err
}

func GetUserByPhones(phones []string) ([]User, error) {
	var users []User
	err := db.Find(&users, "phone in (?)", phones).Error
	return users, err
}
