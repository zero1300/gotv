package history

import (
	"github.com/bwmarrin/snowflake"
	"gorm.io/gorm"
	"gotv/model"
	"gotv/model/vo"
	"gotv/module/user"
	"gotv/module/video"
	"gotv/resp"
	"strconv"
	"time"
)

type HistoryDao struct {
	db       *gorm.DB
	videoDao *video.VideoDao
	userDao  *user.UserDao
}

func NewHistoryDao(db *gorm.DB, videoDao *video.VideoDao, userDao *user.UserDao) *HistoryDao {
	return &HistoryDao{
		db:       db,
		videoDao: videoDao,
		userDao:  userDao,
	}
}

func (h HistoryDao) addHistory(history model.History) error {
	historyPo := h.getHistory(history.VID, history.UID)
	if historyPo.ID == 0 {
		node, _ := snowflake.NewNode(1)
		id := node.Generate().Int64()
		history.ID = uint(id)
		return h.db.Create(&history).Error
	}
	return h.db.Debug().Model(model.History{}).Where("id = ?", historyPo.ID).Updates(model.History{Position: history.Position, CreateTime: time.Now()}).Error
}

func (h HistoryDao) getHistory(vid string, uid uint) model.History {
	var history model.History
	h.db.Model(model.History{}).Where("vid = ? and uid = ?", vid, uid).Find(&history)
	return history
}

func (h HistoryDao) getHistoryVo(vid string, uid uint) vo.HistoryVo {
	var historyVo vo.HistoryVo

	history := h.getHistory(vid, uid)
	videoPo := h.videoDao.GetVideoById(vid)
	userPo := h.userDao.GetUser(strconv.Itoa(int(uid)))

	historyVo.Position = history.Position
	historyVo.CreateTime = history.CreateTime

	historyVo.ID = videoPo.ID
	historyVo.Uri = videoPo.Uri
	historyVo.Title = videoPo.Title
	historyVo.Cover = videoPo.Cover
	historyVo.Duration = videoPo.Duration
	historyVo.Views = videoPo.Views

	historyVo.Nickname = userPo.Nickname
	historyVo.UID = videoPo.UID

	return historyVo
}

func (h HistoryDao) historyList(p model.Page, uid uint) (resp.Pager, error) {
	histories := make([]model.History, 10)
	h.db.Where("uid = ?", uid).Find(&histories).Limit(p.PageSize).Offset((p.PageNum - 1) * p.PageSize)

	historyVos := make([]vo.HistoryVo, 0)
	for _, history := range histories {
		historyVo := h.getHistoryVo(history.VID, history.UID)
		historyVos = append(historyVos, historyVo)
	}

	pager := resp.Pager{}
	pager.List = historyVos
	pager.Total = int64(len(historyVos))
	return pager, nil
}

func (h HistoryDao) DelHistoryByVid(vid string) {
	var history model.History
	h.db.Delete(history, "vid = ?", vid)
}
