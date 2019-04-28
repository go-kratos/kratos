package service

import (
	"context"
	"time"

	"go-common/app/service/main/tag/model"
	xsql "go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"
)

var emptySubs = []*model.Sub{{Tid: -1}}

// AddSub user add sub .
func (s *Service) AddSub(c context.Context, mid int64, tids []int64, ip string) (err error) {
	var (
		ts      []*model.Tag
		addSubs []*model.Sub
		subMap  map[int64]int32
		now     = time.Now()
	)
	if subMap, err = s.subTids(c, mid); err != nil {
		return
	}
	if len(subMap) >= model.MaxSubNum {
		return ecode.TagMaxSub
	}
	if ts, err = s.tags(c, tids); err != nil {
		return
	}
	if len(ts) == 0 {
		return ecode.TagNotExist
	}
	for _, t := range ts {
		if v, ok := subMap[t.ID]; !ok || v != 1 {
			sub := &model.Sub{
				Tid:   t.ID,
				Mid:   mid,
				MTime: xtime.Time(now.Unix()),
			}
			addSubs = append(addSubs, sub)
		}
	}
	addLen := len(addSubs)
	if addLen == 0 {
		return ecode.TagTagIsSubed
	}
	if (len(subMap) + addLen) > model.MaxSubNum {
		return ecode.TagMaxSub
	}
	_, err = s.addSub(c, mid, addSubs, now)
	return
}

// CancelSub user cancel sub tag .
func (s *Service) CancelSub(c context.Context, mid, tid int64, ip string) (err error) {
	var (
		ok  bool
		now = time.Now()
	)
	if ok, err = s.isSubTid(c, mid, tid); err != nil {
		return
	}
	if !ok {
		return ecode.TagNotSub
	}
	_, err = s.cancelSub(c, mid, tid, now)
	return
}

// SubTags user sub tags .
func (s *Service) SubTags(c context.Context, mid int64, pn, ps, order int) (ts []*model.Tag, total int, err error) {
	var (
		ok   bool
		subs []*model.Sub
	)
	if ok, err = s.dao.ExpireSubCache(c, mid); err != nil {
		return
	}
	if ok {
		if subs, err = s.dao.SubListCache(c, mid); err != nil || len(subs) == 0 {
			return
		}
	} else {
		if subs, err = s.dao.SubList(c, mid); err != nil {
			return
		}
	}
	if len(subs) > 0 {
		if ts, total, err = s.sortSubTags(c, subs, pn, ps, order); err != nil {
			return
		}
	} else {
		subs = emptySubs
	}
	if !ok {
		s.cacheCh.Save(func() {
			s.dao.AddSubListCache(context.Background(), mid, subs)
		})
	}
	return
}

func (s *Service) sortSubTags(c context.Context, subs []*model.Sub, pn, ps, order int) (ts []*model.Tag, total int, err error) {
	var (
		tids  []int64
		start = (pn - 1) * ps
		end   = start + ps - 1
	)
	subSort := &model.SubSort{Subs: subs, Order: order}
	tids = subSort.Sort()
	total = len(tids)
	switch {
	case total > start && total > end:
		tids = tids[start : end+1]
	case total > start && total <= end:
		tids = tids[start:]
	default:
		return
	}
	if total == 0 || len(tids) == 0 {
		return
	}
	tm, err := s.tagMap(c, tids)
	if err != nil || len(tm) == 0 {
		return
	}
	for _, tid := range tids {
		if t, ok := tm[tid]; ok && t != nil {
			t.Attention = 1
			ts = append(ts, t)
		}
	}
	return
}

func (s *Service) isSubTid(c context.Context, mid, tid int64) (sub bool, err error) {
	var (
		ok     bool
		subMap map[int64]*model.Sub
	)
	if ok, err = s.dao.ExpireSubCache(c, mid); err != nil {
		return
	}
	if ok {
		sub, err = s.dao.IsSubCache(c, mid, tid)
		return
	}
	if _, subMap, err = s.dao.Sub(c, mid); err != nil {
		return
	}
	_, sub = subMap[tid]
	if len(subMap) == 0 {
		s.cache.Save(func() {
			s.dao.AddSubListCache(context.Background(), mid, emptySubs)
		})
		return
	}
	s.cache.Save(func() {
		s.dao.AddSubMapCache(context.Background(), mid, subMap)
	})
	return
}

