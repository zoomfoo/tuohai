package models

import (
	"fmt"
	"time"
)

type GroupMember struct {
	Id        int    `gorm:"column:id"`
	GroupId   string `gorm:"column:gid"`
	Member    string `gorm:"column:member"`
	Role      uint8  `gorm:"column:role"`
	Status    uint8  `gorm:"column:status"`
	CreatedAt int64  `gorm:"column:created_at"`
	UpdatedAt int64  `gorm:"column:updated_at"`
}

func (t *GroupMember) TableName() string {
	return "tbl_group_member"
}

func (gm *GroupMember) GetGroupMemberByUidAndGid(UserId, GroupId string) error {
	return db.Find(gm, "userId = ? and groupId = ? and status = 0", UserId, GroupId).Error
}

func (gm *GroupMember) IsCreator() bool {
	return gm.Role == GroupCreator
}

func (gm *GroupMember) IsAdmin() bool {
	return gm.Role == GroupAdmin
}

func IsCreator(UserId, GroupId string) (bool, error) {
	gm := new(GroupMember)
	err := gm.GetGroupMemberByUidAndGid(UserId, GroupId)
	if err != nil {
		return false, err
	}
	return gm.IsCreator(), nil
}

func IsAdmin(UserId, GroupId string) (bool, error) {
	gm := new(GroupMember)
	err := gm.GetGroupMemberByUidAndGid(UserId, GroupId)
	if err != nil {
		return false, err
	}
	return gm.IsAdmin(), nil
}

func GetGroupMemberById(id int) (*GroupMember, error) {
	mem := &GroupMember{}
	err := db.Find(mem, "id = ?", id).Error
	return mem, err
}

func GroupMemsId(gid string) ([]GroupMember, error) {
	var (
		mems []GroupMember
	)
	err := db.Find(&mems, "gid = ? and status = 0", gid).Error
	return mems, err
}

func GroupMemByUid(uid string) ([]GroupMember, error) {
	var mems []GroupMember
	err := db.Find(&mems, "member = ? and status = 0", uid).Error
	return mems, err
}

func AssociationGroups(uid string) ([]GroupMember, error) {
	var mems []GroupMember
	err := db.Table((&GroupMember{}).TableName()).Where("`member` = ?", uid).Scan(&mems).Error
	return mems, err
}

//添加群成员√
func AddGroupMember(gid string, GroupMems []string) (*Group, error) {
	var (
		val []interface{}
		now = time.Now().Unix()
		g   = &Group{Gid: gid}
	)

	if err := g.GetGroupById(); err != nil {
		return nil, err
	}

	if len(GroupMems) == 0 {
		return nil, fmt.Errorf("ERROR: %s", "GroupMems length is zero")
	}

	tx := db.Begin()
	//遍历群成员
	for _, mem := range GroupMems {
		//生成redis key
		val = append(val, mem, now)

		gm := &GroupMember{}
		//将status为1 和 0状态的的数据都查出来
		tx.Where("gid = ? and member = ? ", gid, mem).Find(gm)
		//如果Status = Normal 那么判断用户已经存在
		if gm.GroupId != "" && gm.Status == Normal {
			return nil, fmt.Errorf("%s 已经存在", mem)
		}
		//如果 Status = Quit 那么就将Status 修改为 Normal
		if gm.Status == Quit {
			tx.Table(gm.TableName()).Where("gid = ? and member = ?", g.Gid, mem).Updates(map[string]interface{}{"status": Normal, "updated_at": time.Now().Unix()})
			continue
		}

		if err := tx.Create(&GroupMember{
			GroupId:   g.Gid,
			Member:    mem,
			Status:    Normal,
			Role:      OrdinaryMembers,
			CreatedAt: now,
			UpdatedAt: now,
		}).Error; err != nil {
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

	if err := db.Table(g.TableName()).Where("id = ?", gid).
		Updates(map[string]interface{}{"membercnt": len(GroupMems)}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	//查找创建完成的group
	group, err := GetGroupById(g.Gid)
	if err != nil {
		return nil, err
	}

	return group, nil
}

//删除成员√
func DelGroupMember(gid string, GroupMems []string) (*Group, error) {
	var (
		val []interface{}
		gm  GroupMember
	)

	if len(GroupMems) == 0 {
		return nil, fmt.Errorf("ERROR: %s", "GroupMems length is zero")
	}

	tx := db.Begin()
	err := tx.Table(gm.TableName()).
		Where("gid = ? and member in (?)", gid, GroupMems).
		Updates(map[string]interface{}{"status": Quit, "updated_at": time.Now().Unix()}).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	//移除成员并移除成员群组所在的session列表
	for i := 0; i < len(GroupMems); i++ {
		session := (&Session{From: GroupMems[i]})
		fmt.Println(session.TableName())
		if err := tx.Table(session.TableName()).Where("`to` = ? and `from` = ?", gid, GroupMems[i]).
			Updates(map[string]interface{}{"status": deleted}).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	for _, mem := range GroupMems {
		val = append(val, mem)
	}

	c := rpool.Get()
	defer c.Close()
	c.Do("select", "5")
	if _, err := c.Do("hdel", append([]interface{}{fmt.Sprintf("group:member:%s", gid)}, val...)...); err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	//查找创建完成的group
	group, err := GetGroupById(gid)
	if err != nil {
		return nil, err
	}

	return group, nil
}

//添加或删除管理员
//is = true 添加管理员
//is = false 删除管理员
func ManageAdmin(operator uint32, gid string, mems []string, is bool) (*Group, error) {
	// var mem int
	// if is {
	// 	mem = GroupAdmin
	// } else {
	// 	mem = OrdinaryMembers
	// }

	// err := db.Table("IMGroupMember").
	// 	Where("groupId = ? and userId in (?)", gid, mems).
	// 	Updates(map[string]interface{}{"role": mem, "updated": time.Now().Unix()}).Error
	// if err != nil {
	// 	return nil, err
	// }

	// group, err := GroupInfo(operator, g.Id)
	// if err != nil {
	// 	return nil, err
	// }

	// return group, nil
	return nil, nil
}

func GetMyGroupId(id string) ([]string, error) {
	var (
		gm  []GroupMember
		ids []string
	)
	err := db.Find(&gm, "member = ? and status = 0", id).Error
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(gm); i++ {
		ids = append(ids, gm[i].GroupId)
	}

	return ids, nil
}
