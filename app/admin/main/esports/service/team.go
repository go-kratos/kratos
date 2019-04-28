package service

import (
	"context"
	"fmt"

	"go-common/app/admin/main/esports/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

var _emptyTeamList = make([]*model.Team, 0)

// TeamInfo .
func (s *Service) TeamInfo(c context.Context, id int64) (data *model.TeamInfo, err error) {
	var gameMap map[int64][]*model.Game
	team := new(model.Team)
	if err = s.dao.DB.Where("id=?", id).Where("is_deleted=?", _notDeleted).First(&team).Error; err != nil {
		log.Error("TeamInfo Error (%v)", err)
		return
	}
	if gameMap, err = s.gameList(c, model.TypeTeam, []int64{id}); err != nil {
		return
	}
	if games, ok := gameMap[id]; ok {
		data = &model.TeamInfo{Team: team, Games: games}
	} else {
		data = &model.TeamInfo{Team: team, Games: _emptyGameList}
	}
	return
}

// TeamList .
func (s *Service) TeamList(c context.Context, pn, ps int64, title string, status int) (list []*model.TeamInfo, count int64, err error) {
	var (
		teams   []*model.Team
		teamIDs []int64
		gameMap map[int64][]*model.Game
	)
	source := s.dao.DB.Model(&model.Team{})
	if status != _statusAll {
		source = source.Where("is_deleted=?", _notDeleted)
	}
	if title != "" {
		likeStr := fmt.Sprintf("%%%s%%", title)
		source = source.Where("title like ?", likeStr)
	}
	source.Count(&count)
	if err = source.Offset((pn - 1) * ps).Limit(ps).Find(&teams).Error; err != nil {
		log.Error("TeamList Error (%v)", err)
		return
	}
	if len(teams) == 0 {
		return
	}
	for _, v := range teams {
		teamIDs = append(teamIDs, v.ID)
	}
	if gameMap, err = s.gameList(c, model.TypeTeam, teamIDs); err != nil {
		return
	}
	for _, v := range teams {
		if games, ok := gameMap[v.ID]; ok {
			list = append(list, &model.TeamInfo{Team: v, Games: games})
		} else {
			list = append(list, &model.TeamInfo{Team: v, Games: _emptyGameList})
		}
	}
	return
}

// AddTeam .
func (s *Service) AddTeam(c context.Context, param *model.Team, gids []int64) (err error) {
	var (
		games   []*model.Game
		gidMaps []*model.GIDMap
	)
	if err = s.dao.DB.Model(&model.Game{}).Where("status=?", _statusOn).Where("id IN (?)", gids).Find(&games).Error; err != nil {
		log.Error("AddTeam check game ids Error (%v)", err)
		return
	}
	if len(games) == 0 {
		log.Error("AddTeam games(%v) not found", gids)
		err = ecode.RequestErr
		return
	}
	tx := s.dao.DB.Begin()
	if err = tx.Error; err != nil {
		log.Error("s.dao.DB.Begin error(%v)", err)
		return
	}
	if err = tx.Model(&model.Team{}).Create(param).Error; err != nil {
		log.Error("AddTeam s.dao.DB.Model Create(%+v) error(%v)", param, err)
		err = tx.Rollback().Error
		return
	}
	for _, v := range games {
		gidMaps = append(gidMaps, &model.GIDMap{Type: model.TypeTeam, Oid: param.ID, Gid: v.ID})
	}
	if err = tx.Model(&model.GIDMap{}).Exec(model.GidBatchAddSQL(gidMaps)).Error; err != nil {
		log.Error("AddTeam s.dao.DB.Model Create(%+v) error(%v)", param, err)
		err = tx.Rollback().Error
		return
	}
	err = tx.Commit().Error
	return
}

// EditTeam .
func (s *Service) EditTeam(c context.Context, param *model.Team, gids []int64) (err error) {
	var (
		games                    []*model.Game
		preGidMaps, addGidMaps   []*model.GIDMap
		upGidMapAdd, upGidMapDel []int64
	)
	preData := new(model.Team)
	if err = s.dao.DB.Where("id=?", param.ID).First(&preData).Error; err != nil {
		log.Error("EditTeam s.dao.DB.Where id(%d) error(%d)", param.ID, err)
		return
	}
	if err = s.dao.DB.Model(&model.Game{}).Where("status=?", _statusOn).Where("id IN (?)", gids).Find(&games).Error; err != nil {
		log.Error("EditTeam check game ids Error (%v)", err)
		return
	}
	if len(games) == 0 {
		log.Error("EditTeam games(%v) not found", gids)
		err = ecode.RequestErr
		return
	}
	if err = s.dao.DB.Model(&model.GIDMap{}).Where("oid=?", param.ID).Where("type=?", model.TypeTeam).Find(&preGidMaps).Error; err != nil {
		log.Error("EditTeam games(%v) not found", gids)
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
			addGidMaps = append(addGidMaps, &model.GIDMap{Type: model.TypeTeam, Oid: param.ID, Gid: gid})
		}
	}
	tx := s.dao.DB.Begin()
	if err = tx.Error; err != nil {
		log.Error("s.dao.DB.Begin error(%v)", err)
		return
	}
	if err = tx.Model(&model.Team{}).Save(param).Error; err != nil {
		log.Error("EditTeam Team Update(%+v) error(%v)", param, err)
		err = tx.Rollback().Error
		return
	}
	if len(upGidMapAdd) > 0 {
		if err = tx.Model(&model.GIDMap{}).Where("id IN (?)", upGidMapAdd).Updates(map[string]interface{}{"is_deleted": _notDeleted}).Error; err != nil {
			log.Error("EditTeam GIDMap Add(%+v) error(%v)", upGidMapAdd, err)
			err = tx.Rollback().Error
			return
		}
	}
	if len(upGidMapDel) > 0 {
		if err = tx.Model(&model.GIDMap{}).Where("id IN (?)", upGidMapDel).Updates(map[string]interface{}{"is_deleted": _deleted}).Error; err != nil {
			log.Error("EditTeam GIDMap Del(%+v) error(%v)", upGidMapDel, err)
			err = tx.Rollback().Error
			return
		}
	}
	if len(addGidMaps) > 0 {
		if err = tx.Model(&model.GIDMap{}).Exec(model.GidBatchAddSQL(addGidMaps)).Error; err != nil {
			log.Error("EditTeam GIDMap Create(%+v) error(%v)", addGidMaps, err)
			err = tx.Rollback().Error
			return
		}
	}
	err = tx.Commit().Error
	return
}

// ForbidTeam .
func (s *Service) ForbidTeam(c context.Context, id int64, state int) (err error) {
	preTeam := new(model.Team)
	if err = s.dao.DB.Where("id=?", id).First(&preTeam).Error; err != nil {
		log.Error("TeamForbid s.dao.DB.Where id(%d) error(%d)", id, err)
		return
	}
	if err = s.dao.DB.Model(&model.Team{}).Where("id=?", id).Update(map[string]int{"is_deleted": state}).Error; err != nil {
		log.Error("TeamForbid s.dao.DB.Model error(%v)", err)
	}
	return
}
