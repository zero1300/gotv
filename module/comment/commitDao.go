package comment

import (
	"fmt"
	"gotv/model"

	"github.com/bwmarrin/snowflake"
	"gorm.io/gorm"
)

type commentDao struct {
	db *gorm.DB
}

func NewCommentDao(db *gorm.DB) *commentDao {
	return &commentDao{
		db: db,
	}
}

func (c commentDao) addComment(comment model.Comment) error {
	node, _ := snowflake.NewNode(1)
	id := node.Generate().Int64()
	comment.ID = uint(id)
	fmt.Println(comment)
	return c.db.Debug().Create(&comment).Error
}
