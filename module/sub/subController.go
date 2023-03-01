package sub

import (
	"gotv/model"
	"gotv/resp"

	"github.com/gin-gonic/gin"
)

type subController struct {
	subDao *subDao
}

func NewSubController(dao *subDao) *subController {
	return &subController{
		subDao: dao,
	}
}

func (s subController) getSubVo(ctx *gin.Context) model.Sub {
	var sub model.Sub
	obj, _ := ctx.Get("user")
	user := obj.(*model.User)
	uid := ctx.Param("uid")
	sub.Fans = user.ID
	sub.UID = uid
	return sub
}

func (s subController) addSub(ctx *gin.Context) {
	subVo := s.getSubVo(ctx)
	_ = s.subDao.addSub(subVo)
}

func (s subController) getSub(ctx *gin.Context) {
	subVo := s.getSubVo(ctx)
	count := s.subDao.getSub(subVo)
	resp.Success(ctx, count)
}
func (s subController) delSub(ctx *gin.Context) {
	subVo := s.getSubVo(ctx)
	_ = s.subDao.delSub(subVo)
}

func (s subController) subList(ctx *gin.Context) {
	obj, _ := ctx.Get("user")
	user := obj.(*model.User)
	var p model.Page
	if err := ctx.ShouldBind(&p); err != nil {
		resp.Fail(ctx, "参数异常: "+err.Error())
		return
	}
	list := s.subDao.subList(user.ID, p)
	resp.Success(ctx, list)
}

func (s subController) fansList(ctx *gin.Context) {
	obj, _ := ctx.Get("user")
	user := obj.(*model.User)
	var p model.Page
	if err := ctx.ShouldBind(&p); err != nil {
		resp.Fail(ctx, "参数异常: "+err.Error())
		return
	}
	list := s.subDao.fansList(user.ID, p)
	resp.Success(ctx, list)
}

func (s subController) SetUp(admin *gin.RouterGroup, api *gin.RouterGroup) {
	sub := api.Group("sub")
	sub.GET("/addSub/:uid", s.addSub)
	sub.GET("/getSub/:uid", s.getSub)
	sub.GET("/delSub/:uid", s.delSub)
	sub.POST("/subList", s.subList)
	sub.POST("fansList", s.fansList)
}
