package service

import (
	"context"
	"sort"
	"time"

	"go-common/app/interface/main/history/model"
	arcmdl "go-common/app/service/main/archive/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
)

const maxToView = 100

var (
	_empToView       = []*model.ToView{}
	_empArcToView    = []*model.ArcToView{}
	_empWebArcToView = []*model.WebArcToView{}
)

// AddToView add the user a toview item.
// +wd:ignore
func (s *Service) AddToView(c context.Context, mid, aid int64, ip string) (err error) {
	var (
		ok    bool
		count int
		arc   *arcmdl.View3
		now   = time.Now().Unix()
	)
	arcAid := &arcmdl.ArgAid2{Aid: aid}
	if arc, err = s.arcRPC.View3(c, arcAid); err != nil {
		return
	}
	if arc == nil {
		return
	}
	if arc.Rights.UGCPay == 1 {
		return ecode.ToViewPayUGC
	}
	if ok, err = s.toviewDao.Expire(c, mid); err != nil {
		return
	}
	if ok {
		if count, err = s.toviewDao.CntCache(c, mid); err != nil {
			return
		}
	} else {
		if _, count, err = s.toView(c, mid, 1, maxToView, ip); err != nil {
			return
		}
	}
	if count >= maxToView {
		err = ecode.ToViewOverMax
		return
	}
	if err = s.toviewDao.Add(c, mid, aid, now); err != nil {
		return
	}
	s.toviewCache.Do(c, func(ctx context.Context) {
		if errCache := s.toviewDao.AddCache(ctx, mid, aid, now); errCache != nil {
			log.Warn("s.toviewDao.AddCache(%d,%d,%d) err:%v", mid, aid, now, errCache)
		}
	})
	return
}

// AddMultiToView add toview items to user.
// +wd:ignore
func (s *Service) AddMultiToView(c context.Context, mid int64, aids []int64, ip string) (err error) {
	var (
		ok           bool
		count        int
		expectedAids []int64
		arcs         map[int64]*arcmdl.View3
		viewMap      map[int64]*model.ToView
		views        []*model.ToView
		now          = time.Now().Unix()
	)
	arcAids := &arcmdl.ArgAids2{Aids: aids}
	if arcs, err = s.arcRPC.Views3(c, arcAids); err != nil {
		return
	} else if len(arcs) == 0 {
		return
	}
	for _, v := range arcs {
		if v.Rights.UGCPay == 1 {
			return ecode.ToViewPayUGC
		}
	}
	if ok, err = s.toviewDao.Expire(c, mid); err != nil {
		return
	}
	if ok {
		if viewMap, err = s.toviewDao.CacheMap(c, mid); err != nil {
			return
		}
	} else {
		var list []*model.ToView
		list, _, err = s.toView(c, mid, 1, maxToView, ip)
		if err != nil {
			return
		}
		viewMap = make(map[int64]*model.ToView)
		for _, v := range list {
			if v == nil {
				continue
			}
			viewMap[v.Aid] = v
		}
	}
	count = len(viewMap)
	if count >= maxToView {
		err = ecode.ToViewOverMax
		return
	}
	expectedLen := maxToView - count
	for _, aid := range aids {
		if _, exist := viewMap[aid]; !exist {
			expectedAids = append(expectedAids, aid)
			expectedLen--
			if expectedLen == 0 {
				break
			}
		}
	}
	if err = s.toviewDao.Adds(c, mid, expectedAids, now); err != nil {
		return
	}
	if ok {
		s.toviewCache.Do(c, func(ctx context.Context) {
			for _, aid := range expectedAids {
				views = append(views, &model.ToView{Aid: aid, Unix: now})
			}
			if errCache := s.toviewDao.AddCacheList(ctx, mid, views); errCache != nil {
				log.Warn("s.toviewDao.AddCacheList(%d,%v) err:%v", mid, views, errCache)
			}
		})
	}
	return
}

// RemainingToView add toview items to user.
// +wd:ignore
func (s *Service) RemainingToView(c context.Context, mid int64, ip string) (remaining int, err error) {
	var (
		count int
	)
	if _, count, err = s.toView(c, mid, 1, maxToView, ip); err != nil {
		return
	}
	remaining = maxToView - count
	return
}

// ClearToView clear the user toview items.
// +wd:ignore
func (s *Service) ClearToView(c context.Context, mid int64) (err error) {
	if err = s.toviewDao.ClearCache(c, mid); err != nil {
		return
	}
	s.userActionLog(mid, model.ToviewClear)
	return s.toviewDao.Clear(c, mid)
}

