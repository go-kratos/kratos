package service

import (
	"context"
	"fmt"

	"go-common/app/admin/main/esports/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

var _emptyMatchList = make([]*model.Match, 0)

// MatchInfo .
func (s *Service) MatchInfo(c context.Context, id int64) (data *model.MatchInfo, err error) {
	var gameMap map[int64][]*model.Game
	match := new(model.Match)
	if err = s.dao.DB.Where("id=?", id).First(&match).Error; err != nil {
		log.Error("MatchInfo Error (%v)", err)
		return
	}
	if gameMap, err = s.gameList(c, model.TypeMatch, []int64{id}); err != nil {
		return
	}
	if games, ok := gameMap[id]; ok {
		data = &model.MatchInfo{Match: match, Games: games}
	} else {
		data = &model.MatchInfo{Match: match, Games: _emptyGameList}
	}
	return
}

// MatchList .
func (s *Service) MatchList(c context.Context, pn, ps int64, title string) (list []*model.MatchInfo, count int64, err error) {
	var (
		matchs   []*model.Match
		matchIDs []int64
		gameMap  map[int64][]*model.Game
	)
	source := s.dao.DB.Model(&model.Match{})
	if title != "" {
		likeStr := fmt.Sprintf("%%%s%%", title)
		source = source.Where("title like ?", likeStr)
	}
	source.Count(&count)
	if err = source.Offset((pn - 1) * ps).Order("rank DESC,id ASC").Limit(ps).Find(&matchs).Error; err != nil {
		log.Error("MatchList Error (%v)", err)
		return
	}
	if len(matchs) == 0 {
		return
	}
	for _, v := range matchs {
		matchIDs = append(matchIDs, v.ID)
	}
	if gameMap, err = s.gameList(c, model.TypeMatch, matchIDs); err != nil {
		return
	}
	for _, v := range matchs {
		if games, ok := gameMap[v.ID]; ok {
			list = append(list, &model.MatchInfo{Match: v, Games: games})
		} else {
			list = append(list, &model.MatchInfo{Match: v, Games: _emptyGameList})
		}
	}
	return
}

// AddMatch .
func (s *Service) AddMatch(c context.Context, param *model.Match, gids []int64) (err error) {
	// check game ids
	var (
		games   []*model.Game
		gidMaps []*model.GIDMap
	)
	if err = s.dao.DB.Model(&model.Game{}).Where("status=?", _statusOn).Where("id IN (?)", gids).Find(&games).Error; err != nil {
		log.Error("AddMatch check game ids Error (%v)", err)
		return
	}
	if len(games) == 0 {
		log.Error("AddMatch games(%v) not found", gids)
		err = ecode.RequestErr
		return
	}
	tx := s.dao.DB.Begin()
	if err = tx.Error; err != nil {
		log.Error("s.dao.DB.Begin error(%v)", err)
		return
	}
	if err = tx.Model(&model.Match{}).Create(param).Error; err != nil {
		log.Error("AddMatch s.dao.DB.Model Create(%+v) error(%v)", param, err)
		err = tx.Rollback().Error
		return
	}
	for _, v := range games {
		gidMaps = append(gidMaps, &model.GIDMap{Type: model.TypeMatch, Oid: param.ID, Gid: v.ID})
	}
	if err = tx.Model(&model.GIDMap{}).Exec(model.GidBatchAddSQL(gidMaps)).Error; err != nil {
		log.Error("AddMatch s.dao.DB.Model Create(%+v) error(%v)", param, err)
		err = tx.Rollback().Error
		return
	}
	err = tx.Commit().Error
	return
}

// EditMatch .
func (s *Service) EditMatch(c context.Context, param *model.Match, gids []int64) (err error) {
	var (
		games                    []*model.Game
		preGidMaps, addGidMaps   []*model.GIDMap
		upGidMapAdd, upGidMapDel []int64
	)
	preData := new(model.Match)
	if err = s.dao.DB.Where("id=?", param.ID).First(&preData).Error; err != nil {
		log.Error("EditMatch s.dao.DB.Where id(%d) error(%d)", param.ID, err)
		return
	}
	if err = s.dao.DB.Model(&model.Game{}).Where("status=?", _statusOn).Where("id IN (?)", gids).Find(&games).Error; err != nil {
		log.Error("AddMatch check game ids Error (%v)", err)
		return
	}
	if len(games) == 0 {
		log.Error("AddMatch games(%v) not found", gids)
		err = ecode.RequestErr
		return
	}
	if err = s.dao.DB.Model(&model.GIDMap{}).Where("oid=?", param.ID).Where("type=?", model.TypeMatch).Find(&preGidMaps).Error; err != nil {
		log.Error("AddMatch games(%v) not found", gids)
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
			addGidMaps = append(addGidMaps, &model.GIDMap{Type: model.TypeMatch, Oid: param.ID, Gid: gid})
		}
	}
	tx := s.dao.DB.Begin()
	if err = tx.Error; err != nil {
		log.Error("s.dao.DB.Begin error(%v)", err)
		return
	}
	if err = tx.Model(&model.Match{}).Save(param).Error; err != nil {
		log.Error("EditMatch Match Update(%+v) error(%v)", param, err)
		err = tx.Rollback().Error
		return
	}
	if len(upGidMapAdd) > 0 {
		if err = tx.Model(&model.GIDMap{}).Where("id IN (?)", upGidMapAdd).Updates(map[string]interface{}{"is_deleted": _notDeleted}).Error; err != nil {
			log.Error("EditMatch GIDMap Add(%+v) error(%v)", upGidMapAdd, err)
			err = tx.Rollback().Error
			return
		}
	}
	if len(upGidMapDel) > 0 {
		if err = tx.Model(&model.GIDMap{}).Where("id IN (?)", upGidMapDel).Updates(map[string]interface{}{"is_deleted": _deleted}).Error; err != nil {
			log.Error("EditMatch GIDMap Del(%+v) error(%v)", upGidMapDel, err)
			err = tx.Rollback().Error
			return
		}
	}
	if len(addGidMaps) > 0 {
		if err = tx.Model(&model.GIDMap{}).Exec(model.GidBatchAddSQL(addGidMaps)).Error; err != nil {
			log.Error("EditMatch GIDMap Create(%+v) error(%v)", addGidMaps, err)
			err = tx.Rollback().Error
			return
		}
	}
	err = tx.Commit().Error
	return
}

// ForbidMatch .
func (s *Service) ForbidMatch(c context.Context, id int64, state int) (err error) {
	preMatch := new(model.Match)
	if err = s.dao.DB.Where("id=?", id).First(&preMatch).Error; err != nil {
		log.Error("MatchForbid s.dao.DB.Where id(%d) error(%d)", id, err)
		return
	}
	if err = s.dao.DB.Model(&model.Match{}).Where("id=?", id).Update(map[string]int{"status": state}).Error; err != nil {
		log.Error("MatchForbid s.dao.DB.Model error(%v)", err)
	}
	return
}
