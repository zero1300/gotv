package main

import (
	"gotv/common"
	"gotv/load"
	"gotv/module/comment"
	"gotv/module/file"
	"gotv/module/user"
	"gotv/module/video"
	"gotv/router"
)

func main() {
	// init
	db := load.NewDatabase()
	rc := load.NewRedisClient()
	mc := load.NewMinioClient()
	ginRouter := router.NewGinRouter()
	ginRouter.Engine.Use(common.Cors())

	admin := ginRouter.Engine.Group("/admin")
	admin.Use(common.AdminAuth)

	api := ginRouter.Engine.Group("api")

	// user
	userDao := user.NewUserDao(db)
	uerHandler := user.NewUserHandler(userDao, rc)
	uerHandler.SetUp(admin, api)

	api.Use(common.Auth)
	uerHandler.SetUp2(admin, api)

	// video
	videoDao := video.NewVideoDao(db)
	videoHandler := video.NewVideoHandler(videoDao, rc)
	videoHandler.SetUp(admin, api)

	// file
	fileDao := file.NewFileDao(db)
	fileController := file.NewFileController(fileDao, rc, mc)
	fileController.SetUp(admin, api)

	// comment
	commentDao := comment.NewCommentDao(db)
	commentController := comment.NewCommentController(commentDao)
	commentController.SetUp(admin, api)
	ginRouter.Engine.Run()
}
