package dao

import (
	"fmt"

	"go-common/app/admin/main/space/model"
	"go-common/library/log"
)

// BlacklistAdd add blacklist
func (d *Dao) BlacklistAdd(addmids, upmids []int64) (err error) {
	tx := d.DB.Begin()
	if err = tx.Error; err != nil {
		log.Error("dao.BlacklistAdd.Begin error(%v)", err)
		return
	}
	if len(addmids) > 0 {
		if err = tx.Model(&model.Blacklist{}).Exec(model.BlacklistBatchAddSQL(addmids)).Error; err != nil {
			log.Error("dao.BlacklistAdd.BlacklistBatchAddSQL(%+v) error(%v)", addmids, err)
			err = tx.Rollback().Error
			return
		}
	}
	if len(upmids) > 0 {
		if err = tx.Model(&model.Blacklist{}).Exec(model.BlacklistBatchUpdateSQL(upmids)).Error; err != nil {
			log.Error("dao.BlacklistAdd.BlacklistBatchUpdateSQL(%+v) error(%v)", upmids, err)
			err = tx.Rollback().Error
			return
		}
	}
	err = tx.Commit().Error
	return
}

// BlacklistIn query blackist count
func (d *Dao) BlacklistIn(mids []int64) (blacks map[int64]*model.Blacklist, err error) {
	var (
		blacklist []*model.Blacklist
	)
	blacks = make(map[int64]*model.Blacklist, len(mids))
	if len(mids) == 0 {
		return nil, fmt.Errorf("mid不能为空")
	}
	if err = d.DB.Model(&model.Blacklist{}).Where("mid in (?)", mids).Find(&blacklist).
		Error; err != nil {
		log.Error("dao.BlacklistIn.Count(%+v) error(%v)", mids, err)
		return
	}
	for _, v := range blacklist {
		blacks[v.Mid] = v
	}
	return
}

// BlacklistUp blackist update
func (d *Dao) BlacklistUp(id int64, status int) (err error) {
	w := map[string]interface{}{
		"id": id,
	}
	up := map[string]interface{}{
		"status": status,
	}
	if err = d.DB.Model(&model.Blacklist{}).Where(w).Update(up).Error; err != nil {
		log.Error("dao.BlacklistUp.Update error(%v)", err)
		return
	}
	return
}

// BlacklistIndex blackist
func (d *Dao) BlacklistIndex(mid int64, pn, ps int) (pager *model.BlacklistPager, err error) {
	var (
		blacklist []*model.Blacklist
	)
	pager = &model.BlacklistPager{
		Page: model.Page{
			Num:  pn,
			Size: ps,
		},
	}
	query := d.DB.Model(&model.Blacklist{})
	if mid != 0 {
		query = query.Where("mid = ?", mid)
	}
	if err = query.Count(&pager.Page.Total).Error; err != nil {
		log.Error("dao.BlacklistIndex.Count error(%v)", err)
		return
	}
	if err = query.Order("`id` DESC").Offset((pn - 1) * ps).Limit(ps).Find(&blacklist).Error; err != nil {
		log.Error("dao.BlacklistIndex.Find error(%v)", err)
		return
	}
	pager.Item = blacklist
	return
}
