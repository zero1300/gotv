package model

import "github.com/gin-gonic/gin"

type Page struct {
	PageNum  int    `form:"page"`
	PageSize int    `form:"pageSize"`
	Keyword  string `form:"keyword"`
}

func GetPage(ctx *gin.Context) (Page, error) {
	var page Page
	err := ctx.ShouldBind(&page)
	if err != nil {
		return Page{}, err
	}
	return page, nil
}
