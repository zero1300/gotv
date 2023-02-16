package file

import (
	"gotv/model"

	"github.com/bwmarrin/snowflake"
	"gorm.io/gorm"
)

type fileDao struct {
	db *gorm.DB
}

func NewFileDao(db *gorm.DB) *fileDao {
	return &fileDao{
		db: db,
	}
}

func (f fileDao) SaveFile(file model.File) error {
	node, _ := snowflake.NewNode(1)
	id := node.Generate().Int64()
	file.ID = uint(id)
	return f.db.Create(&file).Error
}
