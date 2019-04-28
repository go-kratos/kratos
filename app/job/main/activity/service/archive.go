package service

import (
	"context"
	"sort"
	"time"

	likemdl "go-common/app/interface/main/activity/model/like"
	"go-common/app/job/main/activity/model/like"
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/log"
)

const (
	_rankViewPieceSize = 100
	_rankCount         = 50
)

func (s *Service) subsRankproc() {
	for {
		if s.closed {
			return
		}
		var (
			subs []*like.Subject
			err  error
		)
		now := time.Now()
		if subs, err = s.dao.SubjectList(context.Background(), []int64{likemdl.PHONEVIDEO, likemdl.SMALLVIDEO}, now); err != nil {
			log.Error("viewRankproc s.dao.SubjectList error(%+v)", err)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		if len(subs) == 0 {
			log.Warn("viewRankproc no subjects time(%d)", now.Unix())
			time.Sleep(time.Duration(s.c.Interval.ViewRankInterval))
			continue
		}
		for _, v := range subs {
			s.viewRankproc(v.ID)
			time.Sleep(100 * time.Millisecond)
		}
		time.Sleep(time.Duration(s.c.Interval.ViewRankInterval))
	}
}

func (s *Service) viewRankproc(sid int64) {
	var (
		likeCnt  int
		likes    []*like.Like
		rankArcs []*api.Arc
		err      error
	)
	if likeCnt, err = s.dao.LikeCnt(context.Background(), sid); err != nil {
		log.Error("viewRankproc s.dao.LikeCnt(sid:%d) error(%v)", sid, err)
		return
	}
	if likeCnt == 0 {
		log.Warn("viewRankproc s.dao.LikeCnt(sid:%d) likeCnt == 0", sid)
		return
	}
	for i := 0; i < likeCnt; i += _rankViewPieceSize {
		if likes, err = s.likeList(context.Background(), sid, i, _objectPieceSize, _retryTimes); err != nil {
			log.Error("viewRankproc s.likeList(%d,%d,%d) error(%+v)", sid, i, _objectPieceSize, err)
			time.Sleep(100 * time.Millisecond)
			continue
		} else {
			var aids []int64
			for _, v := range likes {
				if v.Wid > 0 {
					aids = append(aids, v.Wid)
				}
			}
			var arcs map[int64]*api.Arc
			if arcs, err = s.arcs(context.Background(), aids, _retryTimes); err != nil {
				log.Error("viewRankproc s.arcs(%v) error(%v)", aids, err)
				time.Sleep(100 * time.Millisecond)
				continue
			} else {
				for _, aid := range aids {
					if arc, ok := arcs[aid]; ok && arc.IsNormal() {
						rankArcs = append(rankArcs, arc)
					}
				}
				sort.Slice(rankArcs, func(i, j int) bool {
					return rankArcs[i].Stat.View > rankArcs[j].Stat.View
				})
				if len(rankArcs) > _rankCount {
					rankArcs = rankArcs[:_rankCount]
				}
			}
		}
	}
	if len(rankArcs) > 0 {
		var rankAids []int64
		for _, v := range rankArcs {
			rankAids = append(rankAids, v.Aid)
		}
		if err = s.setViewRank(context.Background(), sid, rankAids, _retryTimes); err != nil {
			log.Error("viewRankproc s.setObjectStat(%d,%v) error(%+v)", sid, rankAids, err)
		}
	}
}

func (s *Service) arcs(c context.Context, aids []int64, retryCnt int) (arcs map[int64]*api.Arc, err error) {
	for i := 0; i < retryCnt; i++ {
		if arcs, err = s.arcRPC.Archives3(c, &archive.ArgAids2{Aids: aids}); err == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	return
}

func (s *Service) setViewRank(c context.Context, sid int64, aids []int64, retryTime int) (err error) {
	for i := 0; i < retryTime; i++ {
		if err = s.dao.SetViewRank(c, sid, aids); err == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	return
}
