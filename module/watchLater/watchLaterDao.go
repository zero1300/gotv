package watchLater

import (
	"github.com/bwmarrin/snowflake"
	"gorm.io/gorm"
	"gotv/model"
	"gotv/model/vo"
	"gotv/resp"
)

type watchLaterDao struct {
	db *gorm.DB
}

func NewWatchLater(db *gorm.DB) *watchLaterDao {
	return &watchLaterDao{db: db}
}

func (w watchLaterDao) addWL(wl model.WatchLater) error {
	if w.getWL(wl) == 1 {
		return nil
	}
	// wl: watchLater
	node, _ := snowflake.NewNode(1)
	id := node.Generate().Int64()
	wl.ID = uint(id)
	return w.db.Create(&wl).Error
}

func (w watchLaterDao) getWL(wl model.WatchLater) int64 {
	var count int64
	w.db.Model(model.WatchLater{}).Where("uid = ? and vid = ?", wl.UID, wl.VID).Count(&count)
	return count
}

func (w watchLaterDao) delWL(wl model.WatchLater) error {
	if w.getWL(wl) == 0 {
		return nil
	}
	return w.db.Where("uid = ? and vid = ?", wl.UID, wl.VID).Delete(&wl).Error
}

func (w watchLaterDao) wlList(uid uint, p model.Page) (resp.Pager, error) {
	wls := make([]model.WatchLater, 0)
	w.db.Model(model.WatchLater{}).Where("uid = ?", uid).Limit(p.PageSize).Offset((p.PageNum - 1) * p.PageSize).Find(&wls)
	videoVos := make([]vo.VideoVo, 10)
	for i := 0; i < len(wls); i++ {
		wl := wls[i]
		var videoVo vo.VideoVo
		w.db.Model(model.Video{}).Where("id = ?", wl.VID).Find(&videoVo)
		var user model.User
		user.ID = videoVo.UID
		w.db.Model(model.User{}).First(&user)
		videoVo.Nickname = user.Nickname
		var count int64
		w.db.Model(model.Comment{}).Where("vid = ?", videoVo.ID).Count(&count)
		videoVo.Comments = count
		videoVo.CreateTimeString = videoVo.CreateTime.Format("2006-01-02 15:04")
		videoVos = append(videoVos, videoVo)
	}
	pager := resp.Pager{}
	pager.List = videoVos
	pager.Total = int64(len(wls))
	return pager, nil
}
