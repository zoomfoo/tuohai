package models

import (
	"fmt"
	"time"
)

type Shield struct {
	Id         int    `gorm:"column:id" json:"-"`
	Cid        string `gorm:"column:cid" json:"cid"`
	ShieldedBy string `gorm:"column:shielded_by" json:"shielded_by"`
	IsActived  int    `gorm:"cloumn:is_actived" json:"is_actived"`
	CreatedAt  int64  `gorm:"cloumn:created_at" json:"created_at"`
	UpdatedAt  int64  `gorm:"cloumn:updated_at" json:"updated_at"`
}

func (s *Shield) TableName() string {
	return "tbl_shieldlist"
}

func ShieldProcess(cid, uid string, flag int) error {
	shield := &Shield{
		Cid:        cid,
		ShieldedBy: uid,
		IsActived:  flag,
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	}
	ns := &Shield{}
	err := db.Table(ns.TableName()).Where("cid = ? and shielded_by = ?", cid, uid).Scan(ns).Error
	if err != nil {
		fmt.Println("no find ")
		return db.Create(shield).Error
	} else {
		shield.Id = ns.Id
		return db.Save(shield).Error
	}
}
