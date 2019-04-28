package show

import (
	"go-common/app/admin/main/feed/model/common"
	"go-common/app/admin/main/feed/model/show"
	"go-common/library/log"
)

// PopularStarsAdd add popuar stars card
func (d *Dao) PopularStarsAdd(param *show.PopularStarsAP) (err error) {
	if err = d.DB.Create(param).Error; err != nil {
		log.Error("dao.show.PopularStarsAdd error(%v)", err)
		return
	}
	return
}

// PopularStarsUpdate update popuar stars card
func (d *Dao) PopularStarsUpdate(param *show.PopularStarsUP) (err error) {
	if err = d.DB.Model(&show.PopularStarsUP{}).Update(param).Error; err != nil {
		log.Error("dao.show.PopularStarsUpdate error(%v)", err)
		return
	}
	return
}

// PopularStarsDelete delete popuar stars card
func (d *Dao) PopularStarsDelete(id int64, t string) (err error) {
	up := map[string]interface{}{
		"deleted": common.Deleted,
	}
	w := map[string]interface{}{
		"id":   id,
		"type": t,
	}
	if err = d.DB.Model(&show.PopularStars{}).Where(w).Update(up).Error; err != nil {
		log.Error("dao.show.PopularStarsDelete error(%v)", err)
		return
	}
	return
}

// PopularStarsReject reject popuar stars card
func (d *Dao) PopularStarsReject(id int64, t string) (err error) {
	up := map[string]interface{}{
		"status": common.Rejecte,
	}
	w := map[string]interface{}{
		"id":   id,
		"type": t,
	}
	if err = d.DB.Model(&show.PopularStars{}).Where(w).Update(up).Error; err != nil {
		log.Error("dao.show.PopularStarsReject error(%v)", err)
		return
	}
	return
}
