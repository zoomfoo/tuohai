package models

import (
	"errors"
	"fmt"
	"sync"
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

const (
	DelGroup = 1 //删除群
)

var RecordNotFound = errors.New("record not found")

var GroupPool *sync.Pool

func init() {
	GroupPool = &sync.Pool{
		New: func() interface{} {
			return new(Group)
		},
	}
}

type Group struct {
	Id          int      `gorm:"column:id" json:"-"`
	Gid         string   `gorm:"column:gid" json:"gid"`
	Gname       string   `gorm:"column:gname" json:"g_name"`
	Creator     string   `gorm:"column:creator" json:"creator"`
	Admincnt    uint8    `gorm:"column:admincnt" json:"admin_cnt"`
	Membercnt   uint     `gorm:"column:membercnt" json:"mem_cnt"`
	Version     uint     `gorm:"column:version" json:"version"`
	IsPublic    uint8    `gorm:"column:is_public" json:"is_public"`
	CreatedTime int64    `gorm:"column:created_at" json:"time"`
	UpdatedTime int64    `gorm:"column:updated_at" json:"-"`
	GroupMems   []string `gorm:"-" json:"members" `
}

func (g *Group) TableName() string {
	return "tbl_group"
}

func (g *Group) GetGroupById() error {
	if g.Gid == "" {
		return errors.New("gid is empty")
	}
	return db.Find(g, "gid = ?", g.Gid).Error
}

func (g *Group) validation() bool {
	return true
}

func (g *Group) reset() {
	g.Id = 0
	g.Gid = ""
	g.Gname = ""
	g.Creator = ""
	g.Admincnt = 0
	g.Membercnt = 0
	g.Version = 0
	g.IsPublic = 0
	g.CreatedTime = 0
	g.UpdatedTime = 0
	g.GroupMems = g.GroupMems[:]
}

func (g *Group) GetGroupMemsForSQL() error {
	var (
		mems []GroupMember
	)
	err := db.Find(&mems, "gid = ? and status = 0", g.Gid).Error
	if err != nil {
		return err
	}

	for _, mem := range mems {
		g.GroupMems = append(g.GroupMems, mem.Member)
	}
	return nil
}

func (g *Group) GetGroupMemsForNOSQL() error {
	return nil
}

func (g *Group) RenameGroup(name string) error {
	if !g.validation() {
		return errors.New("group a lack of important value")
	}
	if g.Gname == name {
		return nil
	}
	return db.Save(g).Error
}

func (g *Group) TransferGroup(prince, king uint32) error {
	return nil
}

func AddGroup(m *Group) error {
	return db.Create(m).Error
}

//修改群名√
func RenameGroup(gid, newname string) error {
	g := GroupPool.Get().(*Group)
	g.reset()
	g.Gid = gid

	if err := g.GetGroupById(); err != nil {
		return err
	}
	if err := g.RenameGroup(newname); err != nil {
		return err
	}
	GroupPool.Put(g)
	return nil
}

//获取群组信息√
func GetGroupById(gid string) (*Group, error) {
	g := GroupPool.Get().(*Group)
	g.reset()
	g.Gid = gid

	if err := g.GetGroupById(); err != nil {
		return nil, err
	}
	if err := g.GetGroupMemsForSQL(); err != nil {
		return nil, err
	}
	GroupPool.Put(g)
	return g, nil
}

//获取群组列表√
func GetGroupsByUid(uid string) ([]Group, error) {
	var (
		groups []Group
		gs     []string
	)

	mems, err := GroupMemByUid(uid)
	if err != nil {
		return nil, err
	}
	fmt.Println("mems : ", mems)

	for _, m := range mems {
		gs = append(gs, m.GroupId)
	}

	if err := db.Find(&groups, "gid in (?)", gs).Error; err != nil {
		return nil, err
	}

	for i := 0; i < len(groups); i++ {
		if err := groups[i].GetGroupMemsForSQL(); err != nil {
			fmt.Println(err)
		}
	}

	return groups, nil
}

//创建群√
func CreateGroup(g *Group) (*Group, error) {
	var (
		val []interface{}
		now = time.Now().Unix()
	)

	tx := db.Begin()

	g.Membercnt = uint(len(g.GroupMems))
	g.CreatedTime = now
	g.UpdatedTime = now
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
		gm := &GroupMember{
			GroupId:   g.Gid,
			Member:    g.Creator,
			Role:      GroupCreator,
			Status:    Normal,
			CreatedAt: now,
			UpdatedAt: now,
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

			gm := &GroupMember{
				GroupId:   g.Gid,
				Member:    mem,
				Role:      OrdinaryMembers,
				Status:    Normal,
				CreatedAt: now,
				UpdatedAt: now,
			}
			if err := tx.Create(gm).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		//将数据更新redis中
		c := rpool.Get()
		defer c.Close()
		c.Do("select", "5")
		if _, err := c.Do("hmset", append([]interface{}{fmt.Sprintf("group:member:%s", g.Gid)}, val...)...); err != nil {
			tx.Rollback()
			return nil, err
		}

		if err := tx.Commit().Error; err != nil {
			return nil, err
		}

		//查找创建完成的group
		if err := g.GetGroupById(); err != nil {
			return nil, err
		}

		return g, nil
	}
}

//群组管理
//群组管理
//群组管理

//禅让群主 新加的接口
func TransferGroup(prince, king, groupId uint32) error {
	return nil
}

//解散群组√
func DismissGroup(gid string) error {
	var (
		gm    *GroupMember
		g     Group
		mlist []string
	)

	//设置群组status 为1 删除
	tx := db.Begin()
	err := tx.Table(g.TableName()).Where("gid = ?", gid).Updates(map[string]interface{}{"status": DelGroup}).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	//将群成员表status设置为退出
	if err := tx.Table(gm.TableName()).Where("gid = ?", gid).
		Updates(map[string]interface{}{"status": Quit}).Error; err != nil {
		tx.Rollback()
		return err
	}

	mems, err := GroupMemsId(gid)
	if err != nil {
		tx.Rollback()
		return err
	}

	for i := 0; i < len(mems); i++ {
		mlist = append(mlist, mems[i].Member)
	}

	//将redis群组数据清空，群组未读清空
	if _, err := QuitGroup(gid, mlist); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