// DelToView delete the user to videos.
// +wd:ignore
func (s *Service) DelToView(c context.Context, mid int64, aids []int64, viewed bool, ip string) (err error) {
	var (
		delAids []int64
		list    []*model.ToView
		rhs, hs map[int64]*model.History
		rs      *model.History
	)
	// viewed del all viewed
	if viewed {
		if list, _, err = s.toView(c, mid, 1, maxToView, ip); err != nil {
			return
		}
		for _, l := range list {
			aids = append(aids, l.Aid)
		}
		if len(aids) == 0 {
			return
		}
		if hs, err = s.historyDao.AidsMap(c, mid, aids); err != nil {
			return
		}
		rhs, _ = s.historyDao.CacheMap(c, mid)
		for _, rs = range rhs {
			hs[rs.Aid] = rs
		}
		for k, v := range hs {
			if v.Pro >= 30 || v.Pro == -1 {
				delAids = append(delAids, k)
			}
		}
		if len(delAids) == 0 {
			return
		}
		if err = s.toviewDao.Del(c, mid, delAids); err != nil {
			return
		}
		s.toviewCache.Do(c, func(ctx context.Context) {
			s.toviewDao.DelCaches(ctx, mid, delAids)
		})
		return
	}
	if err = s.toviewDao.Del(c, mid, aids); err != nil {
		return
	}
	s.toviewCache.Do(c, func(ctx context.Context) {
		s.toviewDao.DelCaches(ctx, mid, aids)
	})
	return
}

// WebToView get videos of user view history.
// +wd:ignore
func (s *Service) WebToView(c context.Context, mid int64, pn, ps int, ip string) (res []*model.WebArcToView, count int, err error) {
	var (
		ok          bool
		aids, epids []int64
		avs         map[int64]*arcmdl.View3
		views       []*model.ToView
		v           *model.ToView
		hm          map[int64]*model.History
		av          *arcmdl.View3
		epban       = make(map[int64]*model.Bangumi)
		seasonMap   = make(map[int64]*model.BangumiSeason)
	)
	if views, count, err = s.toView(c, mid, pn, ps, ip); err != nil {
		return
	}
	if len(views) == 0 {
		return
	}
	for _, v = range views {
		if v != nil {
			aids = append(aids, v.Aid)
		}
	}
	argAids := &arcmdl.ArgAids2{Aids: aids}
	if avs, err = s.arcRPC.Views3(c, argAids); err != nil {
		log.Error("s.arcRPC.Views3(arcAids:(%v), arcs) error(%v)", aids, err)
		return
	}
	if len(avs) == 0 {
		return
	}
	seasonMap, epids = s.season(c, mid, aids, ip)
	// bangumi info
	if len(epids) > 0 {
		epban = s.bangumi(c, mid, epids)
	}
	if hm, err = s.toViewPro(c, mid, aids); err != nil {
		err = nil
	}
	res = make([]*model.WebArcToView, 0, len(aids))
	for _, v = range views {
		if v == nil {
			count--
			continue
		}
		// NOTE compat android
		if av, ok = avs[v.Aid]; !ok || av == nil {
			count--
			continue
		}
		// NOTE all no pay
		av.Rights.Movie = 0
		at := &model.WebArcToView{View3: av}
		at.AddTime = v.Unix
		if hm[v.Aid] != nil {
			at.Cid = hm[v.Aid].Cid
			at.Progress = hm[v.Aid].Pro
		}
		if season, ok := seasonMap[v.Aid]; ok && season != nil {
			if bangumi, ok := epban[season.Epid]; ok && bangumi != nil {
				at.BangumiInfo = bangumi
			}
		}
		res = append(res, at)
	}
	if len(res) == 0 {
		res = _empWebArcToView
	}
	return
}

