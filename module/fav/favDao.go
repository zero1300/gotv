package fav

import (
	"fmt"
	"gotv/model"
	"gotv/model/vo"
	"gotv/resp"

	"gorm.io/gorm"
)

type FavDao struct {
	db *gorm.DB
}

func NewFavDao(db *gorm.DB) *FavDao {
	return &FavDao{db: db}
}

func (f *FavDao) addFav(fav model.Fav) error {
	if f.getFav(fav) == 1 {
		return nil
	}
	return f.db.Create(&fav).Error
}

func (f *FavDao) getFav(fav model.Fav) int64 {
	var count int64
	f.db.Debug().Model(model.Fav{}).Where("uid = ? and vid = ?", fav.UID, fav.VID).Count(&count)
	fmt.Println(count)
	return count
}

func (f *FavDao) DelFav(fav model.Fav) error {
	if f.getFav(fav) == 0 {
		return nil
	}
	return f.db.Delete(fav, "uid = ? and vid = ?", fav.UID, fav.VID).Error
}

func (f *FavDao) DelFavByVid(vid string) error {
	var fav model.Fav
	return f.db.Delete(fav, "vid = ?", vid).Error
}

func (f *FavDao) favList(p model.Page, uid uint) (resp.Pager, error) {
	favs := make([]model.Fav, 10)
	var total int64
	err := f.db.Debug().Model(model.Fav{}).Order("create_time desc").Count(&total).Limit(p.PageSize).Offset((p.PageNum - 1) * p.PageSize).Find(&favs).Error
	if err != nil {
		return resp.Pager{}, err
	}
	videoVos := make([]vo.VideoVo, 0)
	for i := 0; i < len(favs); i++ {
		fav := favs[i]
		var videoVo vo.VideoVo
		f.db.Debug().Model(model.Video{}).Where("id = ?", fav.VID).Find(&videoVo)
		fmt.Println(videoVo)
		var user model.User
		user.ID = videoVo.UID
		f.db.Debug().Model(model.User{}).First(&user)
		fmt.Println(user)
		videoVo.Nickname = user.Nickname
		var c int64
		f.db.Model(model.Comment{}).Where("vid = ?", videoVo.ID).Count(&c)
		fmt.Println(c)
		videoVo.Comments = c
		videoVo.CreateTimeString = videoVo.CreateTime.Format("2006-01-02 15:04")
		videoVos = append(videoVos, videoVo)
	}
	pager := resp.Pager{}
	pager.List = videoVos
	pager.Total = total
	return pager, nil
}
