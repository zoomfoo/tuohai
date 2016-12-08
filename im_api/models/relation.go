package models

import (
	"fmt"
	"time"

	"tuohai/internal/convert"
	"tuohai/internal/uuid"
)

const (
	RelationDeleted = 1
)

type Relation struct {
	Id           int    `gorm:"column:id"`
	Rid          string `gorm:"column:rid"`
	SmallId      string `gorm:"column:small_id"`
	BigId        string `gorm:"column:big_id"`
	OriginId     string `gorm:"column:origin_id"`
	Status       int    `gorm:"column:status"`
	SyncFriendId int    `gorm:"column:sync_friend_id"`
	CreatedAt    int64  `gorm:"column:created_at"`
	UpatedAt     int64  `gorm:"column:upated_at"`
}

func (r *Relation) TableName() string {
	return "tbl_relation"
}

func Friends(uuid string) ([]Relation, error) {
	var r []Relation
	err := db.Find(&r, "status = 0 and (small_id = ? or big_id = ?)", uuid, uuid).Error
	return r, err
}

func Friend(uid, fuid string) (*Relation, error) {
	var rel Relation
	small, big := convert.StringSort(uid, fuid)
	fmt.Println(small, big)
	err := db.Find(&rel, "status = 0 and small_id = ? and big_id = ?", small, big).Error
	return &rel, err
}

func GetMyRelationId(id string) ([]string, error) {
	var (
		ids []string
	)
	rel, err := Friends(id)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(rel); i++ {
		ids = append(ids, rel[i].Rid)
	}

	return ids, nil
}

func SyncCreateFriend(small, big string, fid int) error {
	return createRelation(small, big, fid)
}

func createRelation(small, big string, fid int) error {
	cid := "r_" + uuid.NewV4().StringMd5()
	r := &Relation{
		Rid:          cid,
		SmallId:      small,
		BigId:        big,
		OriginId:     small,
		Status:       0,
		SyncFriendId: fid,
		CreatedAt:    time.Now().Unix(),
		UpatedAt:     time.Now().Unix(),
	}
	tx := db.Begin()

	if err := tx.Create(r).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := saveChennelToRedis(cid, []string{small, big}); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func CreateRelation(small, big string) error {
	return createRelation(small, big, 0)
}

func DelRelation(small, big string) error {
	r := &Relation{}
	tx := db.Begin()
	err := tx.Table(r.TableName()).Where("small_id = ? and big_id = ?", small, big).Updates(map[string]interface{}{"status": RelationDeleted}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	c := rpool.Get()
	defer c.Close()
	return nil
	// c.Do(commandName, ...)
}
