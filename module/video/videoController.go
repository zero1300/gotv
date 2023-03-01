package video

import (
	"context"
	"fmt"
	"gotv/common"
	"gotv/model"
	"gotv/resp"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
)

type VideoHandler struct {
	videoDao *VideoDao
	context  context.Context
	rc       *redis.Client
}

func NewVideoHandler(videoDao *VideoDao, rc *redis.Client) *VideoHandler {
	return &VideoHandler{videoDao: videoDao, context: context.Background(), rc: rc}
}

func (v *VideoHandler) addVideo(ctx *gin.Context) {
	obj, _ := ctx.Get("user")
	user := obj.(*model.User)
	var video model.Video
	ctx.ShouldBind(&video)
	url := common.OSS_BASE_URL + video.Uri
	fmt.Println(url)
	duration, err := duration(url)
	if err != nil {

		resp.Fail(ctx, "获取视频时长失败,请稍后再试...")
		return
	}
	// TODO fix

	video.Duration = strconv.Itoa(int(duration))
	video.UID = user.ID
	v.videoDao.addVideo(video)
	resp.Success(ctx, video)
}

func (v *VideoHandler) getVideoListByUid(ctx *gin.Context) {
	var p model.Page
	if err := ctx.ShouldBind(&p); err != nil {
		resp.Fail(ctx, "参数异常: "+err.Error())
		return
	}
	obj, _ := ctx.Get("user")
	user := obj.(*model.User)
	pager, err := v.videoDao.getVideoListByUid(user.ID, p)
	if err != nil {
		resp.Fail(ctx, "查询失败："+err.Error())
		return
	}
	resp.Success(ctx, pager)
}

func (v *VideoHandler) videoRecommend(ctx *gin.Context) {
	var p model.Page
	if err := ctx.ShouldBind(&p); err != nil {
		resp.Fail(ctx, "参数异常: "+err.Error())
		return
	}
	pager, err := v.videoDao.latestVideo(p)
	if err != nil {
		resp.Fail(ctx, "查询失败: "+err.Error())
		return
	}
	resp.Success(ctx, pager)
}

// 热门视频
func (v VideoHandler) hot(ctx *gin.Context) {
	var p model.Page
	if err := ctx.ShouldBind(&p); err != nil {
		resp.Fail(ctx, "参数异常: "+err.Error())
		return
	}
	hot, err := v.videoDao.hot(p)
	if err != nil {
		resp.Fail(ctx, "查询失败: "+err.Error())
		return
	}
	resp.Success(ctx, hot)

}

func (v VideoHandler) dynamic(ctx *gin.Context) {
	var p model.Page
	if err := ctx.ShouldBind(&p); err != nil {
		resp.Fail(ctx, "参数异常: "+err.Error())
		return
	}
	obj, _ := ctx.Get("user")
	user := obj.(*model.User)
	dynamics := v.videoDao.dynamic(p, user.ID)
	resp.Success(ctx, dynamics)
}

func (v *VideoHandler) getOne(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		return
	}
	video := v.videoDao.getVideoById(id)
	resp.Success(ctx, video)

}

func (v *VideoHandler) addViews(ctx *gin.Context) {

	vid := ctx.Param("vid")
	fmt.Println(vid)
	v.videoDao.addViews(vid)
}

// 点赞视频
func (v *VideoHandler) addLike(ctx *gin.Context) {
	obj, _ := ctx.Get("user")
	user := obj.(*model.User)
	id := ctx.Param("vid")
	v.videoDao.addLike(id, user.ID)
}

// 取消点赞视频
func (v *VideoHandler) cancelLike(ctx *gin.Context) {
	obj, _ := ctx.Get("user")
	user := obj.(*model.User)
	id := ctx.Param("vid")
	fmt.Println(id)
	v.videoDao.cancelLike(id, user.ID)
}

// 点踩视频
func (v *VideoHandler) addDisLike(ctx *gin.Context) {
	obj, _ := ctx.Get("user")
	user := obj.(*model.User)
	id := ctx.Param("vid")
	v.videoDao.addDisLike(id, user.ID)
}

