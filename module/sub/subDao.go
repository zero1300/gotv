package sub

import (
	"github.com/bwmarrin/snowflake"
	"gorm.io/gorm"
	"gotv/model"
	"gotv/resp"
)

type subDao struct {
	db *gorm.DB
}

func NewSubDao(db *gorm.DB) *subDao {
	return &subDao{
		db: db,
	}
}

func (s subDao) addSub(sub model.Sub) error {
	if s.getSub(sub) == 1 {
		return nil
	}
	node, _ := snowflake.NewNode(1)
	id := node.Generate().Int64()
	sub.ID = uint(id)
	return s.db.Create(&sub).Error
}
func (s subDao) delSub(sub model.Sub) error {
	if s.getSub(sub) == 0 {
		return nil
	}
	return s.db.Where("uid = ? and fans = ?", sub.UID, sub.Fans).Delete(&sub).Error
}

func (s subDao) getSub(sub model.Sub) int64 {
	var count int64
	s.db.Model(&model.Sub{}).Where("uid = ? and fans = ?", sub.UID, sub.Fans).Count(&count)
	return count
}

func (s subDao) subList(fans uint, p model.Page) resp.Pager {
	subs := make([]model.Sub, 0)
	var count int64
	s.db.Debug().Where("fans = ?", fans).Limit(p.PageSize).Offset((p.PageNum - 1) * p.PageSize).Find(&subs).Count(&count)
	users := make([]model.User, 0)
	for _, sub := range subs {
		s.db.Debug().Model(model.User{}).Where("id = ?", sub.UID).Find(&users)
	}
	pager := resp.Pager{}
	pager.List = users
	pager.Total = count
	return pager
}
