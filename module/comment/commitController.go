package comment

import (
	"fmt"
	"gotv/model"
	"gotv/resp"

	"github.com/gin-gonic/gin"
)

type commentController struct {
	commentDao *commentDao
}

// 构造函数
func NewCommentController(commentDao *commentDao) *commentController {
	return &commentController{
		commentDao: commentDao,
	}
}

// 用户评论视频
func (c commentController) comment(ctx *gin.Context) {
	obj, _ := ctx.Get("user")
	user := obj.(*model.User)
	var comment model.Comment
	comment.UID = user.ID
	err := ctx.ShouldBind(&comment)
	if err != nil {
		fmt.Println(err)
	}
	err = c.commentDao.addComment(comment)
	if err != nil {
		fmt.Println(err)
	}
	resp.Success(ctx, comment)
}

func (c commentController) SetUp(admin *gin.RouterGroup, api *gin.RouterGroup) {
	comment := api.Group("comment")
	comment.POST("/", c.comment)
}
