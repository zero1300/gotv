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

func (v *VideoDao) hot(p model.Page) (resp.Pager, error) {
	videoVos := make([]vo.VideoVo, 10)
	var total int64
	err := v.db.Debug().Model(model.Video{}).Order("views desc").Count(&total).Limit(p.PageSize).Offset((p.PageNum - 1) * p.PageSize).Find(&videoVos).Error
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
func (v VideoDao) countVideoLike(vid string) int64 {
	var count int64
	v.db.Debug().Model(model.LikeRecordModel{}).Where("vid = ?", vid).Count(&count)
	return count
}

func (v VideoDao) countVideoDisLike(vid string) int64 {
	var count int64
	v.db.Debug().Model(model.DislikeRecord{}).Where("vid = ?", vid).Count(&count)
	return count
}

// 动态接口
func (v *VideoDao) dynamic(p model.Page, uid uint) resp.Pager {

	subs := make([]model.Sub, 0)
	v.db.Debug().Where("fans = ?", uid).Find(&subs)
	uids := make([]uint64, 0)
	for _, sub := range subs {
		uids = append(uids, sub.UID)
	}

	var count int64
	videos := make([]model.Video, 0)
	v.db.Debug().Model(model.Video{}).Where("uid in ?", uids).Order("create_time desc").Count(&count).Limit(p.PageSize).Offset((p.PageNum - 1) * p.PageSize).Find(&videos)

	dynamics := make([]vo.Dynamic, 0)
	for _, video := range videos {
		var user model.User
		v.db.Debug().Model(model.User{}).Where("id = ?", video.UID).Find(&user)
		var dynamic vo.Dynamic
		dynamic.Video = video
		dynamic.User = user
		dynamics = append(dynamics, dynamic)
		dynamic.Video.CreateTimeString = dynamic.Video.CreateTime.Format("2006-01-02 15:04")
	}

	for i := 0; i < len(dynamics); i++ {
		dynamics[i].Video.CreateTimeString = dynamics[i].Video.CreateTime.Format("2006-01-02 15:04")
	}

	pager := resp.Pager{}
	pager.List = dynamics
	pager.Total = count
	return pager

}

func (v VideoDao) GetVideoById(vid string) model.Video {
	var video model.Video
	v.db.Where("id = ?", vid).Find(&video)
	return video
}

// 删除视频
func (v VideoDao) delVideo(id string) error {

	var video model.Video

	var f model.Fav
	v.db.Delete(f, "vid = ?", id)

	var h model.History
	v.db.Delete(h, "vid = ?", id)

	return v.db.Delete(video, "id = ?", id).Error
}

// 删除 like
func (v VideoDao) delLike(vid string) error {
	var like model.LikeRecordModel
	return v.db.Delete(like, "vid = ?", vid).Error
}

// 删除dislike
func (v VideoDao) delDislike(vid string) error {
	var disLike model.DislikeRecord
	return v.db.Delete(disLike, "vid = ?", vid).Error
}

// 视频搜索
func (v *VideoDao) searchVideo(p model.Page) (resp.Pager, error) {
	videos := make([]vo.VideoVo, 0)
	v.db.Debug().Model(model.Video{}).Where("title LIKE ? ", "%"+p.Keyword+"%").Limit(p.PageSize).Offset((p.PageNum - 1) * p.PageSize).Find(&videos)

	for i := 0; i < len(videos); i++ {
		var user model.User
		user.ID = videos[i].UID
		v.db.Model(model.User{}).First(&user)
		videos[i].Nickname = user.Nickname
		var count int64
		v.db.Model(model.Comment{}).Where("vid = ?", videos[i].ID).Count(&count)
		videos[i].Comments = count
		videos[i].CreateTimeString = videos[i].CreateTime.Format("2006-01-02 15:04")
	}
	pager := resp.Pager{}
	pager.List = videos
	pager.Total = int64(len(videos))
	return pager, nil
}

// ----- admin -----
func (v VideoDao) videoAdminList(p model.Page) (resp.Pager, error) {
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

func (v VideoDao) getVideoById(id string) model.Video {
	var video model.Video
	v.db.Model(model.Video{}).Where("id = ?", id).Find(&video)
	return video
}

func (v VideoDao) updateVideoInfo(video model.Video) {
	v.db.Model(&video).Updates(video)
}
