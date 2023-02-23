package fav

import (
	"gotv/model"
	"gotv/resp"
	"strconv"

	"github.com/gin-gonic/gin"
)

type favController struct {
	favDao *favDao
}

func NewFavController(favDao *favDao) *favController {
	return &favController{favDao: favDao}
}

func (f favController) addFav(ctx *gin.Context) {
	vid := ctx.Param("vid")
	obj, _ := ctx.Get("user")
	user := obj.(*model.User)
	var fav model.Fav
	fav.UID = user.ID
	fav.VID, _ = strconv.ParseUint(vid, 10, 64)
	f.favDao.addFav(fav)
}

func (f favController) delFav(ctx *gin.Context) {
	vid := ctx.Param("vid")
	obj, _ := ctx.Get("user")
	user := obj.(*model.User)
	var fav model.Fav
	fav.UID = user.ID
	fav.VID, _ = strconv.ParseUint(vid, 10, 64)
	f.favDao.delFav(fav)
}

func (f favController) getFav(ctx *gin.Context) {
	vid := ctx.Param("vid")
	obj, _ := ctx.Get("user")
	user := obj.(*model.User)
	var fav model.Fav
	fav.UID = user.ID
	fav.VID, _ = strconv.ParseUint(vid, 10, 64)
	r := f.favDao.getFav(fav)
	resp.Success(ctx, r)
}

func (f favController) favList(ctx *gin.Context) {
	obj, _ := ctx.Get("user")
	user := obj.(*model.User)
	var p model.Page
	if err := ctx.ShouldBind(&p); err != nil {
		resp.Fail(ctx, "参数异常: "+err.Error())
		return
	}
	pager, _ := f.favDao.favList(p, user.ID)
	resp.Success(ctx, pager)

}

func (f favController) SetUp(admin *gin.RouterGroup, api *gin.RouterGroup) {
	fav := api.Group("fav")
	fav.GET("/addFav/:vid", f.addFav)
	fav.GET("/delFav/:vid", f.delFav)
	fav.GET("/getFav/:vid", f.getFav)
	fav.POST("/favList", f.favList)

}
