package service

import (
	"database/sql"

	"go-common/app/admin/main/tv/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"

	"github.com/jinzhu/gorm"
)

const (
	_seasonPassed  = 1
	_seasonToAudit = 2
	_epPassed      = 3
	_epToAudit     = 1
	_seasonDefault = 4
)

// SeasonCheck manages the season's check status
func (s *Service) SeasonCheck(sn *model.TVEpSeason) (err error) {
	var status int
	if sn.Check != _seasonPassed {
		status = _seasonToAudit
	} else {
		return
	}
	if err = s.DB.Model(&model.TVEpSeason{}).Where("id = ?", sn.ID).Update(map[string]int{"check": status}).Error; err != nil {
		log.Error("tvSrv.SeasonCheck error(%v)", err)
	}
	return
}

// EpCheck manages the ep's check status
func (s *Service) EpCheck(ep *model.Content) (err error) {
	var status int
	if ep.State != _epPassed {
		status = _epToAudit
	} else {
		return // if passed, no need to update the status
	}
	if err = s.DB.Model(&model.Content{}).Where("id = ?", ep.ID).Update(map[string]int{"state": status}).Error; err != nil {
		log.Error("tvSrv.EpCheck error(%v)", err)
	}
	return
}

// EpDel controls the ep's deletion status, 0 or 1
func (s *Service) EpDel(epid int64, action int) (err error) {
	if err = s.DB.Model(&model.TVEpContent{}).Where("id=?", epid).Update(map[string]int{"is_deleted": action}).Error; err != nil {
		log.Error("tvSrv.EpDel error(%v)\n", err)
		return
	}
	if err = s.DB.Model(&model.Content{}).Where("epid=?", epid).Update(map[string]int{"is_deleted": action}).Error; err != nil {
		log.Error("tvSrv.EpDel error(%v)\n", err)
	}
	return
}

// SeasonRemove removes the season and its eps
func (s *Service) SeasonRemove(season *model.TVEpSeason) (err error) {
	var rows *sql.Rows
	// remove season
	if err = s.DB.Model(&model.TVEpSeason{}).Where("id = ?", season.ID).Update(map[string]int{"is_deleted": 1}).Error; err != nil {
		log.Error("tvSrv.removeSeason error(%v)", err)
		return
	}
	// manage season's check status
	if err = s.SeasonCheck(season); err != nil {
		return
	}
	// manage season's eps check status and deletion status
	if rows, err = s.DB.Model(&model.Content{}).Where("season_id = ?", season.ID).Select("id, epid, state").Rows(); err != nil {
		log.Error("tvSrv.removeSeason error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		cont := &model.Content{}
		if err = rows.Scan(&cont.ID, &cont.EPID, &cont.State); err != nil {
			log.Error("rows.Scan error(%v)", err)
			continue
		}
		if err = s.EpDel(int64(cont.EPID), 1); err != nil {
			continue
		}
		if err = s.EpCheck(cont); err != nil {
			continue
		}
	}
	return
}

// SnUpdate receives the season update/create info from PGC side and save them in DB
func (s *Service) SnUpdate(c *bm.Context, req *model.TVEpSeason) (err error) {
	var (
		exist     = model.TVEpSeason{}
		updateMap = make(map[string]interface{})
		reqForm   = c.Request.PostForm
	)
	if err = s.DB.Where("id=?", req.ID).First(&exist).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("SnUpdate Sid %d Err %v", req.ID, err)
		return
	}
	if exist.ID <= 0 { // if data not exist, it's brand new data
		req.Check = _seasonDefault
		if err = s.DB.Create(req).Error; err != nil {
			log.Error("tvSrv.createSeason error(%v)", err)
		}
		return
	}
	if exist.IsDeleted == 1 { // exist but was deleted
		if err = s.DB.Model(&model.TVEpSeason{}).Where("id = ?", req.ID).Update(map[string]int{"is_deleted": 0}).Error; err != nil {
			log.Error("tvSrv.removeSeason error(%v)", err)
			return
		}
	}
	if updateMap = exist.Updated(reqForm); len(updateMap) == 0 {
		log.Warn("SnUpdate Sid %d No change", req.ID)
		return
	}
	if err = s.DB.Model(&model.TVEpSeason{}).Where("id = ?", req.ID).Update(updateMap).Error; err != nil {
		log.Error("SnUpdate Sid %d, Update Err %v", req.ID, err)
		return
	}
	log.Info("SnUpdate Sid %d, Update Fields %v Succ", req.ID, updateMap)
	return s.SeasonCheck(&exist)
}

// EpAct acts on ep item
func (s *Service) EpAct(c *bm.Context, cid int64, act int) (err error) {
	var (
		item = &model.Content{}
		db   = s.DB.Model(&model.Content{})
	)
	if err = db.Where("state=?", _epPassed).Where("is_deleted=?", 0).First(item, cid).Error; err != nil {
		log.Error("query fail(%v)\n", err)
		return ecode.RequestErr
	}
	if err = db.Where("id=?", cid).Update(map[string]int{"valid": act}).Error; err != nil {
		log.Error("online error(%v)\n", err)
		return ecode.RequestErr
	}
	return
}
