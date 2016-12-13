package models

import ()

type FileType int8

type TransactHandle interface {
	Path() (string, error)
}

const (
	FileTypeImage      FileType = 1
	FileTypeDoc        FileType = 2
	FileTypeCompressed FileType = 3
	FileTypeSnippet    FileType = 4
)

type FileInfo struct {
	Id       string   `gorm:"column:id" json:"-"`
	To       string   `gorm:"column:to" json:"cid"`
	Name     string   `gorm:"column:name" json:"name"`
	Path     string   `gorm:"column:path" json:"path"`
	Size     int      `gorm:"column:size" json:"size"`
	Type     FileType `gorm:"column:type" json:"type"`
	Ext      string   `gorm:"column:ext" json:"ext"`
	Category string   `gorm:"column:category" json:"category"`
	Meta     *Image   `gorm:"column:meta" json:"meta"`
	Creator  string   `gorm:"column:creator" json:"creator"`
	Updated  int64    `gorm:"column:updated" json:"-"`
	Created  int64    `gorm:"column:created" json:"time"`
}

type Image struct {
}

func (file *FileInfo) TableName() string {
	return "file_info"
}

func (file *FileInfo) GetFileInfo() {
	if file.Type == FileTypeImage {
		file.GetImage()
	}
}

//to 等于cid
func (file *FileInfo) GetFilesInfo(to []string) ([]FileInfo, error) {
	var infos []FileInfo
	err := db.Find(&infos, "`to` in (?) and `status` = 0", to).Error
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(infos); i++ {
		if file.Type == FileTypeImage {
			infos[i].GetImage()
		}
	}
	return infos, err
}

func (file *FileInfo) GetImage() error {
	var img *Image
	err := db.Find(img, "id = ?", file.Id).Error
	if err != nil {
		return err
	}
	if img == nil {
		file.Meta = &Image{}
	} else {
		file.Meta = img
	}
	return nil
}

func WriteFileToDB(file *FileInfo, transact TransactHandle) error {
	tx := db.Begin()

	if path, err := transact.Path(); err != nil {
		tx.Rollback()
		return err
	} else {
		file.Path = path
	}

	if err := tx.Create(file).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func GetFilesInfo(to []string) ([]FileInfo, error) {
	fi := &FileInfo{}
	return fi.GetFilesInfo(to)
}
