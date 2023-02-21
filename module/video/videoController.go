package video

import (
	"context"
	"fmt"
	"gotv/common"
	"gotv/model"
	"gotv/resp"

	log "github.com/sirupsen/logrus"

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
		log.Error(err.Error())
		resp.Fail(ctx, "获取视频时长失败,请稍后再试...")
		return
	}
	// TODO fix
	_ = duration.String()
	video.Duration = "0"
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

func (v *VideoHandler) getOne(ctx *gin.Context) {
	v.videoDao.GetVideoById(123)
}

func (v *VideoHandler) addViews(ctx *gin.Context) {
	vid := ctx.Param("vid")
	fmt.Println(vid)
	v.videoDao.addViews(vid)
}

func (v *VideoHandler) SetUp(admin *gin.RouterGroup, api *gin.RouterGroup) {
	video := api.Group("/video")
	video.GET("/getOne", v.getOne)
	video.POST("/upload", v.addVideo)
	video.POST("/listByUid", v.getVideoListByUid)
	video.POST("/recommend", v.videoRecommend)
	video.GET("/addViews/:vid", v.addViews)
}
