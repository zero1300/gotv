package video

import (
	"fmt"

	"gorm.io/gorm"
)

type VideoDao struct {
	db *gorm.DB
}

func NewVideoDao(db *gorm.DB) *VideoDao {
	return &VideoDao{
		db: db,
	}
}

func (v *VideoDao) GetVideoById(id int) {
	fmt.Println("...")
}
