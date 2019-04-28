package channel

import (
	"context"
	"fmt"
	"time"

	"go-common/app/admin/main/feed/model/common"
	"go-common/app/admin/main/feed/model/show"
	"go-common/app/admin/main/feed/util"
	"go-common/library/log"
)

const (
	//TabOnline channel tab online
	TabOnline = 1
	//TabDownline channel tab down line
	TabDownline = 2
	//TabWaitOnline channel tab wait to line
	TabWaitOnline = 3
	//OrderTimeDown channel tab oder by stime desc
	OrderTimeDown = 1
	//OrderTimeUp channel tab oder by stime asc
	OrderTimeUp = 2
	//ActionAddCTab log action
	ActionAddCTab = "ActAddChannelTab"
	//ActionUpCTab log action
	ActionUpCTab = "ActUpChannelTab"
	//ActionDelCTab log action
	ActionDelCTab = "ActDelChannelTab"
	//ActionOfflineCTab log action
	ActionOfflineCTab = "ActOfflineChannelTab"
)

//TabList channel tab list
func (s *Service) TabList(lp *show.ChannelTabLP) (pager *show.ChannelTabPager, err error) {
	var (
		eTime int64
		sTime int64
	)
	pager = &show.ChannelTabPager{
		Page: common.Page{
			Num:  lp.Pn,
			Size: lp.Ps,
		},
	}
	w := map[string]interface{}{
		"is_delete": common.NotDeleted,
	}
	query := s.showDao.DB.Model(&show.ChannelTab{})
	if lp.TagID > 0 {
		w["tag_id"] = lp.TagID
	}
	if lp.TabID > 0 {
		w["tab_id"] = lp.TabID
	}
	if lp.Stime > 0 {
		query = query.Where("stime >= ?", lp.Stime)
	}
	if lp.Etime > 0 {
		query = query.Where("etime <= ?", lp.Etime)
	}
	if lp.Person != "" {
		query = query.Where("person like ?", "%"+lp.Person+"%")
	}
	if lp.Status == TabWaitOnline {
		if lp.Stime != 0 {
			if lp.Stime < time.Now().Unix() {
				sTime = time.Now().Unix()
			} else {
				sTime = lp.Stime
			}
		} else {
			sTime = time.Now().Unix()
		}
		query = query.Where("stime >= ?", sTime)
	} else if lp.Status == TabOnline {
		if lp.Stime != 0 {
			if lp.Stime < time.Now().Unix() {
				sTime = time.Now().Unix()
			} else {
				sTime = lp.Stime
			}
		} else {
			sTime = time.Now().Unix()
		}
		if lp.Etime != 0 {
			if lp.Etime > time.Now().Unix() {
				eTime = time.Now().Unix()
			} else {
				eTime = lp.Etime
			}
		} else {
			eTime = time.Now().Unix()
		}
		query = query.Where("stime < ?", sTime).Where("etime >= ?", eTime)
	} else if lp.Status == TabDownline {
		if lp.Etime != 0 {
			if lp.Etime < time.Now().Unix() {
				eTime = lp.Etime
			} else {
				eTime = time.Now().Unix()
			}
		} else {
			eTime = time.Now().Unix()
		}
		query = query.Where("etime < ?", eTime)
	}
	if lp.Order == OrderTimeDown {
		query = query.Order("`stime` ASC")
	} else if lp.Order == OrderTimeUp {
		query = query.Order("`stime` DESC")
	}
	if err = query.Where(w).Count(&pager.Page.Total).Error; err != nil {
		log.Error("chanelSvc.CardSetupList Index count error(%v)", err)
		return
	}

	tabs := []*show.ChannelTab{}
	if err = query.Where(w).Order("`id` DESC").Offset((lp.Pn - 1) * lp.Ps).Limit(lp.Ps).Find(&tabs).Error; err != nil {
		log.Error("chanelSvc.CardSetupList First error(%v)", err)
		return
	}
	for k, v := range tabs {
		//online for fe
		if time.Now().Unix() < v.Stime {
			tabs[k].Status = TabWaitOnline
		} else if time.Now().Unix() >= v.Etime {
			tabs[k].Status = TabDownline
		} else {
			tabs[k].Status = TabOnline
		}
	}
	pager.Item = tabs
	return
}

//AddTab add channel tab
func (s *Service) AddTab(c context.Context, param *show.ChannelTabAP, name string, uid int64) (err error) {
	if err = s.IsValid(0, param.TagID, param.Stime, param.Etime, param.Priority); err != nil {
		return
	}
	if err = s.showDao.ChannelTabAdd(param); err != nil {
		return
	}
	if err = util.AddLogs(common.LogChannelTab, name, uid, 0, ActionAddCTab, param); err != nil {
		log.Error("chanelSvc.AddTab AddLog error(%v)", err)
		return
	}
	return
}

//IsValid validate data
func (s *Service) IsValid(id, tagID, sTime int64, eTime int64, priority int) (err error) {
	var (
		count int
	)
	if sTime > eTime {
		err = fmt.Errorf("开始时间不能大于结束时间")
		return
	}
	if sTime < time.Now().Unix() {
		err = fmt.Errorf("生效时间需要大于当前时间")
		return
	}
	if count, err = s.showDao.ChannelTabValid(id, tagID, sTime, eTime, priority); err != nil {
		return
	}
	if count > 0 {
		err = fmt.Errorf("已有该排序无法创建，请重新选择")
		return
	}
	if count, err = s.showDao.ChannelTabValid(id, tagID, sTime, eTime, 0); err != nil {
		return
	}
	if count >= 3 {
		stimeStr := time.Unix(sTime, 0).Format("2006-01-02 15:04:05")
		etimeStr := time.Unix(eTime, 0).Format("2006-01-02 15:04:05")
		str := "频道在" + stimeStr + " 至 " + etimeStr + " 时间段内已有3个运营tab，无法创建"
		err = fmt.Errorf(str)
		return
	}
	return
}

//UpdateTab update channel tab
func (s *Service) UpdateTab(c context.Context, param *show.ChannelTabUP, name string, uid int64) (err error) {
	if err = s.IsValid(param.ID, param.TagID, param.Stime, param.Etime, param.Priority); err != nil {
		return
	}
	if err = s.showDao.ChannelTabUpdate(param); err != nil {
		return
	}
	if err = util.AddLogs(common.LogChannelTab, name, uid, 0, ActionUpCTab, param); err != nil {
		log.Error("chanelSvc.UpdateTab AddLog error(%v)", err)
		return
	}
	return
}

//DeleteTab delete channel tab
func (s *Service) DeleteTab(id int64, name string, uid int64) (err error) {
	if err = s.showDao.ChannelTabDelete(id); err != nil {
		return
	}
	if err = util.AddLogs(common.LogChannelTab, name, uid, id, ActionDelCTab, id); err != nil {
		log.Error("chanelSvc.DeleteTab AddLog error(%v)", err)
		return
	}
	return
}

//OfflineTab offline channel tab
func (s *Service) OfflineTab(id int64, name string, uid int64) (err error) {
	if err = s.showDao.ChannelTabOffline(id); err != nil {
		return
	}
	if err = util.AddLogs(common.LogChannelTab, name, uid, id, ActionOfflineCTab, id); err != nil {
		log.Error("chanelSvc.DeleteTab AddLog error(%v)", err)
		return
	}
	return
}
