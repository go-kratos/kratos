package service

import (
	"context"
	"time"

	"go-common/app/job/main/up-rating/model"
)

// BaseInfoOffEnd get offset and end
func (s *Service) BaseInfoOffEnd(c context.Context, date time.Time) (offset, end int, err error) {
	start, err := s.dao.BaseInfoStart(c, date)
	if err != nil {
		return
	}
	offset = start - 1
	end, err = s.dao.BaseInfoEnd(c, date)
	return
}

// RatingOffEnd get offset and end
func (s *Service) RatingOffEnd(c context.Context, date time.Time) (offset, end, count int, err error) {
	start, err := s.dao.RatingStart(c, date)
	if err != nil {
		return
	}
	offset = start - 1
	end, err = s.dao.RatingEnd(c, date)
	if err != nil {
		return
	}
	count, err = s.dao.RatingCount(c, date)
	return
}

// BaseInfo get base infos
func (s *Service) BaseInfo(c context.Context, date time.Time, start, end int, ch chan []*model.BaseInfo) (err error) {
	defer close(ch)
	for {
		var bs []*model.BaseInfo
		bs, start, err = s.dao.GetBaseInfo(c, date.Month(), start, end, _limit)
		if err != nil {
			return
		}
		if len(bs) == 0 {
			break
		}
		ch <- bs
	}
	return
}

// BaseTotal get total base
func (s *Service) BaseTotal(c context.Context, date time.Time) (total map[int64]*model.BaseInfo, err error) {
	total = make(map[int64]*model.BaseInfo)
	var id int64
	for {
		var bs []*model.BaseInfo
		bs, err = s.dao.GetBaseTotal(c, date, id, int64(_limit))
		if err != nil {
			return
		}
		for _, b := range bs {
			total[b.MID] = b
		}
		if len(bs) < _limit {
			break
		}
		id = bs[len(bs)-1].ID
	}
	return
}
