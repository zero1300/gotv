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
	commentsDetailVos := make([]vo.CommentDetailVo, 10)

	var total int64
	err := c.db.Debug().Model(model.Comment{}).Count(&total).Where("vid = ?", vid).Order("create_time desc").Limit(p.PageSize).Offset((p.PageNum - 1) * p.PageSize).Find(&commentsDetailVos).Error
	if err != nil {
		return resp.Pager{}, err
	}
	for i := 0; i < len(commentsDetailVos); i++ {
		uid := commentsDetailVos[i].UID
		var user model.User
		user.ID = uid
		c.db.Model(model.User{}).First(&user)
		commentsDetailVos[i].Avatar = user.Avatar
		commentsDetailVos[i].Nickname = user.Nickname
		commentsDetailVos[i].CreateTimeString = commentsDetailVos[i].CreateTime.Format("2006-01-02 15:04")
	}
	pager := resp.Pager{}
	pager.List = commentsDetailVos
	pager.Total = total
	return pager, nil
}

// ----- admin -----
func (c commentDao) adminCommentList(p model.Page) (resp.Pager, error) {
	commentsDetailVos := make([]vo.CommentDetailVo, 10)
	err := c.db.Debug().Model(model.Comment{}).Order("create_time desc").Limit(p.PageSize).Offset((p.PageNum - 1) * p.PageSize).Find(&commentsDetailVos).Error
	if err != nil {
		return resp.Pager{}, err
	}
	for i := 0; i < len(commentsDetailVos); i++ {
		uid := commentsDetailVos[i].UID
		var user model.User
		user.ID = uid
		c.db.Model(model.User{}).First(&user)
		var video model.Video
		vid := commentsDetailVos[i].VID
		video.ID = vid
		c.db.Model(model.Video{}).First(&video)
		commentsDetailVos[i].Title = video.Title
		commentsDetailVos[i].Avatar = user.Avatar
		commentsDetailVos[i].Nickname = user.Nickname
		commentsDetailVos[i].CreateTimeString = commentsDetailVos[i].CreateTime.Format("2006-01-02 15:04")
	}
	pager := resp.Pager{}
	pager.List = commentsDetailVos
	pager.Total = int64(len(commentsDetailVos))
	return pager, nil
}

func (c commentDao) delComment(id string) {
	var comment model.Comment
	c.db.Where("id = ?", id).Delete(&comment)
}

func (c commentDao) DelComment(vid string) error {
	var comment model.Comment
	return c.db.Delete(comment, "vid = ? ", vid).Error
}

func (c commentDao) searchComment(p model.Page) (resp.Pager, error) {
	commentsDetailVos := make([]vo.CommentDetailVo, 10)
	err := c.db.Debug().Model(model.Comment{}).Where("content like ?", "%"+p.Keyword+"%").Order("create_time desc").Limit(p.PageSize).Offset((p.PageNum - 1) * p.PageSize).Find(&commentsDetailVos).Error
	if err != nil {
		return resp.Pager{}, err
	}
	for i := 0; i < len(commentsDetailVos); i++ {
		uid := commentsDetailVos[i].UID
		var user model.User
		user.ID = uid
		c.db.Model(model.User{}).First(&user)
		var video model.Video
		vid := commentsDetailVos[i].VID
		video.ID = vid
		c.db.Model(model.Video{}).First(&video)
		commentsDetailVos[i].Title = video.Title
		commentsDetailVos[i].Avatar = user.Avatar
		commentsDetailVos[i].Nickname = user.Nickname
		commentsDetailVos[i].CreateTimeString = commentsDetailVos[i].CreateTime.Format("2006-01-02 15:04")
	}
	pager := resp.Pager{}
	pager.List = commentsDetailVos
	pager.Total = int64(len(commentsDetailVos))
	return pager, nil
}
