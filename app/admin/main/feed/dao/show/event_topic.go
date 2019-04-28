package show

import (
	"go-common/app/admin/main/feed/model/common"
	"go-common/app/admin/main/feed/model/show"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
)

// EventTopicAdd add event topic
func (d *Dao) EventTopicAdd(param *show.EventTopicAP) (err error) {
	if err = d.DB.Create(param).Error; err != nil {
		log.Error("dao.show.EventTopicAdd error(%v)", err)
		return
	}
	return
}

// EventTopicUpdate update event topic
func (d *Dao) EventTopicUpdate(param *show.EventTopicUP) (err error) {
	if err = d.DB.Model(&show.EventTopicUP{}).Update(param).Error; err != nil {
		log.Error("dao.show.EventTopicUpdate error(%v)", err)
		return
	}
	return
}

// EventTopicDelete delete cevent topic
func (d *Dao) EventTopicDelete(id int64) (err error) {
	up := map[string]interface{}{
		"deleted": common.Deleted,
	}
	if err = d.DB.Model(&show.EventTopic{}).Where("id = ?", id).Update(up).Error; err != nil {
		log.Error("dao.show.EventTopicDelete error(%v)", err)
		return
	}
	return
}

// ETFindByID event topic table find by id
func (d *Dao) ETFindByID(id int64) (topic *show.EventTopic, err error) {
	topic = &show.EventTopic{}
	w := map[string]interface{}{
		"deleted": common.NotDeleted,
	}
	if err = d.DB.Model(&show.EventTopic{}).Where("id = ?", id).Where(w).First(topic).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			topic = nil
			err = nil
		} else {
			log.Error("dao.ormshow.event_topic.findByID error(%v)", err)
		}
	}
	return
}
