package resp

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response
// context 上下文
// httpStatus http 状态码
// code 自己定义的状态码
// data 返回的空接口
// msg 返回的信息
func Response(context *gin.Context, code int, data any, msg string) {
	context.JSON(http.StatusOK, gin.H{
		"code": code,
		"data": data,
		"msg":  msg,
	})
}

func Success(context *gin.Context, data any) {
	context.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": data,
		"msg":  "ok",
	})
}

func Fail(context *gin.Context, msg string) {
	context.JSON(http.StatusOK, gin.H{
		"code": 1,
		"data": nil,
		"msg":  msg,
	})
}
