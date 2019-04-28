package dao

import (
	"go-common/app/admin/ep/tapd/model"
	"go-common/library/ecode"
)

const _wildcards = "%"

// AddHookURL Add Hook URL.
func (d *Dao) AddHookURL(hookURL *model.HookUrl) error {
	return d.db.Create(hookURL).Error
}

// UpdateHookURL Update Hook URL.
func (d *Dao) UpdateHookURL(hookURL *model.HookUrl) error {
	return d.db.Model(&model.HookUrl{}).Where("id=?", hookURL.ID).Update(hookURL).Error
}

// QueryHookURLByID Query Hook URL By ID.
func (d *Dao) QueryHookURLByID(id int64) (hookURL *model.HookUrl, err error) {
	hookURL = &model.HookUrl{}
	err = d.db.Model(&model.HookUrl{}).Where("id = ?", id).First(hookURL).Error
	if err == ecode.NothingFound {
		err = nil
	}
	return
}

// AddHookURLandEvent Add Hook URL and Event.
func (d *Dao) AddHookURLandEvent(hookURL *model.HookUrl, urlEvents []*model.UrlEvent) (err error) {
	tx := d.db.Begin()

	if err = tx.Error; err != nil {
		return
	}

	if err = tx.Create(hookURL).Error; err != nil {
		tx.Rollback()
		return
	}

	for _, urlEvent := range urlEvents {
		urlEvent.UrlID = hookURL.ID
		if err = tx.Create(urlEvent).Error; err != nil {
			tx.Rollback()
			return
		}
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
	}
	return
}

// UpdateHookURLandEvent Update Hook URL and Event.
func (d *Dao) UpdateHookURLandEvent(hookURL *model.HookUrl, urlEvents []*model.UrlEvent) (err error) {
	tx := d.db.Begin()
	if err = tx.Error; err != nil {
		return
	}

	if err = tx.Model(model.HookUrl{}).Where("id=?", hookURL.ID).
		Updates(map[string]interface{}{"url": hookURL.URL, "workspace_id": hookURL.WorkspaceID, "status": hookURL.Status, "update_by": hookURL.UpdateBy}).
		Error; err != nil {

		tx.Rollback()
		return
	}

	for _, urlEvent := range urlEvents {
		if urlEvent.ID != 0 {
			//update
			if err = tx.Model(model.UrlEvent{}).Where("id=?", urlEvent.ID).Update(urlEvent).Error; err != nil {
				tx.Rollback()
				return
			}
		} else {
			//add
			if err = tx.Create(urlEvent).Error; err != nil {
				tx.Rollback()
				return
			}
		}
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
	}
	return
}

//FindHookURLs Find Hook URLs.
func (d *Dao) FindHookURLs(req *model.QueryHookURLReq) (total int64, hookURLs []*model.HookUrl, err error) {
	gDB := d.db.Model(&model.HookUrl{})

	if req.ID > 0 {
		gDB = gDB.Where("id=?", req.ID)
	}
	if req.Status > 0 {
		gDB = gDB.Where("status=?", req.Status)
	}

	if req.UpdateBy != "" {
		gDB = gDB.Where("update_by=?", req.UpdateBy)
	}

	if req.URL != "" {
		gDB = gDB.Where("url like ?", req.URL+_wildcards)
	}

	if err = gDB.Count(&total).Error; err != nil {
		return
	}

	err = gDB.Order("ctime desc").Offset((req.PageNum - 1) * req.PageSize).Limit(req.PageSize).Find(&hookURLs).Error

	return
}
