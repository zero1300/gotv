package model

import "encoding/json"

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id,string"`
	Phone    string `form:"phone" binding:"required" json:"phone" grom:"phone"`
	Code     string `form:"code" binding:"required" gorm:"-" json:"-"`
	Password string `form:"password" json:"-"`
	Avator   string `json:"avator"`
	Nickname string `json:"nickname"`
}

type Tabler interface {
	TableName() string
}

// TableName 会将 User 的表名重写为 `user`
func (User) TableName() string {
	return "user"
}

func (u *User) MarshalBinary() (data []byte, err error) {
	return json.Marshal(u)
}

func (u *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}
