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

// 获得视频下的评论
func (c commentController) commentList(ctx *gin.Context) {
	vid := ctx.Param("vid")

	var p model.Page
	if err := ctx.ShouldBind(&p); err != nil {
		resp.Fail(ctx, "参数异常: "+err.Error())
		return
	}

	data, err := c.commentDao.commentList(vid, p)
	if err != nil {
		resp.Fail(ctx, err.Error())
		return
	}
	resp.Success(ctx, data)

}

// ----- admin -----
func (c commentController) adminCommentList(ctx *gin.Context) {
	var p model.Page
	if err := ctx.ShouldBind(&p); err != nil {
		resp.Fail(ctx, "参数异常: "+err.Error())
		return
	}
	data, err := c.commentDao.adminCommentList(p)
	if err != nil {
		resp.Fail(ctx, err.Error())
		return
	}
	resp.Success(ctx, data)
}

func (c commentController) delComment(ctx *gin.Context) {
	id := ctx.Param("id")
	c.commentDao.delComment(id)
}

func (c commentController) searchComment(ctx *gin.Context) {
	var p model.Page
	if err := ctx.ShouldBind(&p); err != nil || p.Keyword == "" {
		resp.Fail(ctx, "参数异常: "+err.Error())
		return
	}
	comment, err := c.commentDao.searchComment(p)
	if err != nil {
		resp.Fail(ctx, err.Error())
		return
	}
	resp.Success(ctx, comment)
}

func (c commentController) SetUp(admin *gin.RouterGroup, api *gin.RouterGroup) {
	comment := api.Group("comment")
	comment.POST("/", c.comment)
	comment.POST("/list/:vid", c.commentList)

	adminComment := admin.Group("comment")
	adminComment.POST("/list", c.adminCommentList)
	adminComment.POST("/search", c.searchComment)
	adminComment.DELETE("/:id", c.delComment)
}
