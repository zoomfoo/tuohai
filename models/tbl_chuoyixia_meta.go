package models

type TblChuoyixiaMeta struct {
	ID            int    `gorm:"primary_key" json:"-"`
	Chuoid        string `json:"chuoid"`
	From          string `json:"from"`
	Urgent        int8   `json:"urgent"`
	TotalCnt      int8   `json:"total_cnt"`
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

func AddChuoMeta(t *TblChuoyixiaMeta) error {
	return db.Create(t).Error
}

func GetChuoMeta(chid string) (*TblChuoyixiaMeta, error) {
	t := &TblChuoyixiaMeta{}
	if err := db.Find(t, "chuoid = ? and is_del_by_sender = 0", chid).Error; err != nil {
		return nil, err
	}
	return t, nil
}
