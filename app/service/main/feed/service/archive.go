package service

import (
	"context"
	"sync"

	accmdl "go-common/app/service/main/account/model"
	"go-common/app/service/main/archive/api"
	arcmdl "go-common/app/service/main/archive/model/archive"
	"go-common/app/service/main/feed/dao"
	"go-common/library/log"
	"go-common/library/time"

	"go-common/library/sync/errgroup"
)

const _upsArcBulkSize = 50

// AddArc add archive when archive passed.
func (s *Service) AddArc(c context.Context, mid, aid int64, pubDate int64, ip string) (err error) {
	res, err := s.dao.ExpireUppersCache(c, []int64{mid})
	if (err == nil) && res[mid] {
		err = s.dao.AddUpperCache(c, mid, &arcmdl.AidPubTime{Aid: aid, PubDate: time.Time(pubDate)})
	}
	log.Info("service.AddArc(aid: %v, mid: %v, pubDate:%v, ip: %v) err(%v)", aid, mid, pubDate, ip, err)
	dao.PromInfo("archive:稿件过审增加缓存")
	return
}

// DelArc delete archive when archive not passed.
func (s *Service) DelArc(c context.Context, mid, aid int64, ip string) (err error) {
	s.dao.DelArchiveCache(c, aid)
	err = s.dao.DelUpperCache(c, mid, aid)
	log.Info("service.DelArc(aid: %v, mid: %v, ip: %v) err(%v)", aid, mid, ip, err)
	dao.PromInfo("archive:稿件删除缓存")
	return
}

// upArcs get new archives of uppers.
func (s *Service) upArcs(c context.Context, minTotalCount int, ip string, mids ...int64) (res map[int64][]*arcmdl.AidPubTime, err error) {
	var (
		start, end     int
		cache          = true
		cached, missed []int64
		oks            map[int64]bool
		tmpRes         map[int64][]*arcmdl.AidPubTime
		length         = len(mids)
	)
	if length == 0 {
		return
	}
	end = minTotalCount/length + s.c.Feed.MinUpCnt
	if oks, err = s.dao.ExpireUppersCache(c, mids); err != nil {
		cache = false
		missed = mids
	} else {
		for mid, ok := range oks {
			if ok {
				cached = append(cached, mid)
			} else {
				missed = append(missed, mid)
			}
		}
	}
	if len(cached) > 0 {
		if res, err = s.dao.UppersCaches(c, cached, start, end); err != nil {
			dao.PromError("up主缓存", "dao.UppersCaches(%v) error(%v)", cached, err)
			missed = mids
			err = nil
			cache = false
		}
	}
	if res == nil {
		res = make(map[int64][]*arcmdl.AidPubTime, len(mids))
	}
	if len(missed) > 0 {
		if tmpRes, err = s.upsPassed(c, missed, ip); err != nil {
			dao.PromError("up主回源", "upsPassed(%v) error(%v)", missed, err)
			err = nil
			tmpRes, _ = s.dao.UppersCaches(c, missed, start, end)
		} else {
			if cache {
				s.addCache(func() {
					s.dao.AddUpperCaches(context.Background(), tmpRes)
				})
			}
		}
		for mid, arcs := range tmpRes {
			if len(arcs) == 0 {
				continue
			}
			var tmp []*arcmdl.AidPubTime
			if len(arcs) > end+1 {
				tmp = arcs[start : end+1]
			} else {
				tmp = arcs
			}
			res[mid] = tmp
		}
	}
	return
}

// attenUpArcs get new archives of attention uppers.
func (s *Service) attenUpArcs(c context.Context, minTotalCount int, mid int64, ip string) (res map[int64][]*arcmdl.AidPubTime, err error) {
	var mids []int64
	arg := &accmdl.ArgMid{Mid: mid}
	if mids, err = s.accRPC.Attentions3(c, arg); err != nil {
		dao.PromError("关注rpc接口:Attentions", "s.accRPC.Attentions(%d) error(%v)", mid, err)
		return
	}
	return s.upArcs(c, minTotalCount, ip, mids...)
}

