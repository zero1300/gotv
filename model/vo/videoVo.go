package vo

import "gotv/model"

type VideoVo struct {
	model.Video
	Nickname         string `gorm:"-" json:"nickname"`
	CreateTimeString string `gorm:"-" json:"createTime"`
	Comments         int64  `gorm:"-" json:"comments"`
}
