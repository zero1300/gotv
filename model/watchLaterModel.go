package model

type WatchLater struct {
	ID  uint   `gorm:"primaryKey" json:"id,string"`
	UID uint   `json:"uid,string"`
	VID uint64 `gorm:"column:vid" json:"vid"`
}

func (WatchLater) TableName() string {
	return "watch_later"
}
