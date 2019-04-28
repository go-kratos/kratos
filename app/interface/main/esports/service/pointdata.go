package service

import (
	"context"

	"go-common/app/interface/main/esports/model"
)

// Game get game.
func (s *Service) Game(c context.Context, p *model.ParamGame) (rs map[int64]*model.Game, err error) {
	var (
		ok      bool
		games   []*model.Game
		gameMap map[int64]*model.Game
	)
	rs = make(map[int64]*model.Game, len(p.GameIDs))
	if p.Tp == _lolType {
		if games, ok = s.lolGameMap.Data[p.MatchID]; !ok {
			return
		}
	} else if p.Tp == _dotaType {
		if games, ok = s.dotaGameMap.Data[p.MatchID]; !ok {
			return
		}
	}
	count := len(games)
	if count == 0 {
		return
	}
	gameMap = make(map[int64]*model.Game, count)
	for _, game := range games {
		gameMap[game.ID] = game
	}
	for _, id := range p.GameIDs {
		if game, ok := gameMap[id]; ok {
			rs[id] = game
		}
	}
	return
}

// Items get items.
func (s *Service) Items(c context.Context, p *model.ParamLeidas) (rs map[int64]*model.Item, err error) {
	rs = make(map[int64]*model.Item, len(p.IDs))
	if p.Tp == _lolType {
		for _, id := range p.IDs {
			if item, ok := s.lolItemsMap.Data[id]; ok {
				rs[id] = item
			}
		}
	} else if p.Tp == _dotaType {
		for _, id := range p.IDs {
			if item, ok := s.dotaItemsMap.Data[id]; ok {
				rs[id] = item
			}
		}
	}
	return
}

// Heroes lol:champions ; dota2 heroes.
func (s *Service) Heroes(c context.Context, p *model.ParamLeidas) (rs interface{}, err error) {
	var (
		champions  map[int64]*model.Champion
		dotaHeroes map[int64]*model.Hero
	)
	if p.Tp == _lolType {
		champions = make(map[int64]*model.Champion, len(p.IDs))
		for _, id := range p.IDs {
			if item, ok := s.lolChampions.Data[id]; ok {
				champions[id] = item
			}
		}
		rs = champions
	} else if p.Tp == _dotaType {
		dotaHeroes = make(map[int64]*model.Hero, len(p.IDs))
		for _, id := range p.IDs {
			if item, ok := s.dotaHeroes.Data[id]; ok {
				dotaHeroes[id] = item
			}
		}
		rs = dotaHeroes
	}
	return
}

// Abilities lol:spells;dota2:abilities.
func (s *Service) Abilities(c context.Context, p *model.ParamLeidas) (rs interface{}, err error) {
	infos := make(map[int64]*model.LdInfo, len(p.IDs))
	if p.Tp == _lolType {
		for _, id := range p.IDs {
			if info, ok := s.lolSpells.Data[id]; ok {
				infos[id] = info
			}
		}
		rs = infos
	} else if p.Tp == _dotaType {
		for _, id := range p.IDs {
			if info, ok := s.dotaAbilities.Data[id]; ok {
				infos[id] = info
			}
		}
		rs = infos
	}
	return
}

// Players get players.
func (s *Service) Players(c context.Context, p *model.ParamLeidas) (rs interface{}, err error) {
	infos := make(map[int64]*model.LdInfo, len(p.IDs))
	if p.Tp == _lolType {
		for _, id := range p.IDs {
			if info, ok := s.lolPlayers.Data[id]; ok {
				infos[id] = info
			}
		}
		rs = infos
	} else if p.Tp == _dotaType {
		for _, id := range p.IDs {
			if info, ok := s.dotaPlayers.Data[id]; ok {
				infos[id] = info
			}
		}
		rs = infos
	}
	return
}
