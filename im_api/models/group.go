package models

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"tuohai/internal/uuid"
)

const (
	ORDINARY_MEMS = 0 //普通成员
	GROUP_ADMIN   = 1 //群组管理员
	GROUP_CREATOR = 2 //群组创建者/群主
)

const (
	Normal = iota //正常
	Quit          //已经退出群
)

const (
	DelGroup = 1 //删除群
)

type OperationVerb int8

const (
	ADD_GROUP_MEMS OperationVerb = 1 //添加群成员
	DEL_GROUP_MEMS OperationVerb = 2 //删除群成员
	ADD_ADMIN      OperationVerb = 3 //添加管理员
	DEL_ADMIN      OperationVerb = 4 //删除管理员
	RENAME_GROUP   OperationVerb = 5 //重命名
	QUIT_GROUP     OperationVerb = 6 //從群組成員中退出
	DISMISS_GROUP  OperationVerb = 7 //關閉群組并解散群成員
	TRANSFER_GROUP OperationVerb = 8 //轉讓群組
)

type GroupType int8

const (
	NORMAL_GROUP  GroupType = 1
	Project_Group GroupType = 3
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
	Id          int       `gorm:"column:id" json:"-"`
	Gid         string    `gorm:"column:gid" json:"gid"`
	Gname       string    `gorm:"column:gname" json:"name" form:"name"`
	Creator     string    `gorm:"column:creator" json:"creator"`
	Admincnt    uint8     `gorm:"column:admincnt" json:"admin_cnt"`
	Membercnt   uint      `gorm:"column:membercnt" json:"mem_cnt"`
	Version     uint      `gorm:"column:version" json:"version"`
	GType       GroupType `gorm:"column:type" json:"-"`
	IsPublic    uint8     `gorm:"column:is_public" json:"is_public"`
	Status      uint8     `gorm:"column:status" json:"-"`
	CreatedTime int64     `gorm:"column:created_at" json:"time"`
	UpdatedTime int64     `gorm:"column:updated_at" json:"-"`
	GroupMems   []string  `gorm:"-" json:"member" form:"member"`
}

func NewGroup(Gid string) Group {
	g := GroupPool.Get().(*Group)
	g.reset()
	g.Gid = Gid
	GroupPool.Put(g)

	g.GetGroupById()
	return *g
}

func initGroup() *Group {
	g := &Group{}
	g.Gid = "g_" + uuid.NewV4().StringMd5()
	g.CreatedTime = time.Now().Unix()
	g.UpdatedTime = time.Now().Unix()
	g.Status = 0
	g.Version = 1
	return g
}

func (g *Group) TableName() string {
	return "tbl_group"
}

func (g *Group) GetGroupById() error {
	if g.Gid == "" {
		return errors.New("gid is empty")
	}
	return db.Find(g, "gid = ? and status = 0", g.Gid).Error
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
	g.GroupMems = g.GroupMems[0:0]
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

	g.Gname = name
	return db.Save(g).Error
}

func (g *Group) TransferGroup(prince, king uint32) error {
	return nil
}

//修改群名√
func RenameGroup(gid, newname string) error {
	g := NewGroup(gid)
	if err := g.GetGroupById(); err != nil {
		return err
	}
	if err := g.RenameGroup(newname); err != nil {
		return err
	}

	return nil
}

//获取群组信息√
func GetGroupById(gid string) (*Group, error) {
	g := NewGroup(gid)

	if err := g.GetGroupById(); err != nil {
		return nil, err
	}
	if err := g.GetGroupMemsForSQL(); err != nil {
		return nil, err
	}
	return &g, nil
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

	if err := db.Find(&groups, "gid in (?) and status = 0", gs).Error; err != nil {
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
func CreateGroup(creator, gname string, members []string) (*Group, error) {
	var (
		now = time.Now().Unix()
	)

	tx := db.Begin()

	g := initGroup()
	g.Membercnt = uint(len(members) + 1)
	g.Creator = creator
	g.Gname = gname

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
			Role:      GROUP_CREATOR,
			Status:    Normal,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := tx.Create(gm).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		fmt.Println("members: ", members)
		//Add a group of members
		for _, mem := range members {
			//如果成员mem ==0 跳过
			if mem == "" {
				continue
			}

			if mem == g.Creator {
				continue
			}

			gm := &GroupMember{
				GroupId:   g.Gid,
				Member:    mem,
				Role:      ORDINARY_MEMS,
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
		if err := saveChennelToRedis(g.Gid, append(members, creator)); err != nil {
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
		g.GetGroupMemsForSQL()

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
