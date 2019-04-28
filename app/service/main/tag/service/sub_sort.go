package service

import (
	"context"

	"go-common/app/service/main/tag/model"
)

// AddCustomSubTags .
func (s *Service) AddCustomSubTags(c context.Context, mid int64, typ int, tids []int64, ip string) (err error) {
	_, err = s.dao.AddSubSort(c, mid, typ, tids)
	if err != nil {
		return
	}
	err = s.dao.AddSubSortCache(c, mid, typ, tids)
	return
}

// AddCustomSubChannels AddCustomSubChannels.
func (s *Service) AddCustomSubChannels(c context.Context, mid int64, typ int, tids []int64, ip string) (err error) {
	err = s.dao.AddSubChannel(c, mid, typ, tids)
	if err != nil {
		return
	}
	s.cacheCh.Save(func() {
		s.dao.DelSubSortCache(context.Background(), mid, typ)
	})
	return
}

// CustomSubTags user sub tags .
func (s *Service) CustomSubTags(c context.Context, mid int64, typ int, pn, ps, order int) (cts []*model.Tag, ts []*model.Tag, total int, err error) {
	var (
		ok, cached bool
		tidsort    []int64
		subs       []*model.Sub
		subm       = make(map[int64]*model.Sub)
	)
	tidsort, err = s.dao.SubSortCache(c, mid, typ)
	if err != nil {
		return
	}
	if len(tidsort) == 0 {
		tidsort, err = s.dao.CustomSubSort(c, mid, typ)
		if err != nil {
			return
		}
		if len(tidsort) > 0 {
			cached = true
		}
	}
	if ok, err = s.dao.ExpireSubCache(c, mid); err != nil {
		return
	}
	if ok {
		if subs, subm, err = s.dao.SubCache(c, mid); err != nil || len(subs) == 0 {
			return
		}
	} else {
		if subs, subm, err = s.dao.Sub(c, mid); err != nil {
			return
		}
	}
	if len(subs) > 0 {
		var tidsorted []int64
		for _, v := range tidsort {
			if _, ok := subm[v]; ok {
				tidsorted = append(tidsorted, v)
			}
		}
		if len(tidsorted) == 0 {
			tidsort = []int64{-1}
		} else {
			tidsort = tidsorted
		}
		if cts, ts, total, err = s.customSortSubTags(c, tidsort, subs, pn, ps, order); err != nil {
			return
		}
	} else {
		subs = emptySubs
	}
	s.cacheCh.Save(func() {
		s.dao.AddSubListCache(context.Background(), mid, subs)
		if cached {
			s.dao.AddSubSortCache(context.Background(), mid, typ, tidsort)
		}
	})
	return
}

func (s *Service) customSortSubTags(c context.Context, tidsort []int64, subs []*model.Sub, pn, ps, order int) (cts []*model.Tag, ts []*model.Tag, total int, err error) {
	var (
		tids, tidsTmp []int64
		start         = (pn - 1) * ps
		end           = start + ps - 1
		tidMap        = make(map[int64]bool, len(tidsort))
	)
	for _, t := range tidsort {
		tidMap[t] = true
	}
	subSort := &model.SubSort{Subs: subs, Order: order}
	tidsTmp = subSort.Sort()
	for _, v := range tidsTmp {
		if _, ok := tidMap[v]; ok {
			continue
		}
		tids = append(tids, v)
	}
	total = len(tids)
	switch {
	case total > start && total > end:
		tids = tids[start : end+1]
	case total > start && total <= end:
		tids = tids[start:]
	default:
		tids = []int64{}
	}
	tidsTmp = append(tids, tidsort...)
	if len(tidsTmp) == 0 {
		return
	}
	tm, err := s.tagMap(c, tidsTmp)
	if err != nil || len(tm) == 0 {
		return
	}
	for _, tid := range tidsort {
		if t, ok := tm[tid]; ok && t != nil {
			t.Attention = 1
			cts = append(cts, t)
		}
	}
	for _, tid := range tids {
		if t, ok := tm[tid]; ok && t != nil {
			t.Attention = 1
			ts = append(ts, t)
		}
	}
	total += len(cts)
	return
}
