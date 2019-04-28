package service

import (
	"context"

	"go-common/app/admin/main/esports/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// SeasonInfo .
func (s *Service) SeasonInfo(c context.Context, id int64) (data *model.SeasonInfo, err error) {
	var gameMap map[int64][]*model.Game
	season := new(model.Season)
	if err = s.dao.DB.Where("id=?", id).First(&season).Error; err != nil {
		log.Error("SeasonInfo Error (%v)", err)
		return
	}
	if gameMap, err = s.gameList(c, model.TypeSeason, []int64{id}); err != nil {
		return
	}
	if games, ok := gameMap[id]; ok {
		data = &model.SeasonInfo{Season: season, Games: games}
	} else {
		data = &model.SeasonInfo{Season: season, Games: _emptyGameList}
	}
	return
}

// SeasonList .
func (s *Service) SeasonList(c context.Context, mid, pn, ps int64) (list []*model.SeasonInfo, count int64, err error) {
	var (
		seasons   []*model.Season
		seasonIDs []int64
		gameMap   map[int64][]*model.Game
	)
	source := s.dao.DB.Model(&model.Season{})
	if mid > 0 {
		source = source.Where("mid=?", mid)
	}
	source.Count(&count)
	if err = source.Offset((pn - 1) * ps).Order("rank DESC,id ASC").Limit(ps).Find(&seasons).Error; err != nil {
		log.Error("SeasonList Error (%v)", err)
		return
	}
	if len(seasons) == 0 {
		return
	}
	for _, v := range seasons {
		seasonIDs = append(seasonIDs, v.ID)
	}
	if gameMap, err = s.gameList(c, model.TypeSeason, seasonIDs); err != nil {
		return
	}
	for _, v := range seasons {
		if games, ok := gameMap[v.ID]; ok {
			list = append(list, &model.SeasonInfo{Season: v, Games: games})
		} else {
			list = append(list, &model.SeasonInfo{Season: v, Games: _emptyGameList})
		}
	}
	return
}

// AddSeason .
func (s *Service) AddSeason(c context.Context, param *model.Season, gids []int64) (err error) {
	var (
		games   []*model.Game
		gidMaps []*model.GIDMap
	)
	// TODO check name exist
	if err = s.dao.DB.Model(&model.Game{}).Where("status=?", _statusOn).Where("id IN (?)", gids).Find(&games).Error; err != nil {
		log.Error("AddSeason check game ids Error (%v)", err)
		return
	}
	if len(games) == 0 {
		log.Error("AddSeason games(%v) not found", gids)
		err = ecode.RequestErr
		return
	}
	tx := s.dao.DB.Begin()
	if err = tx.Error; err != nil {
		log.Error("s.dao.DB.Begin error(%v)", err)
		return
	}
	if err = tx.Model(&model.Season{}).Create(param).Error; err != nil {
		log.Error("AddSeason s.dao.DB.Model Create(%+v) error(%v)", param, err)
		err = tx.Rollback().Error
		return
	}
	for _, v := range games {
		gidMaps = append(gidMaps, &model.GIDMap{Type: model.TypeSeason, Oid: param.ID, Gid: v.ID})
	}
	if err = tx.Model(&model.GIDMap{}).Exec(model.GidBatchAddSQL(gidMaps)).Error; err != nil {
		log.Error("AddSeason s.dao.DB.Model Create(%+v) error(%v)", param, err)
		err = tx.Rollback().Error
		return
	}
	err = tx.Commit().Error
	return
}

// EditSeason .
func (s *Service) EditSeason(c context.Context, param *model.Season, gids []int64) (err error) {
	var (
		games                    []*model.Game
		preGidMaps, addGidMaps   []*model.GIDMap
		upGidMapAdd, upGidMapDel []int64
	)
	preData := new(model.Season)
	if err = s.dao.DB.Where("id=?", param.ID).First(&preData).Error; err != nil {
		log.Error("EditSeason s.dao.DB.Where id(%d) error(%d)", param.ID, err)
		return
	}
	if err = s.dao.DB.Model(&model.Game{}).Where("status=?", _statusOn).Where("id IN (?)", gids).Find(&games).Error; err != nil {
		log.Error("EditSeason check game ids Error (%v)", err)
		return
	}
	if len(games) == 0 {
		log.Error("EditSeason games(%v) not found", gids)
		err = ecode.RequestErr
		return
	}
	if err = s.dao.DB.Model(&model.GIDMap{}).Where("oid=?", param.ID).Where("type=?", model.TypeSeason).Find(&preGidMaps).Error; err != nil {
		log.Error("EditSeason games(%v) not found", gids)
		return
	}
	gidsMap := make(map[int64]int64, len(gids))
	preGidsMap := make(map[int64]int64, len(preGidMaps))
	for _, v := range gids {
		gidsMap[v] = v
	}
	for _, v := range preGidMaps {
		preGidsMap[v.Gid] = v.Gid
		if _, ok := gidsMap[v.Gid]; ok {
			if v.IsDeleted == 1 {
				upGidMapAdd = append(upGidMapAdd, v.ID)
			}
		} else {
			upGidMapDel = append(upGidMapDel, v.ID)
		}
	}
	for _, gid := range gids {
		if _, ok := preGidsMap[gid]; !ok {
			addGidMaps = append(addGidMaps, &model.GIDMap{Type: model.TypeSeason, Oid: param.ID, Gid: gid})
		}
	}
	tx := s.dao.DB.Begin()
	if err = tx.Error; err != nil {
		log.Error("s.dao.DB.Begin error(%v)", err)
		return
	}
	if err = tx.Model(&model.Season{}).Save(param).Error; err != nil {
		log.Error("EditSeason Match Update(%+v) error(%v)", param, err)
		err = tx.Rollback().Error
		return
	}
	if len(upGidMapAdd) > 0 {
		if err = tx.Model(&model.GIDMap{}).Where("id IN (?)", upGidMapAdd).Updates(map[string]interface{}{"is_deleted": _notDeleted}).Error; err != nil {
			log.Error("EditSeason GIDMap Add(%+v) error(%v)", upGidMapAdd, err)
			err = tx.Rollback().Error
			return
		}
	}
	if len(upGidMapDel) > 0 {
		if err = tx.Model(&model.GIDMap{}).Where("id IN (?)", upGidMapDel).Updates(map[string]interface{}{"is_deleted": _deleted}).Error; err != nil {
			log.Error("EditSeason GIDMap Del(%+v) error(%v)", upGidMapDel, err)
			err = tx.Rollback().Error
			return
		}
	}
	if len(addGidMaps) > 0 {
		if err = tx.Model(&model.GIDMap{}).Exec(model.GidBatchAddSQL(addGidMaps)).Error; err != nil {
			log.Error("EditSeason GIDMap Create(%+v) error(%v)", addGidMaps, err)
			err = tx.Rollback().Error
			return
		}
	}
	err = tx.Commit().Error
	return
}

// ForbidSeason .
func (s *Service) ForbidSeason(c context.Context, id int64, state int) (err error) {
	preSeason := new(model.Season)
	if err = s.dao.DB.Where("id=?", id).First(&preSeason).Error; err != nil {
		log.Error("SeasonForbid s.dao.DB.Where id(%d) error(%d)", id, err)
		return
	}
	if err = s.dao.DB.Model(&model.Season{}).Where("id=?", id).Update(map[string]int{"status": state}).Error; err != nil {
		log.Error("SeasonForbid s.dao.DB.Model error(%v)", err)
	}
	return
}
