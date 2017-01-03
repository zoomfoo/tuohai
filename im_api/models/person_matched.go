package models

import (
//"fmt"
)

type PersonMatched struct {
	Id        int64  `gorm:"column:id" json:"-"`
	From      string `gorm:"column:from" json:"from"`
	Partner   string `gorm:"column:partner" json:"partner"`
	Status    int    `gorm:"column:status" json:"status"`
	CreatedAt int64  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt int64  `gorm:"column:updated_at" json:"updated_at"`
}

func (p *PersonMatched) TableName() string {
	return "tbl_persons_matched"
}

func AddPersonMatched(pm *PersonMatched) error {
	tpm := &PersonMatched{}
	err := db.Find(tpm, "`from` = ? and partner = ?", pm.From, pm.Partner).Error
	if err != nil {
		return db.Create(pm).Error
	}
	return nil
}

func UpdatePersonMatched(pm *PersonMatched) error {
	tpm := &PersonMatched{}
	err := db.Find(tpm, "`from` = ? and partner = ?", pm.From, pm.Partner).Error
	if err != nil {
		return db.Create(pm).Error
	} else {
		pm.Id = tpm.Id
		return db.Save(pm).Error
	}
	return nil
}

func GetPersonMatched(uuid string) ([]PersonMatched, error) {
	var pms []PersonMatched
	err := db.Order("updated_at desc").Find(&pms, "`from` = ?", uuid).Error
	if err != nil {
		return nil, err
	}
	return pms, nil
}
