package service

import (
	"fmt"
	"time"

	"go-common/app/admin/main/space/model"
	"go-common/library/log"
	"go-common/library/queue/databus/report"

	"github.com/jinzhu/gorm"
)

const (
	_BlacklistAdd   = "BlacklistAdd"
	_BlacklistUp    = "BlacklistUp"
	_StatusNotBlack = 1
)

// BlacklistAdd add blacklist
func (s *Service) BlacklistAdd(mids []int64, name string, uid int64) (err error) {
	var (
		blacklist           map[int64]*model.Blacklist
		updateMids, addMids []int64
	)
	if len(mids) > 30 {
		err = fmt.Errorf("黑名单一次最多只能添加30个")
		return
	}
	if blacklist, err = s.dao.BlacklistIn(mids); err != nil {
		return
	}
	for _, v := range mids {
		if blacklist[v] != nil {
			if blacklist[v].Status == _StatusNotBlack {
				updateMids = append(updateMids, blacklist[v].ID)
			}
		} else {
			addMids = append(addMids, v)
		}
	}
	if err = s.dao.BlacklistAdd(addMids, updateMids); err != nil {
		return
	}
	for _, v := range mids {
		if err = report.Manager(&report.ManagerInfo{
			Uname:    name,
			UID:      uid,
			Business: model.NoticeLogID,
			Type:     model.LogBlacklist,
			Oid:      v,
			Action:   _BlacklistAdd,
			Ctime:    time.Now(),
			Content: map[string]interface{}{
				"mids": mids,
			},
		}); err != nil {
			return
		}
	}
	return
}

// BlacklistUp update blacklist
func (s *Service) BlacklistUp(id int64, status int, name string, uid int64) (err error) {
	var (
		mids []int64
	)
	blacklist := &model.Blacklist{}
	if err = s.dao.DB.Model(&model.Blacklist{}).Where("id in (?)", []int64{id}).First(blacklist).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = fmt.Errorf("找不用数据")
			return
		}
		log.Error("Srv.BlacklistUp First error(%v)", err)
		return
	}
	if err = s.dao.BlacklistUp(id, status); err != nil {
		return
	}
	mids = []int64{blacklist.Mid}
	if err = report.Manager(&report.ManagerInfo{
		Uname:    name,
		UID:      uid,
		Business: model.NoticeLogID,
		Type:     model.LogBlacklist,
		Oid:      blacklist.Mid,
		Action:   _BlacklistUp,
		Ctime:    time.Now(),
		Content: map[string]interface{}{
			"id":     id,
			"status": status,
			"mids":   mids,
		},
	}); err != nil {
		return
	}
	return
}

// BlacklistIndex .
func (s *Service) BlacklistIndex(mid int64, pn, ps int) (pager *model.BlacklistPager, err error) {
	if pager, err = s.dao.BlacklistIndex(mid, pn, ps); err != nil {
		return
	}
	return
}
