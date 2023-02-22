package model

type DislikeRecord struct {
	UID uint
	VID uint64 `gorm:"column:vid"`
}

func (DislikeRecord) TableName() string {
	return "dislike_record"
}
