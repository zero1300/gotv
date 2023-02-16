package model

type File struct {
	ID       uint   `gorm:"primaryKey" json:"id,string"`
	Filename string `json:"filename"`
	Path     string `json:"path"`
	Key      string `json:"key"`
}

func (File) TableName() string {
	return "file"
}
