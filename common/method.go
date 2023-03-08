package common

import (
	"github.com/gin-gonic/gin"
	"gotv/model"
)

func GetPage(ctx *gin.Context) (model.Page, error) {
	var page model.Page
	err := ctx.ShouldBind(&page)
	if err != nil {
		return model.Page{}, err
	}
	return page, nil
}
