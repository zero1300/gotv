package video

import (
	"fmt"
	"gotv/model"
	"gotv/resp"

	"github.com/bwmarrin/snowflake"
	"gorm.io/gorm"
)

type VideoDao struct {
	db *gorm.DB
}

func NewVideoDao(db *gorm.DB) *VideoDao {
	return &VideoDao{
		db: db,
	}
}

func (v *VideoDao) addVideo(video model.Video) {
	node, _ := snowflake.NewNode(1)
	id := node.Generate().Int64()
	video.ID = uint(id)
	v.db.Create(&video)
}

func (v *VideoDao) getVideoListByUid(uid uint, p model.Page) (resp.Pager, error) {
	videos := make([]model.Video, 10)
	var total int64
	err := v.db.Debug().Model(model.Video{}).Order("create_time").Where("uid = ?", uid).Count(&total).Limit(p.PageSize).Offset((p.PageNum - 1) * p.PageSize).Find(&videos).Error
	if err != nil {
		return resp.Pager{}, err
	}
	pager := resp.Pager{}
	pager.List = videos
	pager.Total = total
	return pager, nil
}

func (v *VideoDao) GetVideoById(id int) {
	fmt.Println("...")
}
