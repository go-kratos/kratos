package service

import (
	"context"

	"go-common/app/admin/main/esports/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// AddDetail .
func (s *Service) AddDetail(c context.Context, param *model.MatchDetail) (err error) {
	if err = s.dao.DB.Model(&model.MatchDetail{}).Create(param).Error; err != nil {
		log.Error("AddDetail MatchDetail db.Model Create error(%v)", err)
	}
	return
}

// EditDetail .
func (s *Service) EditDetail(c context.Context, param *model.MatchDetail) (err error) {
	upFields := map[string]interface{}{"ma_id": param.MaID, "game_type": param.GameType, "stime": param.Stime,
		"etime": param.Etime, "game_stage": param.GameStage, "knockout_type": param.KnockoutType,
		"winner_type": param.WinnerType, "ScoreID": param.ScoreID, "status": param.Status, "online": param.Online}
	if err = s.dao.DB.Model(&model.MatchDetail{}).Where("id=?", param.ID).Update(upFields).Error; err != nil {
		log.Error("EditDetail MatchDetail db.Model Update error(%v)", err)
	}
	return
}

// ForbidDetail .
func (s *Service) ForbidDetail(c context.Context, id int64, state int) (err error) {
	if err = s.dao.DB.Model(&model.MatchDetail{}).Where("id=?", id).Updates(map[string]interface{}{"status": state}).Error; err != nil {
		log.Error("ForbidDetail MatchDetail db.Model Updates(%d) error(%v)", id, err)
	}
	return
}

// UpOnline .
func (s *Service) UpOnline(c context.Context, id int64, onLine int64) (err error) {
	if onLine == _online {
		var count int64
		treeDB := s.dao.DB.Model(&model.Tree{}).Where("mad_id=?", id).Where("is_deleted=0")
		if err = treeDB.Error; err != nil {
			log.Error("upOnline  treeDB  Error (%v)", err)
			return
		}
		treeDB.Count(&count)
		if count == 0 {
			err = ecode.EsportsTreeEmptyErr
			return
		}
	}
	if err = s.dao.DB.Model(&model.MatchDetail{}).Where("id=?", id).Updates(map[string]interface{}{"online": onLine}).Error; err != nil {
		log.Error("UpOnline s.dao.DB.Model  Updates(%+v) error(%v)", id, err)
	}
	return
}

// ListDetail .
func (s *Service) ListDetail(c context.Context, maID, pn, ps int64) (rs []*model.MatchDetail, count int64, err error) {
	db := s.dao.DB.Model(&model.MatchDetail{}).Offset((pn-1)*ps).Where("ma_id=?", maID).Order("id ASC").Limit(ps).Find(&rs)
	if err = db.Error; err != nil {
		log.Error("ListDetail MatchDetail db.Model Find error(%v)", err)
	}
	db.Count(&count)
	return
}
