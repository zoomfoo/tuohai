package models

import (
	"fmt"
	"sort"
	"time"

	"tuohai/im_api/options"
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
	Way          int    `gorm:"column:way"`
	Note         string `gorm:"column:note"`
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
	return createRelation(small, big, fid, 0, 0, "")
}

func createRelation(a, b string, fid, rtype int, way int, note string) (string, error) {
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
		Way:          way,
		Note:         note,
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

func CreateRelation(a, b string, rtype int, way int, note string) (string, error) {
	return createRelation(a, b, 0, rtype, way, note)
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
// rtype,0：普通关系，1:临时好友，2:系统好友
func IsRelation(a, b string, rtype int) string {
	small, big := convert.StringSortByRune(a, b)
	r := &Relation{}
	err := db.Find(r, "small_id = ? and big_id = ? and status = 0 and rtype = ?", small, big, rtype).Error
	if err != nil {
		return ""
	}
	return r.Rid
}

func GetSysRid(sys, x string) string {
	small, big := convert.StringSortByRune(sys, x)
	r := &Relation{}
	err := db.Find(r, "small_id = ? and big_id = ? and status = 0 and rtype = 2", small, big).Error
	if err != nil || r == nil {
		rid, err := createRelation(sys, x, 0, 2, 0, "")
		if err != nil {
			fmt.Println("create system friend fails,sys:%s,uid:%s", sys, x)
			return ""
		}
		return rid
	}
	return r.Rid
}

func MatchFriends(uid string, ps []User) (map[string]interface{}, error) {
	pu := []string{}
	pf := map[string]string{}
	for i, _ := range ps {
		pu = append(pu, ps[i].Phone)
		ir := IsRelation(uid, ps[i].Uuid, 0)
		if ir != "" {
			pf[ps[i].Phone] = ir
		} else {
			go func() {
				pm := &PersonMatched{
					From:      uid,
					Partner:   ps[i].Uuid,
					Status:    10,
					CreatedAt: time.Now().Unix(),
					UpdatedAt: time.Now().Unix(),
				}
				err := AddPersonMatched(pm)
				if err != nil {
					fmt.Printf("save person matched error:%s", err)
				}
			}()
		}
	}
	ret := map[string]interface{}{
		"platform_users": pu,
		"friends":        pf,
	}
	return ret, nil
}

type NewPerson struct {
	Id        string `json:"id"`
	Uuid      string `json:"uuid"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Avatar    string `json:"avatar"`
	Email     string `json:"email"`
	Way       int    `json:"way"`
	Attach    string `json:"attach"`
	Status    int    `json:"status"`
	UpdatedAt int64  `json:"updated_at"`
	IsApply   int    `json:"is_apply"`
}

type NP []*NewPerson

func (n NP) Len() int {
	return len(n)
}

func (n NP) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

func (n NP) Less(i, j int) bool {
	return n[i].UpdatedAt > n[j].UpdatedAt
}

func NewPersons(uid, token string) ([]*NewPerson, error) {
	var np NP
	pm, err1 := GetPersonMatched(uid)
	if err1 != nil {
		fmt.Printf("get person matched error:%s", err1)
	}
	apply, err2 := GetFriendApplyes(uid)
	if err2 != nil {
		fmt.Printf("get friend apply error:%s", err2)
	}
	if err1 != nil && err2 != nil {
		return np, nil
	}
	if len(pm) == 0 && len(apply) == 0 {
		return np, nil
	}
	if len(pm) == 0 {
		for _, a := range apply {
			params := []string{"user_ids=" + a.ApplyUid}
			ua, err := GetBatchUsersFromMain(token, options.Opts.AuthHost, params)
			if err != nil {
				fmt.Printf("get user error,err:%s", err)
				continue
			}
			t := &NewPerson{
				Id:        a.Fid,
				Uuid:      a.ApplyUid,
				Name:      ua[0].Uname,
				Phone:     ua[0].Phone,
				Avatar:    ua[0].Avatar,
				Email:     ua[0].Email,
				Way:       int(a.Way),
				Attach:    a.Attach,
				Status:    int(a.Status),
				UpdatedAt: a.ConfirmTime,
				IsApply:   1,
			}
			np = append(np, t)
		}
	} else if len(apply) == 0 {
		for _, m := range pm {
			params := []string{"user_ids=" + m.Partner}
			ua, err := GetBatchUsersFromMain(token, options.Opts.AuthHost, params)
			if err != nil {
				fmt.Printf("get user error,err:%s", err)
				continue
			}
			t := &NewPerson{
				Uuid:      m.Partner,
				Name:      ua[0].Uname,
				Phone:     ua[0].Phone,
				Avatar:    ua[0].Avatar,
				Email:     ua[0].Email,
				Way:       int(AddressbookWay),
				Attach:    "",
				Status:    m.Status,
				UpdatedAt: m.UpdatedAt,
				IsApply:   0,
			}
			np = append(np, t)
		}
	} else {
		for _, a := range apply {
			params := []string{"user_ids=" + a.ApplyUid}
			ua, err := GetBatchUsersFromMain(token, options.Opts.AuthHost, params)
			if err != nil {
				fmt.Printf("get user error,err:%s", err)
				continue
			}
			t := &NewPerson{
				Id:        a.Fid,
				Uuid:      a.ApplyUid,
				Name:      ua[0].Uname,
				Phone:     ua[0].Phone,
				Avatar:    ua[0].Avatar,
				Email:     ua[0].Email,
				Way:       int(a.Way),
				Attach:    a.Attach,
				Status:    int(a.Status),
				UpdatedAt: a.ConfirmTime,
				IsApply:   1,
			}
			np = append(np, t)
		}
		for _, m := range pm {
			flag := false
			for _, a := range apply {
				if a.ApplyUid == m.Partner {
					flag = true
					break
				}
			}
			if !flag {
				params := []string{"user_ids=" + m.Partner}
				ua, err := GetBatchUsersFromMain(token, options.Opts.AuthHost, params)
				if err != nil {
					fmt.Printf("get user error,err:%s", err)
					continue
				}
				t := &NewPerson{
					Uuid:      m.Partner,
					Name:      ua[0].Uname,
					Phone:     ua[0].Phone,
					Avatar:    ua[0].Avatar,
					Email:     ua[0].Email,
					Way:       int(AddressbookWay),
					Attach:    "",
					Status:    m.Status,
					UpdatedAt: m.UpdatedAt,
					IsApply:   0,
				}
				np = append(np, t)
			}
		}
	}
	sort.Sort(np)
	return np, nil
}
