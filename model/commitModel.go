package model

import "time"

type Comment struct {
	ID         uint      `gorm:"primaryKey" json:"id,string"`
	UID        uint      `json:"uid,string"`
	VID        uint      `form:"vid" gorm:"column:vid" json:"vid,string" binding:"required"`
	Content    string    `form:"content" json:"content" binding:"required"`
	Like       uint      `json:"like"`
	DisLike    uint      `json:"dislike"`
	CreateTime time.Time `json:"createTime" gorm:"autoCreateTime"`
}

func (Comment) TableName() string {
	return "comment"
}
