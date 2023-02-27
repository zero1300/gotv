package model

import "time"

type History struct {
	ID         uint      `gorm:"primaryKey" json:"id,string"`
	UID        uint      `gorm:"column:uid" json:"uid,string"`
	VID        string    `gorm:"column:vid" form:"vid"`
	Position   string    `json:"position" form:"position"`
	CreateTime time.Time `json:"createTime" gorm:"autoCreateTime" `
}

func (History) TableName() string {
	return "history"
}
