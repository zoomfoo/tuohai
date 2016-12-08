package models

import (
	"time"

	"tuohai/internal/convert"
	"tuohai/internal/uuid"
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

func Friend(token, fuuid string) (*Relation, error) {
	var rel Relation
	small, big := "", ""
	if convert.RuneAccumulation(token) > convert.RuneAccumulation(fuuid) {
		small, big = fuuid, token
	} else {
		small, big = token, fuuid
	}

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
	r := &Relation{
		Rid:          "r_" + uuid.NewV4().StringMd5(),
		SmallId:      small,
		BigId:        big,
		OriginId:     small,
		Status:       0,
		SyncFriendId: fid,
		CreatedAt:    time.Now().Unix(),
		UpatedAt:     time.Now().Unix(),
	}
	return db.Create(r).Error
}

func CreateRelation(small, big string, fid int) error {
	return createRelation(small, big, 0)
}
