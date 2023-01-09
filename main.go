package main

import (
	"gotv/common"
	"gotv/load"
	"gotv/module/user"
	"gotv/module/video"
	"gotv/router"
)

func main() {
	// init
	db := load.NewDatabase()
	rc := load.NewRedisClient()
	ginRouter := router.NewGinRouter()

	admin := ginRouter.Engine.Group("/admin")
	admin.Use(common.Cors())
	admin.Use(common.AdminAuth)

	api := ginRouter.Engine.Group("api")
	api.Use(common.Cors())

	// user
	userDao := user.NewUserDao(db)
	uerHandler := user.NewUserHandler(userDao, rc)
	uerHandler.SetUp(admin, api)

	api.Use(common.Auth)

	// video
	videoDao := video.NewVideoDao(db)
	videoHandler := video.NewVideoHandler(videoDao, rc)
	videoHandler.SetUp(admin, api)

	ginRouter.Engine.Run()
}
