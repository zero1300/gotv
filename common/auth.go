package common

import (
	"context"
	"fmt"
	"gotv/load"
	"gotv/model"
	"gotv/resp"

	"github.com/gin-gonic/gin"
)

// 前台业务系统校验
func Auth(ctx *gin.Context) {
	context := context.Background()
	token := ctx.GetHeader("token")
	if token == "" || len(token) != 64 {
		ctx.JSON(203, gin.H{"msg": "token不能为空或token格式异常"})
		ctx.Abort()
		return
	}

	rc := load.NewRedisClient()
	key := "token_" + token
	ret, _ := rc.Get(context, key).Result()

	if ret == "" {
		resp.Fail(ctx, "token无效1")
		ctx.Abort()
		return
	}
	user := new(model.User)
	user.UnmarshalBinary([]byte(ret))
	fmt.Println(user)
	if user.ID == 0 {
		resp.Fail(ctx, "token无效2")
		ctx.Abort()
		return
	}
	ctx.Set("user", user)
}

// 后台管理系统校验
func AdminAuth(ctx *gin.Context) {
	context := context.Background()
	token := ctx.GetHeader("token")
	if token == "" || len(token) != 64 {
		resp.Fail(ctx, "token不能为空或token格式异常")
		ctx.Abort()
		return
	}
	rc := load.NewRedisClient()
	key := "token_" + token
	ret, _ := rc.Get(context, key).Result()
	if ret == "" {
		resp.Fail(ctx, "token无效1")
		ctx.Abort()
		return
	}
	user := new(model.User)
	user.UnmarshalBinary([]byte(ret))
	fmt.Println(ret)
	if user.ID == 0 {
		resp.Fail(ctx, "token无效2")
		ctx.Abort()
		return
	}
	res, err := rc.SIsMember(context, "admins", user.ID).Result()
	fmt.Println(res)
	if err != nil || !res {
		resp.Fail(ctx, "非管理员")
		ctx.Abort()
		return
	}
	ctx.Set("user", user)
}
