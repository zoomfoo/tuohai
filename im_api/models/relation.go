package models

import (
	"time"

	"tuohai/internal/console"
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
	Rtype        int    `gorm:"column:rtype"`
}

func (r *Relation) TableName() string {
	return "tbl_relation"
}

func Friends(uuid string) ([]Relation, error) {
	var r []Relation
	err := db.Find(&r, "status = 0 and (small_id = ? or big_id = ?)", uuid, uuid).Error
	return r, err
}

//uid用户id
//fuid 是cid 或者是rid是关系id
func Friend(uid, fuid string) (*Relation, error) {
	var rel Relation
	err := db.Find(&rel, "status = 0 and rid = ? and (small_id = ? or big_id = ?)", fuid, uid, uid).Error
	//fmt.Println("获取好友详细信息: ", rel)
	return &rel, err
}

func FriendSmallAndBig(uid, fid string) (*Relation, error) {
	var rel Relation
	small, big := convert.StringSortByRune(uid, fid)
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

func SyncCreateFriend(small, big string, fid int) (string, error) {
	return createRelation(small, big, fid, 0)
}

func createRelation(a, b string, fid, rtype int) (string, error) {
	small, big := convert.StringSortByRune(a, b)
	if cid := IsRelation(small, big, rtype); cid != "" {
		return cid, nil
	}
	prefix := "r_"
	if rtype == 1 {
		prefix = "t_"
	}
	cid := prefix + uuid.NewV4().StringMd5()
	r := &Relation{
		Rid:          cid,
		SmallId:      small,
		BigId:        big,
		OriginId:     small,
		Status:       0,
		SyncFriendId: fid,
		CreatedAt:    time.Now().Unix(),
		UpatedAt:     time.Now().Unix(),
		Rtype:        rtype,
	}
	tx := db.Begin()

	if err := tx.Create(r).Error; err != nil {
		tx.Rollback()
		return "", err
	}
	//同步到redis
	if err := saveChannelToRedis(cid, []string{small, big}); err != nil {
		tx.Rollback()
		return "", err
	}

	return cid, tx.Commit().Error
}

func CreateRelation(a, b string, rtype int) (string, error) {
	return createRelation(a, b, 0, rtype)
}

func DelRelation(cid string) error {
	r := &Relation{}
	tx := db.Begin()
	// 系统好友不能删除
	if err := tx.Find(r, "rid = ? and rtype != 2", cid).Error; err != nil {
		tx.Rollback()
		return err
	}

	err := tx.Table(r.TableName()).Where("rid = ?", cid).Updates(map[string]interface{}{"status": RelationDeleted}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	c := rpool.Get()
	defer c.Close()

	if _, err := c.Do("hdel", "channel:member:"+r.Rid, r.SmallId, r.BigId); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// 这个接口需要重构或者二次封装,以方便使用
func IsRelation(a, b string, rtype int) string {
	small, big := convert.StringSortByRune(a, b)
	r := &Relation{}
	err := db.Find(r, "small_id = ? and big_id = ? and status = 0 and rtype = ?", small, big, rtype).Error
	if err != nil {
		console.StdLog.Error(err)
		return ""
	}
	return r.Rid
}

func GetSysRid(sys, x string) string {
	small, big := convert.StringSortByRune(sys, x)
	r := &Relation{}
	err := db.Find(r, "small_id = ? and big_id = ? and status = 0 and rtype = 2", small, big).Error
	if err != nil {
		return ""
	}
	return r.Rid
}
