package service

import (
	"context"
	"fmt"

	"go-common/app/admin/main/esports/model"
	"go-common/library/log"
)

var _emptyGameList = make([]*model.Game, 0)

// GameInfo .
func (s *Service) GameInfo(c context.Context, id int64) (game *model.Game, err error) {
	game = new(model.Game)
	if err = s.dao.DB.Where("id=?", id).First(&game).Error; err != nil {
		log.Error("GameInfo Error (%v)", err)
	}
	return
}

// GameList .
func (s *Service) GameList(c context.Context, pn, ps int64, title string) (list []*model.Game, count int64, err error) {
	source := s.dao.DB.Model(&model.Game{})
	if title != "" {
		likeStr := fmt.Sprintf("%%%s%%", title)
		source = source.Where("title like ?", likeStr)
	}
	source.Count(&count)
	if err = source.Offset((pn - 1) * ps).Limit(ps).Find(&list).Error; err != nil {
		log.Error("GameList Error (%v)", err)
	}
	return
}

// AddGame .
func (s *Service) AddGame(c context.Context, param *model.Game) (err error) {
	// TODO check name exist
	if err = s.dao.DB.Model(&model.Game{}).Create(param).Error; err != nil {
		log.Error("AddGame s.dao.DB.Model Create(%+v) error(%v)", param, err)
	}
	return
}

// EditGame .
func (s *Service) EditGame(c context.Context, param *model.Game) (err error) {
	preGame := new(model.Game)
	if err = s.dao.DB.Where("id=?", param.ID).First(&preGame).Error; err != nil {
		log.Error("EditGame s.dao.DB.Where id(%d) error(%d)", param.ID, err)
		return
	}
	if err = s.dao.DB.Model(&model.Game{}).Update(param).Error; err != nil {
		log.Error("EditGame s.dao.DB.Model Update(%+v) error(%v)", param, err)
	}
	return
}

// ForbidGame .
func (s *Service) ForbidGame(c context.Context, id int64, state int) (err error) {
	preGame := new(model.Game)
	if err = s.dao.DB.Where("id=?", id).First(&preGame).Error; err != nil {
		log.Error("GameForbid s.dao.DB.Where id(%d) error(%d)", id, err)
		return
	}
	if err = s.dao.DB.Model(&model.Game{}).Where("id=?", id).Update(map[string]int{"status": state}).Error; err != nil {
		log.Error("GameForbid s.dao.DB.Model error(%v)", err)
	}
	return
}

// gameList return game info map with oid key.
func (s *Service) gameList(c context.Context, typ int, oids []int64) (list map[int64][]*model.Game, err error) {
	var (
		gidMaps []*model.GIDMap
		gids    []int64
		games   []*model.Game
	)
	if len(oids) == 0 {
		return
	}
	if err = s.dao.DB.Model(&model.GIDMap{}).Where("is_deleted=?", _notDeleted).Where("type=?", typ).Where("oid IN(?)", oids).Find(&gidMaps).Error; err != nil {
		log.Error("gameList gidMap Error (%v)", err)
		return
	}
	if len(gidMaps) == 0 {
		return
	}
	gidMap := make(map[int64]int64, len(gidMaps))
	oidGidMap := make(map[int64][]int64, len(gidMaps))
	for _, v := range gidMaps {
		oidGidMap[v.Oid] = append(oidGidMap[v.Oid], v.Gid)
		if _, ok := gidMap[v.Gid]; ok {
			continue
		}
		gids = append(gids, v.Gid)
		gidMap[v.Gid] = v.Gid
	}
	if err = s.dao.DB.Model(&model.Game{}).Where("status=?", _statusOn).Where("id IN(?)", gids).Find(&games).Error; err != nil {
		log.Error("gameList games Error (%v)", err)
		return
	}
	if len(games) == 0 {
		return
	}
	gameMap := make(map[int64]*model.Game, len(games))
	for _, v := range games {
		gameMap[v.ID] = v
	}
	list = make(map[int64][]*model.Game, len(oids))
	for _, oid := range oids {
		if ids, ok := oidGidMap[oid]; ok {
			for _, id := range ids {
				if game, ok := gameMap[id]; ok {
					list[oid] = append(list[oid], game)
				}
			}
		}
	}
	return
}

// Types return data page game types.
func (s *Service) Types(c context.Context) (list map[int64]string, err error) {
	list = make(map[int64]string, len(s.c.GameTypes))
	for _, tp := range s.c.GameTypes {
		list[tp.ID] = tp.Name
	}
	return
}
