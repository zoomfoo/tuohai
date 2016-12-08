package models

type ApplyType int8

const (
	UntreatedApply ApplyType = 0
	AgreedApply    ApplyType = 1
	RefusedApply   ApplyType = 2
)

type ApplyWay int8

const (
	FriendSeek ApplyWay = 0
	GroupSeek  ApplyWay = 1
)

type FriendApply struct {
	Id          string    `gorm:"column:id" json:"id"`
	ApplyUid    string    `gorm:"column:apply_uid" json:"uuid"`
	TargetUid   string    `gorm:"column:target_uid" json:"-"`
	Way         ApplyWay  `gorm:"column:way" json:"way"`
	Attach      string    `gorm:"column:attach" json:"attach"`
	Status      ApplyType `gorm:"column:status" json:"status"`
	LaunchTime  int64     `gorm:"column:launch_time" json:"time"`
	ConfirmTime int64     `gorm:"column:confirm_time" json:"-"`
}

func (fa *FriendApply) TableName() string {
	return "friend_apply"
}

func (fa *FriendApply) ValidationField() string {
	if fa.Id == "" {
		return "id 不能为空"
	}

	if fa.ApplyUid == "" {
		return "ApplyUid 不能为空"
	}

	if fa.TargetUid == "" {
		return "TargetUid 不能为空"
	}

	switch fa.Way {
	case FriendSeek, GroupSeek:
	default:
		return "未知的Way值"
	}

	switch fa.Status {
	case UntreatedApply, AgreedApply, RefusedApply:
	default:
		return "未知的Status值"
	}
	return ""
}

func FriendApplyById(id string) (*FriendApply, error) {
	apply := &FriendApply{}
	err := db.Find(apply, "id = ?", id).Error
	return apply, err
}

//获取自己的好友申请列表
func FriendApplys(uid string) ([]FriendApply, error) {
	var applys []FriendApply
	err := db.Find(&applys, "target_uid = ?", uid).Error
	return applys, err
}

//更新
func SaveFriendApply(apply *FriendApply) error {
	return db.Table(apply.TableName()).Where("id = ?", apply.Id).Updates(apply).Error
}

//创建
func CreateFriendApply(apply *FriendApply) error {
	return db.Create(apply).Error
}
