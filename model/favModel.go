package model

import "time"

type Fav struct {
	UID        uint      `json:"uid"`
	VID        uint64    `gorm:"column:vid" json:"vid"`
	CreateTime time.Time `json:"-" gorm:"autoCreateTime" `
}
