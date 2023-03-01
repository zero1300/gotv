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

func (u *UserDao) DelUser(id string) error {
	var users []model.User
	tx := u.db.Delete(users, id)
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
		u.db.Debug().Model(model.User{}).Order("id").Where("nickname LIKE ?", "%"+p.Keyword+"%").Count(&total).Limit(p.PageSize).Offset((p.PageNum - 1) * p.PageSize).Find(&users)
	} else {
		u.db.Debug().Model(model.User{}).Order("id").Count(&total).Limit(p.PageSize).Offset((p.PageNum - 1) * p.PageSize).Find(&users)
	}
	err := u.db.Error
	pager.List = users
	pager.Total = total
	if err != nil {
		return resp.Pager{}, err
	}
	return pager, nil
}

func (u *UserDao) GetUser(id string) *model.User {
	var user model.User
	u.db.Where("id = ?", id).First(&user)
	return &user
}

func (u *UserDao) updateUser(user model.User) {
	user.Phone = ""
	u.db.Model(&user).Updates(user)
}

// ----- 统计接口 -----
// 统计用户动态数量
func (u UserDao) countUserDynamic(uid uint) int64 {
	var count int64
	u.db.Model(model.Video{}).Where("uid = ?", uid).Count(&count)
	return count
}

// 统计用户关注数
func (u UserDao) countSub(uid uint) int64 {
	var count int64
	u.db.Model(model.Sub{}).Where("fans = ?", uid).Count(&count)
	return count
}

// 统计用户粉丝数量
func (u UserDao) countFans(uid uint) int64 {
	var count int64
	u.db.Model(model.Sub{}).Where("uid = ?", uid).Count(&count)
	return count
}
