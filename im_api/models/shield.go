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
		err = db.Create(shield).Error
	} else {
		shield.Id = ns.Id
		err = db.Save(shield).Error
	}
	// delete session
	go func() {
		session := &Session{From: uid}
		//fmt.Println(session.TableName())
		err := db.Table(session.TableName()).Where("`to` = ? and `from` = ?", cid, uid).
			Updates(map[string]interface{}{"status": deleted, "updated_at": time.Now().Unix()}).Error
		if err != nil {
			fmt.Printf("tmp session delete fails,from:%s,to:%s\n", uid, cid)
		}
	}()
	return err
}
