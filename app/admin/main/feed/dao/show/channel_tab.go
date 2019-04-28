package show

import (
	"time"

	"go-common/app/admin/main/feed/model/common"
	"go-common/app/admin/main/feed/model/show"
	"go-common/library/log"
)

// ChannelTabAdd add channel tab
func (d *Dao) ChannelTabAdd(param *show.ChannelTabAP) (err error) {
	if err = d.DB.Create(param).Error; err != nil {
		log.Error("dao.show.ChannelTabAdd error(%v)", err)
		return
	}
	return
}

// ChannelTabUpdate update channel tab
func (d *Dao) ChannelTabUpdate(param *show.ChannelTabUP) (err error) {
	if err = d.DB.Model(&show.ChannelTabUP{}).Update(param).Error; err != nil {
		log.Error("dao.show.ChannelTabUpdate error(%v)", err)
		return
	}
	return
}

// ChannelTabDelete delete channel tab
func (d *Dao) ChannelTabDelete(id int64) (err error) {
	up := map[string]interface{}{
		"is_delete": common.Deleted,
	}
	if err = d.DB.Model(&show.ChannelTab{}).Where("id = ?", id).Update(up).Error; err != nil {
		log.Error("dao.show.ChannelTabDelete error(%v)", err)
		return
	}
	return
}

// ChannelTabOffline offline channel tab
func (d *Dao) ChannelTabOffline(id int64) (err error) {
	up := map[string]interface{}{
		"etime": time.Now().Unix(),
	}
	if err = d.DB.Model(&show.ChannelTab{}).Where("id = ?", id).Update(up).Error; err != nil {
		log.Error("dao.show.ChannelTabOffline error(%v)", err)
		return
	}
	return
}

// ChannelTabValid conflict check
func (d *Dao) ChannelTabValid(id, tagID, sTime int64, eTime int64, priority int) (count int, err error) {
	w := map[string]interface{}{
		"is_delete": common.NotDeleted,
		"tag_id":    tagID,
	}
	if priority != 0 {
		w["priority"] = priority
	}
	query := d.DB.Model(&show.ChannelTab{}).Where("stime < ?", eTime).Where("etime > ?", sTime)
	if id != 0 {
		query = query.Where("id != ?", id)
	}
	if err = query.Where(w).Count(&count).Error; err != nil {
		log.Error("dao.show.ChannelTabValid error(%v)", err)
		return
	}
	return
}
