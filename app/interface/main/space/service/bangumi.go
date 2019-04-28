package service

import (
	"context"

	"go-common/app/interface/main/space/model"
)

var _emptyBangumiList = make([]*model.Bangumi, 0)

// BangumiList get bangumi list by mid.
func (s *Service) BangumiList(c context.Context, mid, vmid int64, pn, ps int) (data []*model.Bangumi, count int, err error) {
	if mid != vmid {
		if err = s.privacyCheck(c, vmid, model.PcyBangumi); err != nil {
			return
		}
	}
	if data, count, err = s.dao.BangumiList(c, vmid, pn, ps); err != nil {
		return
	}
	if len(data) == 0 {
		data = _emptyBangumiList
	}
	return
}

// BangumiConcern bangumi concern.
func (s *Service) BangumiConcern(c context.Context, mid, seasonID int64) (err error) {
	return s.dao.BangumiConcern(c, mid, seasonID)
}

// BangumiUnConcern bangumi unconcern.
func (s *Service) BangumiUnConcern(c context.Context, mid, seasonID int64) (err error) {
	return s.dao.BangumiUnConcern(c, mid, seasonID)
}
