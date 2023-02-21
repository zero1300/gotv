package model

import (
	"encoding/json"
	"time"
)

type Video struct {
	ID         uint      `gorm:"primaryKey" json:"id,string"`
	UID        uint      `json:"uid,string"`
	Title      string    `form:"title" binding:"required" json:"title"`
	Intro      string    `form:"intro" binding:"required" json:"intro"`
	Uri        string    `form:"uri" binding:"required" json:"uri"`
	Cover      string    `form:"cover" binding:"required" json:"cover"`
	Like       uint      `json:"like,string"`
	Duration   string    `json:"duration"`
	Views      uint      `json:"views,string"`
	CreateTime time.Time `json:"createTime" gorm:"autoCreateTime" `
}

func (Video) TableName() string {
	return "video"
}

func (v *Video) MarshalBinary() (data []byte, err error) {
	return json.Marshal(v)
}

func (v *Video) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, v)
}
