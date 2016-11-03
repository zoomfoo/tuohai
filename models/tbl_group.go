package models

import (
	"fmt"
	"time"
)

const (
	OrdinaryMembers = iota //普通成员
	GroupAdmin             //群组管理员
	GroupCreator           //群组创建者/群主
)

const (
	Normal = iota //正常
	Quit          //已经退出群
)

type TblGroup struct {
	Id          int      `gorm:"column:id"`
	Gid         string   `gorm:"column:gid"`
	Gname       string   `gorm:"column:gname"`
	Creator     string   `gorm:"column:creator"`
	Admincnt    uint8    `gorm:"column:admincnt"`
	Membercnt   uint     `gorm:"column:membercnt"`
	Version     uint     `gorm:"column:version"`
	IsPublic    uint8    `gorm:"column:is_public"`
	CreatedTime uint     `gorm:"column:created_time"`
	UpdatedTime uint     `gorm:"column:updated_time"`
	GroupMems   []string `json:"group_member_list" gorm:"-"`
}

func (t *TblGroup) TableName() string {
	return "tbl_group"
}

// AddTblGroup insert a new TblGroup into database and returns
// last inserted Id on success.
func AddTblGroup(m *TblGroup) error {
	return db.Create(m).Error
}

// GetTblGroupById retrieves TblGroup by Id. Returns error if
// Id doesn't exist
func GetTblGroupById(gid string) (*TblGroup, error) {
	g := &TblGroup{}
	err := db.Find(g, "gid = ?", gid).Error
	return g, err
}

func CreateGroup(g *TblGroup) (*TblGroup, error) {
	var (
		val []interface{}
		now = time.Now().Unix()
	)

	tx := db.Begin()

	g.Membercnt = uint(len(g.GroupMems))
	g.Admincnt = 0
	g.CreatedTime = uint(now)
	g.UpdatedTime = uint(now)
	g.Version = 1

	resg := tx.Create(g)
	if err := resg.Error; err != nil {
		tx.Rollback()
		return nil, err
	} else {
		if err := resg.Find(g).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		//Create a group manager
		gm := &TblGroupMember{
			GroupId:   g.Gid,
			Member:    g.Creator,
			Role:      GroupCreator,
			Status:    Normal,
			CreatedAt: uint(time.Now().Unix()),
			UpdatedAt: uint(time.Now().Unix()),
		}

		if err := tx.Create(gm).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		val = append(val, gm.Member, now)

		//Add a group of members
		for _, mem := range g.GroupMems {
			//如果成员mem ==0 跳过
			if mem == "" {
				continue
			}
			//拼接redis 参数
			val = append(val, mem, now)

			if mem == g.Creator {
				continue
			}

			t := time.Now().Unix()
			gm := &TblGroupMember{
				GroupId:   g.Gid,
				Member:    mem,
				Status:    Normal,
				CreatedAt: uint(t),
				UpdatedAt: uint(t),
				Role:      OrdinaryMembers,
			}
			if err := tx.Create(gm).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		if err := tx.Commit().Error; err != nil {
			return nil, err
		}

		//查找创建完成的group
		group, err := GetTblGroupById(g.GroupMems[0])
		if err != nil {
			return nil, err
		}

		//将数据更新redis中
		c := rpool.Get()
		defer c.Close()
		c.Do("select", "5")
		if _, err := c.Do("hmset", append([]interface{}{fmt.Sprintf("group:member:%d", g.Id)}, val...)...); err != nil {
			return nil, err
		}

		return group, nil
	}
}

// GetAllTblGroup retrieves all TblGroup matches certain condition. Returns empty list if
// no records exist
func GetAllTblGroup(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	return nil, nil
}

// UpdateTblGroup updates TblGroup by Id and returns error if
// the record to be updated doesn't exist
func UpdateTblGroupById(m *TblGroup) (err error) {
	return nil
}

// DeleteTblGroup deletes TblGroup by Id and returns error if
// the record to be deleted doesn't exist
func DeleteTblGroup(id int) (err error) {
	return nil
}
