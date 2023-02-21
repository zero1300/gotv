package vo

import "gotv/model"

type VideoVo struct {
	model.Video
	Nickname string `gorm:"-" json:"nickname"`
}