func (s *Service) isSubTids(c context.Context, mid int64, tids []int64) (res map[int64]int32, err error) {
	var (
		ok   bool
		subs []*model.Sub
	)
	if ok, err = s.dao.ExpireSubCache(c, mid); err != nil {
		return
	}
	if ok {
		return s.dao.IsSubsCache(c, mid, tids)
	}
	if subs, err = s.dao.SubList(c, mid); err != nil {
		return
	}
	res = make(map[int64]int32, len(subs))
	if len(subs) > 0 {
		for _, sub := range subs {
			res[sub.Tid] = 1
		}
	} else {
		subs = emptySubs
	}
	s.cache.Save(func() {
		s.dao.AddSubListCache(context.Background(), mid, subs)
	})
	return
}

func (s *Service) subTids(c context.Context, mid int64) (tids map[int64]int32, err error) {
	var subs []*model.Sub
	if tids, err = s.dao.SubTidsCache(c, mid); err != nil {
		log.Error("s.dao.SubTidsCache(%d ) error(%v)", mid, err)
		return nil, err
	}
	if tids != nil {
		return
	}
	if subs, err = s.dao.SubList(c, mid); err != nil {
		log.Error("s.dao.SubLis(%d ) error(%v)", mid, err)
		return nil, err
	}
	if subs == nil {
		return nil, nil
	}
	for _, sub := range subs {
		tids[sub.Tid] = 1
	}
	s.cache.Save(func() {
		s.dao.AddSubListCache(context.Background(), mid, subs)
	})
	return
}

func (s *Service) addSub(c context.Context, mid int64, subs []*model.Sub, now time.Time) (row int64, err error) {
	var (
		tx      *xsql.Tx
		tids    []int64
		addSubs []*model.Sub
	)
	if tx, err = s.dao.BeginTran(c); err != nil {
		log.Error("s.sub.beginTran() error(%v)", err)
		return
	}
	for _, sub := range subs {
		if row, err = s.dao.TxAddSub(tx, mid, sub.Tid); err != nil {
			tx.Rollback()
			return
		}
		if row, err = s.dao.TxUpTagSubCount(tx, sub.Tid, 1); err != nil {
			tx.Rollback()
			return
		}
		addSubs = append(addSubs, sub)
		tids = append(tids, sub.Tid)
	}
	if err = tx.Commit(); err != nil {
		log.Error("AddSub tx.Commit(), error(%v)", err)
		return
	}
	if len(tids) > 0 {
		s.cacheCh.Save(func() {
			s.dao.AddSubListCache(context.Background(), mid, addSubs)
			s.dao.DelCountsCache(context.Background(), tids)
		})
	}
	return
}

func (s *Service) cancelSub(c context.Context, mid, tid int64, now time.Time) (row int64, err error) {
	var (
		tx *xsql.Tx
	)
	if tx, err = s.dao.BeginTran(c); err != nil {
		log.Error("s.sub.beginTran() error(%v)", err)
		return
	}
	if row, err = s.dao.TxDelSub(tx, mid, tid); err != nil || row == 0 {
		tx.Rollback()
		return
	}
	if row, err = s.dao.TxUpTagSubCount(tx, tid, -1); err != nil || row == 0 {
		tx.Rollback()
		return
	}
	if err = tx.Commit(); err != nil {
		log.Error("AddSub tx.Commit(), error(%v)", err)
		return
	}
	s.cacheCh.Save(func() {
		s.dao.DelCountCache(context.Background(), tid)
		s.dao.DelSubTidCache(context.Background(), mid, tid)
		s.updateSortSub(context.Background(), mid, tid)
	})
	return
}

func (s *Service) updateSortSub(c context.Context, mid, tid int64) (err error) {
	tidMap, err := s.dao.AllCustomSubSort(c, mid)
	if err != nil {
		return
	}
	newSortMap := make(map[int32][]int64, len(tidMap))
	for tp, tids := range tidMap {
		newTids := make([]int64, 0, len(tids))
		var change bool
		for _, v := range tids {
			if v == tid {
				change = true
				continue
			}
			newTids = append(newTids, v)
		}
		if change {
			newSortMap[tp] = newTids
		}
	}
	if len(newSortMap) == 0 {
		return
	}
	err = s.dao.AddSubChannels(c, mid, newSortMap)
	for tp := range newSortMap {
		s.dao.DelSubSortCache(context.Background(), mid, int(tp))
	}
	return
}
