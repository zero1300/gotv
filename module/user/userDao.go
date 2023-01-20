package user

import (
	"gotv/model"
	"gotv/resp"

	"gorm.io/gorm"
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}

func (u *UserDao) AddUser(user model.User) error {
	return u.db.Create(&user).Error

}

func (u *UserDao) DelUser(ids []int) error {
	var users []model.User
	tx := u.db.Delete(users, ids)
	return tx.Error
}

func (u *UserDao) QueryByPhone(phone string) *model.User {
	var user model.User
	u.db.Where("phone = ?", phone).First(&user)
	return &user
}

func (u *UserDao) userList(p model.Page) (resp.Pager, error) {

	users := make([]model.User, 10)

	pager := resp.Pager{}

	var total int64

	if p.Keyword != "" {
		u.db.Debug().Model(model.User{}).Where("nickname LIKE ?", "%"+p.Keyword+"%").Count(&total).Limit(p.PageSize).Offset((p.PageNum - 1) * p.PageSize).Find(&users)
	} else {
		u.db.Debug().Model(model.User{}).Count(&total).Limit(p.PageSize).Offset((p.PageNum - 1) * p.PageSize).Find(&users)
	}
	err := u.db.Error
	pager.List = users
	pager.Total = total
	if err != nil {
		return resp.Pager{}, err
	}
	return pager, nil

}