// 取消点踩视频
func (v *VideoHandler) cancelDislike(ctx *gin.Context) {
	obj, _ := ctx.Get("user")
	user := obj.(*model.User)
	id := ctx.Param("vid")
	v.videoDao.cancelDislike(id, user.ID)
}

func (v *VideoHandler) getLikeRecord(ctx *gin.Context) {
	obj, _ := ctx.Get("user")
	user := obj.(*model.User)
	id := ctx.Param("vid")
	b := v.videoDao.getLikeRecord(id, user.ID)
	resp.Success(ctx, b)
}

func (v *VideoHandler) getDislikeRecord(ctx *gin.Context) {
	obj, _ := ctx.Get("user")
	user := obj.(*model.User)
	id := ctx.Param("vid")
	b := v.videoDao.getDislikeRecord(id, user.ID)
	resp.Success(ctx, b)
}

func (v *VideoHandler) searchVideo(ctx *gin.Context) {

	var page model.Page
	err := ctx.ShouldBind(&page)
	if err != nil {
		resp.Fail(ctx, err.Error())
		return
	}
	videos, _ := v.videoDao.searchVideo(page)
	resp.Success(ctx, videos)
}

// ----- admin -----
// 后台管理视频列表
func (v *VideoHandler) videoAdminList(ctx *gin.Context) {
	var p model.Page
	if err := ctx.ShouldBind(&p); err != nil {
		resp.Fail(ctx, "参数异常: "+err.Error())
		return
	}
	pager, err := v.videoDao.videoAdminList(p)
	if err != nil {
		resp.Fail(ctx, "查询失败: "+err.Error())
		return
	}
	resp.Success(ctx, pager)
}

func (v VideoHandler) delVideo(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		return
	}
	err := v.videoDao.delVideo(id)
	if err != nil {
		resp.Fail(ctx, err.Error())
		return
	}
	resp.Success(ctx, "ok")
}

func (v VideoHandler) changeVideoInfo(ctx *gin.Context) {
	var video model.Video
	ctx.ShouldBind(&video)
	v.videoDao.updateVideoInfo(video)
	resp.Success(ctx, video)
}

func (v VideoHandler) countVideoLike(ctx *gin.Context) {
	id := ctx.Param("vid")
	like := v.videoDao.countVideoLike(id)
	resp.Success(ctx, like)
}

func (v VideoHandler) countVideoDislike(ctx *gin.Context) {
	id := ctx.Param("vid")
	dislike := v.videoDao.countVideoDisLike(id)
	resp.Success(ctx, dislike)
}

func (v *VideoHandler) SetUp(admin *gin.RouterGroup, api *gin.RouterGroup) {
	video := api.Group("/video")
	video.GET("/getOne", v.getOne)
	video.POST("/upload", v.addVideo)
	video.POST("/listByUid", v.getVideoListByUid)
	video.POST("/recommend", v.videoRecommend)
	video.POST("/hot", v.hot)
	video.GET("/addViews/:vid", v.addViews)
	video.GET("/addLike/:vid", v.addLike)
	video.GET("/cancelLike/:vid", v.cancelLike)
	video.GET("/addDislike/:vid", v.addDisLike)
	video.GET("/cancelDislike/:vid", v.cancelDislike)
	video.GET("getLikeRecord/:vid", v.getLikeRecord)
	video.GET("/getDislikeRecord/:vid", v.getDislikeRecord)
	video.POST("/dynamic", v.dynamic)
	video.GET("countVideoLike/:vid", v.countVideoLike)
	video.GET("/countVideoDislike/:vid", v.countVideoDislike)
	video.POST("/searchVideo", v.searchVideo)
	//video.DELETE("/:id", v.delVideo)

	videoAdmin := admin.Group("/video")
	videoAdmin.POST("/list", v.videoAdminList)
	videoAdmin.DELETE("/:id", v.delVideo)
	videoAdmin.GET("/:id", v.getOne)
	videoAdmin.POST("/update", v.changeVideoInfo)
}