// ToView get videos of user view history.
// +wd:ignore
func (s *Service) ToView(c context.Context, mid int64, pn, ps int, ip string) (res []*model.ArcToView, count int, err error) {
	var (
		ok    bool
		aids  []int64
		avs   map[int64]*arcmdl.View3
		views []*model.ToView
		v     *model.ToView
		hm    map[int64]*model.History
		av    *arcmdl.View3
	)
	res = _empArcToView
	if views, count, err = s.toView(c, mid, pn, ps, ip); err != nil {
		return
	}
	if len(views) == 0 {
		return
	}
	for _, v = range views {
		if v != nil {
			aids = append(aids, v.Aid)
		}
	}
	argAids := &arcmdl.ArgAids2{Aids: aids}
	if avs, err = s.arcRPC.Views3(c, argAids); err != nil {
		log.Error("s.arcRPC.Views3(%v) error(%v)", aids, err)
		return
	}
	if len(avs) == 0 {
		return
	}
	if hm, err = s.toViewPro(c, mid, aids); err != nil {
		err = nil
	}
	res = make([]*model.ArcToView, 0, len(aids))
	for _, v = range views {
		if v == nil {
			count--
			continue
		}
		// NOTE compat android
		if av, ok = avs[v.Aid]; !ok || av.Archive3 == nil {
			count--
			continue
		}
		// NOTE all no pay
		av.Rights.Movie = 0
		at := &model.ArcToView{
			Archive3: av.Archive3,
			Count:    len(av.Pages),
		}
		at.AddTime = v.Unix
		// get cid and progress
		if hm[v.Aid] != nil {
			at.Cid = hm[v.Aid].Cid
			at.Progress = hm[v.Aid].Pro
		}
		for n, p := range av.Pages {
			if p.Cid == at.Cid {
				p.Page = int32(n + 1)
				at.Page = p
				break
			}
		}
		res = append(res, at)
	}
	return
}

// toView return ToSee of After the paging data.
func (s *Service) toView(c context.Context, mid int64, pn, ps int, ip string) (res []*model.ToView, count int, err error) {
	var (
		ok    bool
		start = (pn - 1) * ps
		end   = start + ps - 1
	)
	if ok, err = s.toviewDao.Expire(c, mid); err != nil {
		return
	}
	if ok {
		if end > maxToView {
			end = maxToView
		}
		if res, err = s.toviewDao.Cache(c, mid, start, end); err != nil {
			return
		}
		count, err = s.toviewDao.CntCache(c, mid)
		if count > maxToView {
			count = maxToView
		}
		return
	}
	var views []*model.ToView
	var viewMap = make(map[int64]*model.ToView)
	if viewMap, err = s.toviewDao.MapInfo(c, mid, nil); err != nil {
		return
	}
	if len(viewMap) == 0 {
		res = _empToView
		return
	}
	for _, v := range viewMap {
		views = append(views, v)
	}
	sort.Sort(model.ToViews(views))
	if count = len(views); count > maxToView {
		count = maxToView
		views = views[:count]
	}
	switch {
	case count > start && count > end:
		res = views[start : end+1]
	case count > start && count <= end:
		res = views[start:]
	default:
		res = _empToView
	}
	s.toviewCache.Do(c, func(ctx context.Context) {
		if errCache := s.toviewDao.AddCacheList(ctx, mid, views); errCache != nil {
			log.Warn("s.toviewDao.AddCacheList(%d,%v) err:%v", mid, views, errCache)
		}
	})
	return
}

func (s *Service) toViewPro(c context.Context, mid int64, aids []int64) (res map[int64]*model.History, err error) {
	var (
		miss []int64
		hm   map[int64]*model.History
	)
	if res, miss, err = s.historyDao.Cache(c, mid, aids); err != nil {
		err = nil
	} else if len(res) == len(aids) {
		return
	}
	if len(res) == 0 {
		res = make(map[int64]*model.History)
		miss = aids
	}
	if hm, err = s.historyDao.AidsMap(c, mid, miss); err != nil {
		err = nil
	}
	for k, v := range hm {
		res[k] = v
	}
	return
}

func (s *Service) season(c context.Context, mid int64, aids []int64, ip string) (seasonMap map[int64]*model.BangumiSeason, epids []int64) {
	var (
		n       = 50
		seasonM = make(map[int64]*model.BangumiSeason, n)
	)
	seasonMap = make(map[int64]*model.BangumiSeason, n)
	for len(aids) > 0 {
		if n > len(aids) {
			n = len(aids)
		}
		seasonM, _ = s.historyDao.BangumisByAids(c, mid, aids[:n], ip)
		aids = aids[n:]
		for k, v := range seasonM {
			epids = append(epids, v.Epid)
			seasonMap[k] = v
		}
	}
	return
}

// ManagerToView manager get mid toview list.
// +wd:ignore
func (s *Service) ManagerToView(c context.Context, mid int64, ip string) ([]*model.ToView, error) {
	return s.toviewDao.ListInfo(c, mid, nil)
}
