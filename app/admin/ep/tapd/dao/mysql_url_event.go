package dao

import (
	"go-common/app/admin/ep/tapd/model"
	"go-common/library/ecode"
)

// AddURLEvent Add URL Event.
func (d *Dao) AddURLEvent(urlEvent *model.UrlEvent) error {
	return d.db.Create(urlEvent).Error
}

// QueryURLEventByUrlAndEvent Query URL Event By Url And Event.
func (d *Dao) QueryURLEventByUrlAndEvent(urlID int64, eventType string) (urlEvents []*model.UrlEvent, err error) {
	err = d.db.Model(&model.HookUrl{}).Where("url_id = ? and event = ?", urlID, eventType).Find(&urlEvents).Error
	if err == ecode.NothingFound {
		err = nil
	}
	return
}

// QueryURLEventByUrl Query URL Event By Url.
func (d *Dao) QueryURLEventByUrl(urlID int64) (urlEvents []*model.UrlEvent, err error) {
	err = d.db.Model(&model.HookUrl{}).Where("url_id = ?", urlID).Find(&urlEvents).Error
	if err == ecode.NothingFound {
		err = nil
	}
	return
}

// QueryURLEventByEventAndStatus Query URL Event By Event and status.
func (d *Dao) QueryURLEventByEventAndStatus(event string, status int) (urlEvents []*model.UrlEvent, err error) {
	err = d.db.Model(&model.HookUrl{}).Where("event = ? and status = ?", event, status).Find(&urlEvents).Error
	if err == ecode.NothingFound {
		err = nil
	}
	return
}

// QueryURLEventByStatus Query URL Event By Status.
func (d *Dao) QueryURLEventByStatus(status int) (urlEvents []*model.UrlEvent, err error) {
	err = d.db.Model(&model.HookUrl{}).Where("status = ?", status).Find(&urlEvents).Error
	if err == ecode.NothingFound {
		err = nil
	}
	return
}

// UpdateURLEventStatus Update URL Event status.
func (d *Dao) UpdateURLEventStatus(id int64, status int) error {
	return d.db.Model(&model.UrlEvent{}).Where("id=?", id).Update("status", status).Error
}
