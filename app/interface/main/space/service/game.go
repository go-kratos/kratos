package service

import (
	"context"

	"go-common/app/interface/main/space/model"
)

var (
	_emptyGameList    = make([]*model.Game, 0)
	_emptyAppGameList = make([]*model.AppGame, 0)
)

// LastPlayGame get last play game by mid
func (s *Service) LastPlayGame(c context.Context, mid, vmid int64) (data []*model.Game, err error) {
	if mid != vmid {
		if err = s.privacyCheck(c, vmid, model.PcyGame); err != nil {
			return
		}
	}
	if data, err = s.dao.LastPlayGame(c, vmid); err != nil {
		err = nil
		data = _emptyGameList
		return
	}
	if len(data) == 0 {
		data = _emptyGameList
	}
	return
}

// AppPlayedGame get app played games.
func (s *Service) AppPlayedGame(c context.Context, mid, vmid int64, platform string, pn, ps int) (data []*model.AppGame, count int, err error) {
	if mid != vmid {
		if err = s.privacyCheck(c, vmid, model.PcyGame); err != nil {
			return
		}
	}
	if data, count, err = s.dao.AppPlayedGame(c, vmid, platform, pn, ps); err != nil {
		err = nil
		data = _emptyAppGameList
		return
	}
	if len(data) == 0 {
		data = _emptyAppGameList
	}
	return
}
