package video

import (
	"fmt"
	"gotv/model"
	"gotv/model/vo"
	"gotv/resp"
	"strconv"

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
		var count int64
		v.db.Model(model.Comment{}).Where("vid = ?", videoVos[i].ID).Count(&count)
		videoVos[i].Comments = count
		videoVos[i].CreateTimeString = videoVos[i].CreateTime.Format("2006-01-02 15:04")
	}

	pager := resp.Pager{}
	pager.List = videoVos
	pager.Total = total
	return pager, nil
}

func (v *VideoDao) addViews(vid string) {
	v.db.Debug().Model(model.Video{}).Where("id = ?", vid).UpdateColumn("views", gorm.Expr("views + ?", 1))
}

func (v *VideoDao) addLike(vid string, uid uint) {
	fmt.Println(uid)
	if v.getLikeRecord(vid, uid) == 1 {
		return
	}
	v.db.Debug().Model(model.Video{}).Where("id = ?", vid).UpdateColumn("like", gorm.Expr(`"like" + ?`, 1))
	var like_record model.LikeRecordModel
	node, _ := snowflake.NewNode(1)
	id := node.Generate().Int64()
	like_record.ID = uint(id)
	like_record.UID = uid
	u64, _ := strconv.ParseUint(vid, 10, 64)
	like_record.VID = u64
	v.db.Debug().Model(model.LikeRecordModel{}).Create(&like_record)
}

func (v *VideoDao) cancelLike(vid string, uid uint) {
	if v.getLikeRecord(vid, uid) == 0 {
		return
	}
	v.db.Debug().Model(model.Video{}).Where("id = ?", vid).UpdateColumn("like", gorm.Expr(`"like" - ?`, 1))
	v.db.Debug().Model(model.LikeRecordModel{}).Where("vid = ? and uid = ?", vid, uid).Delete(&model.LikeRecordModel{})
}

func (v *VideoDao) addDisLike(vid string, uid uint) {
	if v.getDislikeRecord(vid, uid) == 1 {
		return
	}
	var dislike_record model.DislikeRecord
	u64, _ := strconv.ParseUint(vid, 10, 64)
	dislike_record.VID = u64
	dislike_record.UID = uid
	v.db.Debug().Model(model.DislikeRecord{}).Create(&dislike_record)

}

func (v *VideoDao) cancelDislike(vid string, uid uint) {
	if v.getDislikeRecord(vid, uid) == 0 {
		return
	}
	v.db.Debug().Model(model.DislikeRecord{}).Where("vid = ? and uid = ?", vid, uid).Delete(&model.LikeRecordModel{})
}
func (v *VideoDao) getLikeRecord(vid string, uid uint) int64 {
	var count int64
	v.db.Debug().Model(model.LikeRecordModel{}).Where("vid = ? and uid = ?", vid, uid).Count(&count)
	return count
}

func (v *VideoDao) getDislikeRecord(vid string, uid uint) int64 {
	var count int64
	v.db.Debug().Model(model.DislikeRecord{}).Where("vid = ? and uid = ?", vid, uid).Count(&count)
	return count
}
