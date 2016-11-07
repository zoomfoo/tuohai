package models

type Relation struct {
	Id        int    `gorm:"column:id"`
	Rid       string `gorm:"column:rid"`
	SmallId   string `gorm:"column:small_id"`
	BigId     string `gorm:"column:big_id"`
	OriginId  string `gorm:"column:origin_id"`
	Status    int    `gorm:"column:status"`
	CreatedAt int64  `gorm:"column:created_at"`
	UpatedAt  int64  `gorm:"column:upated_at"`
}

func (r *Relation) TableName() string {
	return "tbl_relation"
}

func Friends(uuid string) ([]Relation, error) {
	var r []Relation
	err := db.Find(&r, "status = 0 and (small_id == ? and big_id = ?)", uuid, uuid).Error
	return r, err
}
