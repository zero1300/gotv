package video

import (
	"fmt"
	"gotv/model"
	"gotv/model/vo"
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

func (v *VideoDao) latestVideo(p model.Page) (resp.Pager, error) {
	videoVos := make([]vo.VideoVo, 10)
	var total int64
	err := v.db.Debug().Model(model.Video{}).Order("create_time desc").Count(&total).Limit(p.PageSize).Offset((p.PageNum - 1) * p.PageSize).Find(&videoVos).Error
	if err != nil {
		return resp.Pager{}, err
	}
	for i := 0; i < len(videoVos); i++ {
		var user model.User
		user.ID = videoVos[i].UID
		v.db.Model(model.User{}).First(&user)
		videoVos[i].Nickname = user.Nickname
		fmt.Println(videoVos[i])
	}

	pager := resp.Pager{}
	pager.List = videoVos
	pager.Total = total
	return pager, nil
}

func (v *VideoDao) GetVideoById(id int) {
	fmt.Println("...")
}
