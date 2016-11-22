package models

import ()

type TblChuoyixiaMeta struct {
	ID            int    `gorm:"primary_key" json:"-"`
	Chuoid        string `json:"chuoid"`
	Sender        string `json:"from"`
	Urgent        int8   `json:"urgent"`
	TotalCnt      int8   `json:"total_cnt"`
	ConfirmedCnt  int8   `json:"confirmed_cnt"`
	Cid           string `json:"cid"`
	MsgId         string `json:"msg_id"`
	Content       string `json:"content"`
	IsDelBySender int8   `json:"-"`
	CreatedAt     int    `json:"created_at"`
	UpdatedAt     int    `json:"-"`
}

func (t *TblChuoyixiaMeta) TableName() string {
	return "tbl_chuoyixia_meta"
}

func AddChuo(t *TblChuoyixiaMeta, tos []string) error {
	tr := db.Begin()
	if err := tr.Create(t).Error; err != nil {
		tr.Rollback()
		return err
	}
	for _, to := range tos {
		te := &TblChuoyixia{
			Chuoid:    t.Chuoid,
			Rcv:       to,
			CreatedAt: t.CreatedAt,
		}
		if err := tr.Create(te).Error; err != nil {
			tr.Rollback()
			return err
		}
	}
	tr.Commit()
	return nil
}

func GetChuoMeta(chid string) (*TblChuoyixiaMeta, error) {
	t := &TblChuoyixiaMeta{}
	if err := db.Find(t, "chuoid = ? and is_del_by_sender = 0", chid).Error; err != nil {
		return nil, err
	}
	return t, nil
}

func GetChuoFrom(from string) ([]TblChuoyixiaMeta, error) {
	var t []TblChuoyixiaMeta
	if err := db.Find(&t, "sender = ? and is_del_by_sender = 0", from).Error; err != nil {
		return nil, err
	}
	return t, nil
}
