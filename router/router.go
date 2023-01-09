package router

import (
	"github.com/gin-gonic/gin"
)

type GinRouter struct {
	Engine *gin.Engine
}

func NewGinRouter() *GinRouter {
	httpRouter := gin.Default()
	return &GinRouter{
		Engine: httpRouter,
	}
}
