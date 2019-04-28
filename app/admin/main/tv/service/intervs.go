package service

import (
	"go-common/app/admin/main/tv/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

// 0=not found, 1=pgc, 2=cms, 3=license
const (
	ErrNotFound  = 0
	_TypeDefault = 0
	_TypePGC     = 1
	_TypeUGC     = 2
)

func (s *Service) getSeason(sid int64) (res *model.TVEpSeason, err error) {
	sn := model.TVEpSeason{}
	if err = s.DB.Where("id = ?", sid).First(&sn).Error; err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Error("GetSeason error(%v)\n", err)
		return
	}
	return &sn, nil
}

// RemoveInvalids removes invalid interventions
func (s *Service) RemoveInvalids(invalids []*model.RankError) (err error) {
	tx := s.DB.Begin()
	for _, v := range invalids {
		if err = tx.Model(&model.Rank{}).Where("id=?", v.ID).Update(map[string]int{"is_deleted": 1}).Error; err != nil {
			log.Error("tvSrv.RemoveInvalids error(%v)", err)
			tx.Rollback()
			return
		}
	}
	tx.Commit()
	log.Info("Remove Invalid Interventions: %d", len(invalids))
	return
}

// Intervs pick the intervention and combine the season data
func (s *Service) Intervs(req *model.IntervListReq) (res *model.RankList, err error) {
	var (
		intervs  []*model.Rank
		items    []*model.SimpleRank
		invalids []*model.RankError
		db       = req.BuildDB(s.DB).Order("position asc")
	)
	if err = db.Find(&intervs).Error; err != nil {
		log.Error("[Intervs] DB query fail(%v)", err)
		return
	}
	items, invalids = s.intervsValid(intervs)
	err = s.RemoveInvalids(invalids)
	res = &model.RankList{
		List: items,
	}
	return
}

func (s *Service) intervsValid(intervs []*model.Rank) (items []*model.SimpleRank, invalids []*model.RankError) {
	// check Its Season Status, pick invalid ones
	for _, v := range intervs {
		switch v.ContType {
		case _TypePGC, _TypeDefault:
			isValid, sn := s.snValid(v.ContID)
			if !isValid {
				invalids = append(invalids, v.BeError())
				continue
			}
			items = append(items, v.BeSimpleSn(sn, s.pgcCatToName))
		case _TypeUGC:
			isValid, arc := s.arcValid(v.ContID)
			if !isValid {
				invalids = append(invalids, v.BeError())
				continue
			}
			items = append(items, v.BeSimpleArc(arc, s.arcPName))
		default:
			log.Error("[Intervs] Rank Error Cont_Type %d, RankID:%d", v.ContType, v.ID)
			continue
		}
	}
	return
}

// snValid Distinguish whether the Season is existing and valid
func (s *Service) snValid(sid int64) (res bool, season *model.TVEpSeason) {
	var err error
	if season, err = s.getSeason(sid); err != nil || season == nil {
		return
	}
	res = errTyping(int(season.Check), int(season.Valid), int(season.IsDeleted))
	return
}

// arcValid Distinguish whether the archive is existing and valid
func (s *Service) arcValid(aid int64) (res bool, arc *model.SimpleArc) {
	var err error
	if arc, err = s.ExistArc(aid); err != nil || arc == nil {
		return
	}
	res = errTyping(arc.Result, arc.Valid, arc.Deleted)
	return
}

func errTyping(check, valid, isDeleted int) (res bool) {
	if check == 1 && valid == 1 && isDeleted == 0 {
		return true
	}
	return false
}

// RefreshIntervs is used to delete the previous interventions
func (s *Service) RefreshIntervs(req *model.IntervPubReq) (invalid *model.RankError, err error) {
	var (
		tx       = s.DB.Begin()
		txDel    = req.BuildDB(tx)
		position = 1
		title    string
	)
	if err = txDel.Delete(&model.Rank{}).Error; err != nil { // delete old intervs
		log.Error("Del Previsou Intervs error(%v)\n", err)
		tx.Rollback()
		return
	}
	for _, v := range req.Items {
		if invalid, title = s.checkInterv(req, v); invalid != nil {
			tx.Rollback()
			return
		}
		if err = tx.Create(v.BeComplete(req, title, position)).Error; err != nil { // create new ones
			log.Error("Create New Intervs %v ,Error(%v)\n", v, err)
			tx.Rollback()
			return
		}
		position = position + 1
	}
	tx.Commit()
	log.Info("RefreshIntervs Success")
	return
}

// checkInterv checks whether the to-publish intervention is valid
func (s *Service) checkInterv(req *model.IntervPubReq, v *model.SimpleRank) (invalid *model.RankError, title string) {
	var (
		isValid bool
		sn      *model.TVEpSeason
		arc     *model.SimpleArc
		rankErr = &model.RankError{
			ID:       int(v.ID),
			SeasonID: int(v.ContID),
		}
	)
	if req.IsIdx() {
		if int(req.Category) != v.ContType+model.RankIdxBase { // if ugc, we can't accept pgc data
			isValid = false
			return rankErr, ""
		}
	}
	switch v.ContType {
	case _TypePGC, _TypeDefault:
		isValid, sn = s.snValid(v.ContID)
		if isValid && req.IsIdx() { // if index, check it's the pgc category's season
			isValid = (sn.Category == int(req.Rank))
		}
	case _TypeUGC:
		isValid, arc = s.arcValid(v.ContID)
		if isValid && req.IsIdx() { // if index, check it's the first level type's archive
			pid := s.GetArchivePid(arc.TypeID)
			isValid = pid == int32(req.Rank)
		}
	default:
		log.Error("[Intervs] Rank Error Cont_Type %d, RankID:%d", v.ContType, v.ID)
		return rankErr, ""
	}
	if !isValid {
		log.Error("snValid (%d) Not passed", v.ContID)
		return rankErr, ""
	}
	if v.ContType == _TypeUGC && arc != nil {
		title = arc.Title
	}
	if (v.ContType == _TypePGC || v.ContType == _TypeDefault) && sn != nil {
		title = sn.Title
	}
	return
}

func (s *Service) pgcCatToName(cat int) (res string) {
	if res, ok := s.pgcCatName[cat]; ok {
		return res
	}
	return ""
}
