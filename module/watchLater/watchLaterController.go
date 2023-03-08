package watchLater

import (
	"github.com/gin-gonic/gin"
	"gotv/common"
	"gotv/model"
	"gotv/module/user"
	"gotv/resp"
	"strconv"
)

type watchLaterController struct {
	watchLaterDao *watchLaterDao
}

func NewWatchLaterController(dao *watchLaterDao) *watchLaterController {
	return &watchLaterController{watchLaterDao: dao}
}

func (w watchLaterController) addWL(ctx *gin.Context) {
	var wl model.WatchLater
	uid := user.GetUseIdrByCtx(ctx)
	vid := ctx.Param("vid")
	wl.UID = uid
	wl.VID, _ = strconv.ParseUint(vid, 10, 64)
	_ = w.watchLaterDao.addWL(wl)
}

func (w watchLaterController) delWL(ctx *gin.Context) {
	var wl model.WatchLater
	uid := user.GetUseIdrByCtx(ctx)
	vid := ctx.Param("vid")
	wl.UID = uid
	wl.VID, _ = strconv.ParseUint(vid, 10, 64)
	_ = w.watchLaterDao.delWL(wl)
}

func (w watchLaterController) wlList(ctx *gin.Context) {

	page, err := common.GetPage(ctx)
	if err != nil {
		resp.Fail(ctx, err.Error())
	}
	uid := user.GetUseIdrByCtx(ctx)

	list, err := w.watchLaterDao.wlList(uid, page)
	if err != nil {
		return
	}
	resp.Success(ctx, list)
}

func (w watchLaterController) Setup(admin *gin.RouterGroup, api *gin.RouterGroup) {
	wl := api.Group("wl")
	wl.GET("/addWL/:vid", w.addWL)
	wl.GET("/delWL/:vid", w.delWL)
	wl.POST("/wlList", w.wlList)
}
