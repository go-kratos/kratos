package service

import (
	"context"
	"encoding/json"

	"go-common/app/admin/main/esports/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

var (
	_emptyContestList = make([]*model.Contest, 0)
	_emptyContestData = make([]*model.ContestData, 0)
)

const (
	_sortDesc         = 1
	_sortASC          = 2
	_ReplyTypeContest = "27"
)

// ContestInfo .
func (s *Service) ContestInfo(c context.Context, id int64) (data *model.ContestInfo, err error) {
	var (
		gameMap map[int64][]*model.Game
		teamMap map[int64]*model.Team
		teamIDs []int64
		hasTeam bool
	)
	contest := new(model.Contest)
	if err = s.dao.DB.Where("id=?", id).First(&contest).Error; err != nil {
		log.Error("ContestInfo Error (%v)", err)
		return
	}
	if gameMap, err = s.gameList(c, model.TypeContest, []int64{id}); err != nil {
		return
	}
	if contest.HomeID > 0 {
		teamIDs = append(teamIDs, contest.HomeID)
	}
	if contest.AwayID > 0 {
		teamIDs = append(teamIDs, contest.AwayID)
	}
	if ids := unique(teamIDs); len(ids) > 0 {
		var teams []*model.Team
		if err = s.dao.DB.Model(&model.Team{}).Where("id IN (?)", ids).Find(&teams).Error; err != nil {
			log.Error("ContestList team Error (%v)", err)
			return
		}
		if len(teams) > 0 {
			hasTeam = true
		}
		teamMap = make(map[int64]*model.Team, len(teams))
		for _, v := range teams {
			teamMap[v.ID] = v
		}
	}
	data = &model.ContestInfo{Contest: contest}
	if len(gameMap) > 0 {
		if games, ok := gameMap[id]; ok {
			data.Games = games
		}
	}
	if len(data.Games) == 0 {
		data.Games = _emptyGameList
	}
	if hasTeam {
		if team, ok := teamMap[contest.HomeID]; ok {
			data.HomeName = team.Title
		}
		if team, ok := teamMap[contest.AwayID]; ok {
			data.AwayName = team.Title
		}
	}
	var cDatas []*model.ContestData
	if err = s.dao.DB.Model(&model.ContestData{}).Where(map[string]interface{}{"is_deleted": _notDeleted}).Where("cid IN (?)", []int64{id}).Find(&cDatas).Error; err != nil {
		log.Error("ContestInfo Find ContestData Error (%v)", err)
		return
	}
	data.Data = cDatas
	return
}

// ContestList .
func (s *Service) ContestList(c context.Context, mid, sid, pn, ps, srt int64) (list []*model.ContestInfo, count int64, err error) {
	var contests []*model.Contest
	source := s.dao.DB.Model(&model.Contest{})
	if srt == _sortDesc {
		source = source.Order("stime DESC")
	} else if srt == _sortASC {
		source = source.Order("stime ASC")
	}
	if mid > 0 {
		source = source.Where("mid=?", mid)
	}
	if sid > 0 {
		source = source.Where("sid=?", sid)
	}
	source.Count(&count)
	if err = source.Offset((pn - 1) * ps).Limit(ps).Find(&contests).Error; err != nil {
		log.Error("ContestList Error (%v)", err)
		return
	}
	if len(contests) == 0 {
		contests = _emptyContestList
		return
	}
	if list, err = s.contestInfos(c, contests, true); err != nil {
		log.Error("s.contestInfos Error (%v)", err)
	}
	return
}

func (s *Service) contestInfos(c context.Context, contests []*model.Contest, useGame bool) (list []*model.ContestInfo, err error) {
	var (
		conIDs, teamIDs            []int64
		gameMap                    map[int64][]*model.Game
		teamMap                    map[int64]*model.Team
		cDataMap                   map[int64][]*model.ContestData
		hasGame, hasTeam, hasCData bool
	)
	for _, v := range contests {
		conIDs = append(conIDs, v.ID)
		if v.HomeID > 0 {
			teamIDs = append(teamIDs, v.HomeID)
		}
		if v.AwayID > 0 {
			teamIDs = append(teamIDs, v.AwayID)
		}
		if v.SuccessTeam > 0 {
			teamIDs = append(teamIDs, v.SuccessTeam)
		}
	}
	if useGame {
		if gameMap, err = s.gameList(c, model.TypeContest, conIDs); err != nil {
			return
		} else if len(gameMap) > 0 {
			hasGame = true
		}
	}
	if ids := unique(teamIDs); len(ids) > 0 {
		var teams []*model.Team
		if err = s.dao.DB.Model(&model.Team{}).Where("id IN (?)", ids).Find(&teams).Error; err != nil {
			log.Error("ContestList team Error (%v)", err)
			return
		}
		if len(teams) > 0 {
			hasTeam = true
		}
		teamMap = make(map[int64]*model.Team, len(teams))
		for _, v := range teams {
			teamMap[v.ID] = v
		}
	}
	if len(conIDs) > 0 {
		var cDatas []*model.ContestData
		if err = s.dao.DB.Model(&model.ContestData{}).Where(map[string]interface{}{"is_deleted": _notDeleted}).Where("cid IN (?)", conIDs).Find(&cDatas).Error; err != nil {
			log.Error("ContestList Find ContestData Error (%v)", err)
			return
		}
		if len(cDatas) > 0 {
			hasCData = true
		}
		cDataMap = make(map[int64][]*model.ContestData, len(cDatas))
		for _, v := range cDatas {
			cDataMap[v.CID] = append(cDataMap[v.CID], v)
		}
	}
	for _, v := range contests {
		contest := &model.ContestInfo{Contest: v}
		if hasGame {
			if games, ok := gameMap[v.ID]; ok {
				contest.Games = games
			}
		}
		if len(contest.Games) == 0 {
			contest.Games = _emptyGameList
		}
		if hasTeam {
			if team, ok := teamMap[v.HomeID]; ok {
				contest.HomeName = team.Title
			}
			if team, ok := teamMap[v.AwayID]; ok {
				contest.AwayName = team.Title
			}
			if team, ok := teamMap[v.SuccessTeam]; ok {
				contest.SuccessName = team.Title
			}
		}
		if hasCData {
			if cData, ok := cDataMap[v.ID]; ok {
				contest.Data = cData
			}
		} else {
			contest.Data = _emptyContestData
		}
		list = append(list, contest)
	}
	return
}

// AddContest .
func (s *Service) AddContest(c context.Context, param *model.Contest, gids []int64) (err error) {
	// TODO check name exist
	// check sid
	season := new(model.Season)
	if err = s.dao.DB.Where("id=?", param.Sid).Where("status=?", _statusOn).First(&season).Error; err != nil {
		log.Error("AddContest s.dao.DB.Where id(%d) error(%d)", param.Sid, err)
		return
	}
	// check mid
	match := new(model.Match)
	if err = s.dao.DB.Where("id=?", param.Mid).Where("status=?", _statusOn).First(&match).Error; err != nil {
		log.Error("AddContest s.dao.DB.Where id(%d) error(%d)", param.Mid, err)
		return
	}
	// check game idsEsportsCDataErr
	var (
		games       []*model.Game
		gidMaps     []*model.GIDMap
		contestData []*model.ContestData
	)
	if param.DataType == 0 {
		param.Data = ""
		param.MatchID = 0
	}
	if param.Data != "" {
		if err = json.Unmarshal([]byte(param.Data), &contestData); err != nil {
			err = ecode.EsportsContestDataErr
			return
		}
	}
	if err = s.dao.DB.Model(&model.Game{}).Where("status=?", _statusOn).Where("id IN (?)", gids).Find(&games).Error; err != nil {
		log.Error("AddContest check game ids Error (%v)", err)
		return
	}
	if len(games) == 0 {
		log.Error("AddContest games(%v) not found", gids)
		err = ecode.RequestErr
		return
	}
	tx := s.dao.DB.Begin()
	if err = tx.Model(&model.Contest{}).Create(param).Error; err != nil {
		log.Error("AddContest tx.Model Create(%+v) error(%v)", param, err)
		err = tx.Rollback().Error
		return
	}
	for _, v := range games {
		gidMaps = append(gidMaps, &model.GIDMap{Type: model.TypeContest, Oid: param.ID, Gid: v.ID})
	}
	if err = tx.Model(&model.GIDMap{}).Exec(model.GidBatchAddSQL(gidMaps)).Error; err != nil {
		log.Error("AddContest tx.Model Create(%+v) error(%v)", param, err)
		err = tx.Rollback().Error
		return
	}
	if len(contestData) > 0 {
		if err = tx.Model(&model.Module{}).Exec(model.BatchAddCDataSQL(param.ID, contestData)).Error; err != nil {
			log.Error("AddContest Module tx.Model Create(%+v) error(%v)", param, err)
			err = tx.Rollback().Error
			return
		}
	}
	if err = tx.Commit().Error; err != nil {
		return
	}
	// register reply
	if err = s.dao.RegReply(c, param.ID, param.Adid, _ReplyTypeContest); err != nil {
		err = nil
	}
	return
}

// EditContest .
func (s *Service) EditContest(c context.Context, param *model.Contest, gids []int64) (err error) {
	var (
		games                    []*model.Game
		preGidMaps, addGidMaps   []*model.GIDMap
		upGidMapAdd, upGidMapDel []int64
	)
	// TODO check name exist
	// check sid
	season := new(model.Season)
	if err = s.dao.DB.Where("id=?", param.Sid).Where("status=?", _statusOn).First(&season).Error; err != nil {
		log.Error("EditContest s.dao.DB.Where id(%d) error(%d)", param.Sid, err)
		return
	}
	// check mid
	match := new(model.Match)
	if err = s.dao.DB.Where("id=?", param.Mid).Where("status=?", _statusOn).First(&match).Error; err != nil {
		log.Error("EditContest s.dao.DB.Where id(%d) error(%d)", param.Mid, err)
		return
	}
	preData := new(model.Contest)
	if err = s.dao.DB.Where("id=?", param.ID).First(&preData).Error; err != nil {
		log.Error("EditContest s.dao.DB.Where id(%d) error(%d)", param.ID, err)
		return
	}
	if err = s.dao.DB.Model(&model.Game{}).Where("status=?", _statusOn).Where("id IN (?)", gids).Find(&games).Error; err != nil {
		log.Error("EditContest check game ids Error (%v)", err)
		return
	}
	if len(games) == 0 {
		log.Error("EditContest games(%v) not found", gids)
		err = ecode.RequestErr
		return
	}
	if err = s.dao.DB.Model(&model.GIDMap{}).Where("oid=?", param.ID).Where("type=?", model.TypeContest).Find(&preGidMaps).Error; err != nil {
		log.Error("EditContest games(%v) not found", gids)
		return
	}
	var (
		newCData []*model.ContestData
	)
	if param.DataType == 0 {
		param.Data = ""
		param.MatchID = 0
	}
	if param.Data != "" {
		if err = json.Unmarshal([]byte(param.Data), &newCData); err != nil {
			err = ecode.EsportsContestDataErr
			return
		}
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
			addGidMaps = append(addGidMaps, &model.GIDMap{Type: model.TypeContest, Oid: param.ID, Gid: gid})
		}
	}
	tx := s.dao.DB.Begin()
	if err = tx.Error; err != nil {
		log.Error("s.dao.DB.Begin error(%v)", err)
		return
	}
	if err = tx.Model(&model.Contest{}).Save(param).Error; err != nil {
		log.Error("EditContest Update(%+v) error(%v)", param, err)
		err = tx.Rollback().Error
		return
	}
	if len(upGidMapAdd) > 0 {
		if err = tx.Model(&model.GIDMap{}).Where("id IN (?)", upGidMapAdd).Updates(map[string]interface{}{"is_deleted": _notDeleted}).Error; err != nil {
			log.Error("EditContest GIDMap Add(%+v) error(%v)", upGidMapAdd, err)
			err = tx.Rollback().Error
			return
		}
	}
	if len(upGidMapDel) > 0 {
		if err = tx.Model(&model.GIDMap{}).Where("id IN (?)", upGidMapDel).Updates(map[string]interface{}{"is_deleted": _deleted}).Error; err != nil {
			log.Error("EditContest GIDMap Del(%+v) error(%v)", upGidMapDel, err)
			err = tx.Rollback().Error
			return
		}
	}
	if len(addGidMaps) > 0 {
		if err = tx.Model(&model.GIDMap{}).Exec(model.GidBatchAddSQL(addGidMaps)).Error; err != nil {
			log.Error("EditContest GIDMap Create(%+v) error(%v)", addGidMaps, err)
			err = tx.Rollback().Error
			return
		}
	}
	var (
		mapOldCData, mapNewCData    map[int64]*model.ContestData
		upCData, addCData, oldCData []*model.ContestData
		delCData                    []int64
	)
	if len(newCData) > 0 {
		// check module
		if err = s.dao.DB.Model(&model.ContestData{}).Where("cid=?", param.ID).Where("is_deleted=?", _notDeleted).Find(&oldCData).Error; err != nil {
			log.Error("EditContest s.dao.DB.Model Find (%+v) error(%v)", param.ID, err)
			return
		}
		mapOldCData = make(map[int64]*model.ContestData, len(oldCData))
		for _, v := range oldCData {
			mapOldCData[v.ID] = v
		}
		//新数据在老数据中 更新老数据。新的数据不在老数据 添加新数据
		for _, cData := range newCData {
			if _, ok := mapOldCData[cData.ID]; ok {
				upCData = append(upCData, cData)
			} else {
				addCData = append(addCData, cData)
			}
		}
		mapNewCData = make(map[int64]*model.ContestData, len(oldCData))
		for _, v := range newCData {
			mapNewCData[v.ID] = v
		}
		//老数据在新中 上面已经处理。老数据不在新数据中 删除老数据
		for _, cData := range oldCData {
			if _, ok := mapNewCData[cData.ID]; !ok {
				delCData = append(delCData, cData.ID)
			}
		}
		if len(upCData) > 0 {
			if err = tx.Model(&model.ContestData{}).Exec(model.BatchEditCDataSQL(upCData)).Error; err != nil {
				log.Error("EditContest s.dao.DB.Model tx.Model Exec(%+v) error(%v)", upCData, err)
				err = tx.Rollback().Error
				return
			}
		}
		if len(delCData) > 0 {
			if err = tx.Model(&model.ContestData{}).Where("id IN (?)", delCData).Updates(map[string]interface{}{"is_deleted": _deleted}).Error; err != nil {
				log.Error("EditContest s.dao.DB.Model Updates(%+v) error(%v)", delCData, err)
				err = tx.Rollback().Error
				return
			}
		}
		if len(addCData) > 0 {
			if err = tx.Model(&model.ContestData{}).Exec(model.BatchAddCDataSQL(param.ID, addCData)).Error; err != nil {
				log.Error("EditContest s.dao.DB.Model Create(%+v) error(%v)", addCData, err)
				err = tx.Rollback().Error
				return
			}
		}
	} else {
		if err = tx.Model(&model.ContestData{}).Where("cid = ?", param.ID).Updates(map[string]interface{}{"is_deleted": _deleted}).Error; err != nil {
			log.Error("EditContest s.dao.DB.Model Updates(%+v) error(%v)", param.ID, err)
			err = tx.Rollback().Error
			return
		}
	}
	err = tx.Commit().Error
	return
}

// ForbidContest .
func (s *Service) ForbidContest(c context.Context, id int64, state int) (err error) {
	preContest := new(model.Contest)
	if err = s.dao.DB.Where("id=?", id).First(&preContest).Error; err != nil {
		log.Error("ContestForbid s.dao.DB.Where id(%d) error(%d)", id, err)
		return
	}
	if err = s.dao.DB.Model(&model.Contest{}).Where("id=?", id).Update(map[string]int{"status": state}).Error; err != nil {
		log.Error("ContestForbid s.dao.DB.Model error(%v)", err)
	}
	return
}
