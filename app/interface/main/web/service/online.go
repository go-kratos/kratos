package service

import (
	"context"
	"sync/atomic"
	"time"

	"go-common/app/interface/main/web/model"
	arcmdl "go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/log"
)

const (
	_onlineListNum      = 50
	_onlinePubdateLimit = 86400
	_onlineDanmuLimit   = 3
)

// OnlineArchiveCount Get Archive Count.
func (s *Service) OnlineArchiveCount(c context.Context) (rs *model.Online) {
	rs = &model.Online{
		RegionCount: s.regionCount,
		AllCount:    s.allArchivesCount,
		PlayOnline:  s.playOnline,
		WebOnline:   s.webOnline,
	}
	return
}

// OnlineList online archive list.
func (s *Service) OnlineList(c context.Context) (res []*model.OnlineArc, err error) {
	res = s.onlineArcs
	if len(res) > 0 {
		return
	}
	return s.dao.OnlineListBakCache(c)
}

func (s *Service) newCountproc() {
	var (
		allCount int64
		reIDs    []int16
		err      error
		res      map[int16]int
	)
	for {
		for len(s.rids) == 0 {
			time.Sleep(time.Second)
		}
		reIDs = []int16{}
		allCount = 0
		for regionID := range s.rids {
			reIDs = append(reIDs, int16(regionID))
		}
		arg := &archive.ArgRankTopsCount2{ReIDs: reIDs}
		if res, err = s.arc.RanksTopCount2(context.Background(), arg); err != nil {
			log.Error("s.arc.RanksTopCount2(%v) error (%v)", arg, err)
			time.Sleep(time.Second)
			continue
		} else if len(res) == 0 {
			log.Error("s.arc.RanksTopCount2(%v) res len(%d) == 0", arg, res)
			time.Sleep(time.Second)
			continue
		}
		s.regionCount = res
		for _, count := range res {
			allCount += int64(count)
		}
		if allCount > 0 {
			atomic.StoreInt64(&s.allArchivesCount, allCount)
		}
		time.Sleep(time.Duration(s.c.WEB.PullRegionInterval))
	}
}

func (s *Service) onlineCountproc() {
	var (
		count                 *model.OnlineCount
		liveCount             *model.LiveOnlineCount
		playOnline, webOnline int64
		err                   error
	)
	for {
		if count, err = s.dao.OnlineCount(context.Background()); err != nil {
			time.Sleep(time.Second)
			continue
		} else if count != nil {
			playOnline = count.ConnCount
			webOnline = count.IPCount
		}
		if liveCount, err = s.dao.LiveOnlineCount(context.Background()); err != nil || liveCount == nil {
			time.Sleep(time.Second)
			continue
		} else if liveCount != nil {
			playOnline += liveCount.TotalOnline
			webOnline += liveCount.IPConnect
		}
		if playOnline > 0 && webOnline > 0 {
			atomic.StoreInt64(&s.playOnline, playOnline)
			atomic.StoreInt64(&s.webOnline, webOnline)
		}
		time.Sleep(time.Duration(s.c.WEB.PullRegionInterval))
	}
}

func (s *Service) onlineListproc() {
	var (
		err  error
		aids []*model.OnlineAid
		arcs *arcmdl.ArcsReply
	)
	for {
		if aids, err = s.dao.OnlineList(context.Background(), _onlineListNum); err != nil {
			time.Sleep(time.Second)
			continue
		} else if len(aids) == 0 {
			log.Error("s.dao.OnlineList data len == 0")
			time.Sleep(time.Second)
			continue
		}
		var aidArg []int64
		for _, v := range aids {
			aidArg = append(aidArg, v.Aid)
		}
		archivesArgLog("onlineListproc", aidArg)
		if arcs, err = s.arcClient.Arcs(context.Background(), &arcmdl.ArcsRequest{Aids: aidArg}); err != nil {
			log.Error("s.arcClient.Arcs(%v) error (%v)", aidArg, err)
			time.Sleep(time.Second)
			continue
		} else {
			var onlineArcs []*model.OnlineArc
			for _, v := range aids {
				if arc, ok := arcs.Arcs[v.Aid]; ok && arc != nil && arc.IsNormal() && arc.AttrVal(archive.AttrBitNoRank) == archive.AttrNo {
					if arc.AttrVal(archive.AttrBitIsBangumi) == archive.AttrNo && arc.AttrVal(archive.AttrBitIsMovie) == archive.AttrNo {
						if time.Now().Unix()-int64(arc.PubDate) > _onlinePubdateLimit {
							if arc.Stat.Danmaku == 0 || v.Count/int64(arc.Stat.Danmaku) > _onlineDanmuLimit {
								continue
							}
						}
					}
					onlineArcs = append(onlineArcs, &model.OnlineArc{Arc: arc, OnlineCount: v.Count})
					if len(onlineArcs) >= s.c.WEB.OnlineCount {
						break
					}
				}
			}
			if len(onlineArcs) >= s.c.WEB.OnlineCount {
				s.onlineArcs = onlineArcs
				s.cache.Do(context.Background(), func(c context.Context) {
					s.dao.SetOnlineListBakCache(context.Background(), onlineArcs)
				})
			} else {
				log.Error("s.dao.OnlineList data len(%d) error", len(onlineArcs))
				time.Sleep(time.Second)
				continue
			}
		}
		time.Sleep(time.Duration(s.c.WEB.PullOnlineInterval))
	}
}
