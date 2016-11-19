package models

type TblChuoyixia struct {
	ID          int `gorm:"primary_key"`
	Chuoid      string
	To          string
	IsConfirmed int8
	IsDelByRcv  int8
	CreatedAt   int
	ConfirmedAt int
}

func (t *TblChuoyixia) TableName() string {
	return "tbl_chuoyixia"
}

func AddChuo(t *TblChuoyixia) error {
	return db.Create(t).Error
}
