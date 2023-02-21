package comment

import (
	"fmt"
	"gotv/model"
	"gotv/model/vo"
	"gotv/resp"

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

func (c commentDao) commentList(vid string, p model.Page) (resp.Pager, error) {
	comemntsDetailVos := make([]vo.CommentDetailVo, 10)

	var total int64
	err := c.db.Debug().Model(model.Comment{}).Count(&total).Where("vid = ?", vid).Order("create_time desc").Limit(p.PageSize).Offset((p.PageNum - 1) * p.PageSize).Find(&comemntsDetailVos).Error
	if err != nil {
		return resp.Pager{}, err
	}
	for i := 0; i < len(comemntsDetailVos); i++ {
		uid := comemntsDetailVos[i].UID
		var user model.User
		user.ID = uid
		c.db.Model(model.User{}).First(&user)
		comemntsDetailVos[i].Avatar = user.Avatar
		comemntsDetailVos[i].Nickname = user.Nickname
	}
	pager := resp.Pager{}
	pager.List = comemntsDetailVos
	pager.Total = total
	return pager, nil
}
