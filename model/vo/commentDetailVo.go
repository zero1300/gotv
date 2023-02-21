package vo

import "gotv/model"

type CommentDetailVo struct {
	model.Comment
	Nickname string `json:"nickname" gorm:"-"`
	Avatar   string `json:"avatar" gorm:"-"`
}
