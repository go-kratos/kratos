package popular

import (
	"context"
	"go-common/app/admin/main/feed/model/common"
	"go-common/app/admin/main/feed/model/show"
	"go-common/app/admin/main/feed/util"
	"go-common/library/log"
)

const (
	//ActionAddCEventTopic log action
	ActionAddCEventTopic = "ActAddEventTopic"
	//ActionUpCEventTopic log action
	ActionUpCEventTopic = "ActUpEventTopic"
	//ActionDelCEventTopic log action
	ActionDelCEventTopic = "ActDelEventTopic"
)

//EventTopicList channel EventTopic list
func (s *Service) EventTopicList(lp *show.EventTopicLP) (pager *show.EventTopicPager, err error) {
	pager = &show.EventTopicPager{
		Page: common.Page{
			Num:  lp.Pn,
			Size: lp.Ps,
		},
	}
	w := map[string]interface{}{
		"deleted": common.NotDeleted,
	}
	query := s.showDao.DB.Model(&show.EventTopic{})
	if lp.ID > 0 {
		w["id"] = lp.ID
	}
	if lp.Person != "" {
		query = query.Where("person like ?", "%"+lp.Person+"%")
	}
	if lp.Title != "" {
		query = query.Where("title like ?", "%"+lp.Title+"%")
	}
	if err = query.Where(w).Count(&pager.Page.Total).Error; err != nil {
		log.Error("popularSvc.EventTopicList count error(%v)", err)
		return
	}
	EventTopics := make([]*show.EventTopic, 0)
	if err = query.Where(w).Order("`id` DESC").Offset((lp.Pn - 1) * lp.Ps).Limit(lp.Ps).Find(&EventTopics).Error; err != nil {
		log.Error("popularSvc.EventTopicList Find error(%v)", err)
		return
	}
	pager.Item = EventTopics
	return
}

//AddEventTopic add channel EventTopic
func (s *Service) AddEventTopic(c context.Context, param *show.EventTopicAP, name string, uid int64) (err error) {
	if err = s.showDao.EventTopicAdd(param); err != nil {
		return
	}
	if err = util.AddLogs(common.LogEventTopic, name, uid, 0, ActionAddCEventTopic, param); err != nil {
		log.Error("popularSvc.AddEventTopic AddLog error(%v)", err)
		return
	}
	return
}

//UpdateEventTopic update channel EventTopic
func (s *Service) UpdateEventTopic(c context.Context, param *show.EventTopicUP, name string, uid int64) (err error) {
	if err = s.showDao.EventTopicUpdate(param); err != nil {
		return
	}
	if err = util.AddLogs(common.LogEventTopic, name, uid, 0, ActionUpCEventTopic, param); err != nil {
		log.Error("popularSvc.UpdateEventTopic AddLog error(%v)", err)
		return
	}
	return
}

//DeleteEventTopic delete channel EventTopic
func (s *Service) DeleteEventTopic(id int64, name string, uid int64) (err error) {
	if err = s.showDao.EventTopicDelete(id); err != nil {
		return
	}
	if err = util.AddLogs(common.LogEventTopic, name, uid, id, ActionDelCEventTopic, id); err != nil {
		log.Error("popularSvc.DeleteEventTopic AddLog error(%v)", err)
		return
	}
	return
}
