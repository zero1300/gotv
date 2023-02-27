package history

import (
	"github.com/gin-gonic/gin"
	"gotv/model"
	"gotv/resp"
)

type historyController struct {
	historyDao *HistoryDao
}

func NewHistory(dao *HistoryDao) *historyController {
	return &historyController{
		historyDao: dao,
	}
}

func buildHistory(ctx *gin.Context) (model.History, error) {
	var history model.History
	err := ctx.ShouldBind(&history)
	if err != nil {
		return model.History{}, err
	}
	obj, _ := ctx.Get("user")
	user := obj.(*model.User)
	history.UID = user.ID
	return history, nil
}

func (h historyController) addHistory(ctx *gin.Context) {
	history, err := buildHistory(ctx)
	if err != nil {
		resp.Fail(ctx, "Form 表单参数异常: "+err.Error())
	}
	err = h.historyDao.addHistory(history)
	if err != nil {
		resp.Fail(ctx, err.Error())
	}
}

func (h historyController) historyList(ctx *gin.Context) {
	obj, _ := ctx.Get("user")
	user := obj.(*model.User)
	var p model.Page
	if err := ctx.ShouldBind(&p); err != nil {
		resp.Fail(ctx, "参数异常: "+err.Error())
		return
	}
	list, _ := h.historyDao.historyList(p, user.ID)
	resp.Success(ctx, list)
}

func (h historyController) Setup(admin *gin.RouterGroup, api *gin.RouterGroup) {
	history := api.Group("history")
	history.POST("/addHistory", h.addHistory)
	history.POST("/list", h.historyList)
}
