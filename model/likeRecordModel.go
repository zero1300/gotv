package model

type LikeRecordModel struct {
	ID  uint
	UID uint
	VID uint64 `gorm:"column:vid"`
}

func (LikeRecordModel) TableName() string {
	return "like_record"
}
