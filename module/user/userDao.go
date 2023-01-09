package user

import (
	"gotv/model"

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

func (u *UserDao) userList(p model.Page) ([]model.User, error) {

	users := make([]model.User, 10)

	if p.Keyword != "" {
		u.db.Debug().Limit(p.PageSize).Offset((p.PageNum-1)*p.PageSize).Where("nickname LIKE ?", "%"+p.Keyword+"%").Find(&users)
	} else {
		u.db.Debug().Limit(p.PageSize).Offset((p.PageNum - 1) * p.PageSize).Find(&users)
	}
	err := u.db.Error
	if err != nil {
		return nil, err
	}
	return users, nil

}
