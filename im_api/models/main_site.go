package models

import (
	"fmt"
	// "sort"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"tuohai/internal/convert"
)

var (
	msdb *gorm.DB
)

type MainSiteFriend struct {
	Id             int       `gorm:"column:id"`
	UserId         int       `gorm:"column:user_id"`
	ClientId       int       `gorm:"column:client_id`
	Status         string    `gorm:"column:status`
	CreateAt       time.Time `gorm:"column:create_at`
	TargetUserId   int       `gorm:"column:target_user_id`
	TargetClientId int       `gorm:"column:target_client_id`
}

type MainSiteUser struct {
	Id       int    `gorm:"column:id"`
	Phone    string `gorm:"column:phone"`
	Email    string `gorm:"column:email"`
	Nickname string `gorm:"column:nickname"`
	Avatar   string `gorm:"column:avatar"`
	Uuid     string `gorm:column:uuid`
}

func (f *MainSiteFriend) TableName() string {
	return "friend"
}

func InitMainSiteDB(MysqlOptions string) error {
	var err error
	msdb, err = gorm.Open("mysql", MysqlOptions)
	if err != nil {
		return err
	}

	msdb.DB().SetMaxOpenConns(100)
	msdb.DB().SetMaxIdleConns(10)
	return nil
}

func SyncFriends() error {
	var (
		friends []MainSiteFriend
		rel     Relation
	)
	//获取本地数据库增量friend id
	if err := db.Order("sync_friend_id desc").Limit(1).Find(&rel).Error; err != nil {
		if err.Error() == "record not found" {
			rel.SyncFriendId = 0
		} else {
			return err
		}
	}
	fmt.Println("SyncFriendId: ", rel.SyncFriendId)
	//获取主站好友列表
	err := msdb.Find(&friends, "user_id != 0 and id > ?", rel.SyncFriendId).Error
	if err != nil {
		return err
	}

	//这里使用range 因为friends里面保存对象这里会导致 值拷贝。
	//优化建议使用for i:=0;i<count;i++ 代替
	for _, friend := range friends {
		user := &MainSiteUser{}
		team := &MainSiteUser{}
		if err := msdb.Table("user").Where("id = ?", friend.UserId).Scan(user).
			Error; err != nil {
			fmt.Println(err)
			continue
		}
		if err := msdb.Table("team").Where("id = ?", friend.TargetClientId).Scan(team).
			Error; err != nil {
			fmt.Println(err)
			continue
		}

		if err := CreateUser(&User{
			Uuid:   user.Uuid,
			Uname:  user.Nickname,
			Avatar: user.Avatar,
			Phone:  user.Phone,
			Email:  user.Email,
		}); err != nil {
			fmt.Println(err)
			continue
		}

		small, big := convert.StringSortByRune(user.Uuid, team.Uuid)
		fmt.Println(small, big, friend.Id)
		if _, err := SyncCreateFriend(small, big, friend.Id); err != nil {
			fmt.Println(err)
		}
	}
	return nil
}
