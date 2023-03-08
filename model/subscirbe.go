package model

type Sub struct {
	ID   uint   `gorm:"primaryKey" json:"id,string"`
	UID  uint64 `json:"uid"`
	Fans uint   `json:"-"`
}

func (Sub) TableName() string {
	return "sub"
}
