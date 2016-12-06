package models

import (
	"errors"
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

var RecordNotFound = errors.New("record not found")

type TblGroup struct {
	Id        int      `gorm:"column:id" json:"-"`
	Gid       string   `gorm:"column:gid" json:"gid"`
	Gname     string   `gorm:"column:gname" json:"g_name"`
	Creator   string   `gorm:"column:creator" json:"creator"`
	Admincnt  uint8    `gorm:"column:admincnt" json:"admin_cnt"`
	Membercnt uint     `gorm:"column:membercnt" json:"mem_cnt"`
	Version   uint     `gorm:"column:version" json:"version"`
	IsPublic  uint8    `gorm:"column:is_public" json:"is_public"`
	CreatedAt uint     `gorm:"column:created_at" json:"time"`
	UpdatedAt uint     `gorm:"column:updated_at" json:"-"`
	GroupMems []string `gorm:"-" json:"members" `
}

func (t *TblGroup) TableName() string {
	return "tbl_group"
}

func AddTblGroup(m *TblGroup) error {
	return db.Create(m).Error
}

//
func GetTblGroupById(gid string) (*TblGroup, error) {
	g := &TblGroup{}
	err := db.Find(g, "gid = ?", gid).Error
	if err != nil {
		return nil, err
	}
	mems, err := GroupMemsId(g.Gid)
	if err != nil {
		return nil, err
	}

	for _, mem := range mems {
		g.GroupMems = append(g.GroupMems, mem.Member)
	}
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
	g.CreatedAt = uint(now)
	g.UpdatedAt = uint(now)
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