func (s *Service) upsPassed(c context.Context, mids []int64, ip string) (res map[int64][]*arcmdl.AidPubTime, err error) {
	dao.MissedCount.Add("up", int64(len(mids)))
	var (
		group      *errgroup.Group
		errCtx     context.Context
		midsLen, i int
		mutex      = sync.Mutex{}
	)
	res = make(map[int64][]*arcmdl.AidPubTime)
	group, errCtx = errgroup.WithContext(c)
	midsLen = len(mids)
	for ; i < midsLen; i += _upsArcBulkSize {
		var partMids []int64
		if i+_upsArcBulkSize > midsLen {
			partMids = mids[i:]
		} else {
			partMids = mids[i : i+_upsArcBulkSize]
		}
		group.Go(func() (err error) {
			var tmpRes map[int64][]*arcmdl.AidPubTime
			arg := &arcmdl.ArgUpsArcs2{Mids: partMids, Pn: 1, Ps: s.c.MultiRedis.MaxArcsNum, RealIP: ip}
			if tmpRes, err = s.arcRPC.UpsPassed2(errCtx, arg); err != nil {
				dao.PromError("up稿件回源RPC接口:UpsPassed2", "s.arcRPC.UpsPassed2(%+v) error(%v)", arg, err)
				err = nil
				return
			}
			mutex.Lock()
			for mid, arcs := range tmpRes {
				res[mid] = arcs
			}
			mutex.Unlock()
			return
		})
	}
	group.Wait()
	return
}

func (s *Service) archives(c context.Context, aids []int64, ip string) (res map[int64]*api.Arc, err error) {
	var (
		missed   []int64
		mutex    = sync.Mutex{}
		bulkSize = s.c.Feed.BulkSize
		addCache = true
	)
	if len(aids) == 0 {
		return
	}
	res, missed, err = s.dao.ArchivesCache(c, aids)
	if err != nil {
		dao.PromError("稿件缓存", "dao.ArchivesCache() error(%v)", err)
		missed = aids
		addCache = false
		err = nil
	} else if res != nil && len(missed) == 0 {
		return
	}
	if res == nil {
		res = make(map[int64]*api.Arc, len(aids))
	}
	group, errCtx := errgroup.WithContext(c)
	missedLen := len(missed)
	for i := 0; i < missedLen; i += bulkSize {
		var partAids []int64
		if i+bulkSize < missedLen {
			partAids = missed[i : i+bulkSize]
		} else {
			partAids = missed[i:missedLen]
		}
		group.Go(func() error {
			var (
				tmpRes map[int64]*api.Arc
				arcErr error
				arg    *arcmdl.ArgAids2
			)
			arg = &arcmdl.ArgAids2{Aids: partAids, RealIP: ip}
			if tmpRes, arcErr = s.arcRPC.Archives3(errCtx, arg); arcErr != nil {
				dao.PromError("稿件rpc接口:Archives2", "s.arcRPC.Archives3() error(%v)", err)
				// only log err message
				return nil
			}
			mutex.Lock()
			for aid, arc := range tmpRes {
				res[aid] = arc
			}
			mutex.Unlock()
			if addCache {
				s.addCache(func() {
					s.dao.AddArchivesCacheMap(context.Background(), tmpRes)
				})
			}
			return nil
		})
	}
	group.Wait()
	// check state
	for aid, arc := range res {
		if !arc.IsNormal() {
			delete(res, aid)
		}
	}
	return
}

func (s *Service) archive(c context.Context, aid int64, ip string) (res *api.Arc, err error) {
	arg := &arcmdl.ArgAid2{Aid: aid, RealIP: ip}
	res, err = s.arcRPC.Archive3(c, arg)
	return
}

// ChangeAuthor refresh feed cache
func (s *Service) ChangeAuthor(c context.Context, aid int64, oldMid int64, newMid int64, ip string) (err error) {
	s.dao.DelArchiveCache(c, aid)
	s.dao.DelUpperCache(c, oldMid, aid)
	arc, err := s.archive(c, aid, ip)
	if err != nil {
		dao.PromError("稿件转移", "s.archive(%v) error(%v)", aid, err)
		return
	}
	if !arc.IsNormal() {
		return
	}
	arc.Author.Mid = newMid
	res, err := s.dao.ExpireUppersCache(c, []int64{newMid})
	if (err == nil) && res[newMid] {
		err = s.dao.AddUpperCache(c, newMid, &arcmdl.AidPubTime{Aid: arc.Aid, PubDate: arc.PubDate})
	}
	return
}
