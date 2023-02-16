package video

import (
	"context"

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

func (v *VideoHandler) getOne(ctx *gin.Context) {
	v.videoDao.GetVideoById(123)
}

func (v *VideoHandler) SetUp(admin *gin.RouterGroup, api *gin.RouterGroup) {

	video := api.Group("/video")
	video.GET("/getOne", v.getOne)
}
