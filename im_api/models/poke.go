package models

import (
	"time"
)

type TblChuoyixia struct {
	ID          int    `gorm:"primary_key" json:"-"`
	Chuoid      string `json:"poke_id"`
	Rcv         string `json:"rcv"`
	IsConfirmed int8   `json:"is_confirmed"`
	IsDelByRcv  int8   `json:"is_del_by_rcv"`
	CreatedAt   int    `json:"created"`
	ConfirmedAt int    `json:"-"`
}

func (t *TblChuoyixia) TableName() string {
	return "tbl_chuoyixia"
}

func ConfirmChuo(chid, uid string) error {
	t := &TblChuoyixia{}
	if err := db.Find(t, "chuoid = ? and rcv = ? and is_del_by_rcv = 0", chid, uid).Error; err != nil {
		return err
	}
	if t.IsConfirmed != 0 {
		return nil
	}
	tr := db.Begin()
	t.IsConfirmed = 1
	t.ConfirmedAt = int(time.Now().Unix())
	if err := tr.Save(t).Error; err != nil {
		tr.Rollback()
		return err
	}
	if err := tr.Exec("update tbl_chuoyixia_meta set confirmed_cnt=confirmed_cnt+1,updated_at=? where chuoid=?", time.Now().Unix(), chid).Error; err != nil {
		tr.Rollback()
		return err
	}
	tr.Commit()
	return nil
}

func GetChuo(chid string) ([]TblChuoyixia, error) {
	var t []TblChuoyixia
	if err := db.Find(&t, "chuoid = ?", chid).Error; err != nil {
		return nil, err
	}
	return t, nil
}

func GetChuoRcv(to string) ([]TblChuoyixia, error) {
	var rs []TblChuoyixia
	if err := db.Find(&rs, "rcv = ? and is_del_by_rcv = 0", to).Error; err != nil {
		return nil, err
	}
	return rs, nil
}
