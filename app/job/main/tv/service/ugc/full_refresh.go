package ugc

import (
	"context"
	"fmt"
	"time"

	"go-common/app/job/main/tv/dao/app"
	"go-common/app/job/main/tv/model/ugc"
	"go-common/library/log"
)

const (
	_errSleep  = 500 * time.Millisecond
	_succSleep = 10 * time.Millisecond
)

func errMid(funcName string, mid int64, err error) {
	log.Error("Func:[%s], Step:[%s], Mid:[%d], Err:[%v]", "fullRefresh", funcName, mid, err)
}

func errArcPce(funcName string, mid int64, numPce int, err error) {
	log.Error("Func:[%s], Step:[%s], Mid:[%d], NumPce:[%d], Err:[%v]", "fullRefresh-ArcPce", funcName, mid, numPce, err)
}

func infoArc(funcName string, aid int64, msg string) {
	log.Info("Func:[%s], Step:[%s], Aid:[%d], Msg:[%s]", "fullRefresh-ArcPce-Arc", funcName, aid, msg)
	time.Sleep(_errSleep)
}

func errArc(funcName string, aid int64, err error) {
	log.Error("Func:[%s], Step:[%s], Aid:[%d], Err:[%v]", "fullRefresh-ArcPce-Arc", funcName, aid, err)
	time.Sleep(_errSleep)
}

func errArcVideos(funcName string, aid int64, cids []int64, err error) {
	log.Error("Func:[%s], Step:[%s], Aid:[%d], Cids: [%v], Err:[%v]", "fullRefresh-ArcPce-Arc-Videos", funcName, aid, cids, err)
	time.Sleep(_errSleep)
}

func infoArcVideos(funcName string, aid int64, cids []int64, msg string) {
	log.Info("Func:[%s], Step:[%s], Aid:[%d], Cids: [%v], Msg:[%s]", "fullRefresh-ArcPce-Arc-Videos", funcName, aid, cids, msg)
	time.Sleep(_succSleep)
}

func (s *Service) fullRefreshproc() {
	for {
		s.fullRefresh()
		time.Sleep(time.Duration(s.c.UgcSync.Frequency.FullRefreshFre))
	}
}

func (s *Service) fullRefresh() {
	var (
		fullName  = "fullRefresh"
		pagesize  = s.c.UgcSync.Batch.ProducerPS
		begin     = time.Now()
		totalArcs = 0
		treatedUp = 0
		totalUp   = len(s.activeUps)
	)
	if totalUp == 0 {
		log.Error("[%s] ActiveUps Empty", fullName)
		return
	}
	log.Info("fullRefresh Total Uppers Len %d", totalUp)
	for mid := range s.activeUps {
		var (
			upArcCnt int
			err      error
		)
		if upArcCnt, err = s.dao.UpArcsCnt(ctx, int64(mid)); err != nil {
			errMid("CountUpArcs", mid, err)
			continue
		}
		if upArcCnt == 0 {
			errMid("CountUpArcs", mid, fmt.Errorf("Empty Arcs"))
			continue
		}
		for arcPce := 0; arcPce < app.NumPce(int(upArcCnt), pagesize); arcPce++ { // travel the upper's archive by piece
			var upArcs []*ugc.ArcFull
			if upArcs, err = s.dao.PickUpArcs(ctx, int(mid), arcPce, pagesize); err != nil {
				errArcPce("PickUpArcs", mid, arcPce, err)
				continue
			}
			if len(upArcs) == 0 {
				errArcPce("PickUpArcs", mid, arcPce, fmt.Errorf("Empty Arcs, Stop Picking"))
				break
			}
			if err = s.fullArcs(ctx, upArcs); err != nil {
				errArcPce("FullArcs", mid, arcPce, err)
			}
			time.Sleep(time.Duration(s.c.UgcSync.Frequency.FullRefArcFre)) // pause between each archives pce treatment
		}
		treatedUp = treatedUp + 1
		totalArcs = totalArcs + upArcCnt
		log.Info("fullRefresh Total Up %d, Treated Up %d, Treated Arcs %d, Time Used %v", totalUp, treatedUp, totalArcs, time.Since(begin))
	}
	log.Info("fullRefresh Ends! Len Uppers %d, Time Used %v", len(s.activeUps), time.Since(begin))
}

func (s *Service) fullArcs(ctx context.Context, arcs []*ugc.ArcFull) (err error) {
	for _, arc := range arcs {
		var (
			arcOk, actVideos, shouldAudit bool
			aid                           = arc.AID
			transCids                     []int64
			arcAllow                      = &ugc.ArcAllow{}
		)
		if err = s.dao.SetArcCMS(ctx, &arc.ArcCMS); err != nil { // set cache
			errArc("SetArcCMS", aid, err) // cache error, ignore
		}
		if arc.Deleted == 1 {
			if actVideos, err = s.dao.ActVideos(ctx, aid); err != nil {
				errArc("actVideos", aid, err) // db error
				continue
			}
			if !actVideos {
				infoArc("actVideos", aid, "Arc Deleted && No Active Videos, Jump to the next")
				continue
			} else {
				if err = s.dao.DelVideos(ctx, aid); err != nil { // delete also the videos
					errArc("actVideos", aid, err)
					continue
				}
				infoArc("actVideos", aid, "Arc Deleted, So we delete the rest videos")
			}
		}
		arcAllow.FromArcFull(arc)
		if arcOk = s.arcAllowImport(arcAllow); !arcOk {
			log.Warn("[fullRefresh-ArcPce-Arc]")
			continue
		}
		if arcOk, transCids, err = s.transFailTreat(ctx, aid); err != nil {
			errArcVideos("TransFailVideos-DelVideos", aid, transCids, err) // db error
			continue
		}
		if !arcOk {
			continue
		}
		if shouldAudit, err = s.dao.ShouldAudit(ctx, aid); err != nil {
			errArc("ShouldAudit", aid, err)
			continue
		}
		if shouldAudit {
			log.Info("fullRefresh addAudCid cAid %d", aid)
			s.audAidCh <- []int64{aid} // add aid into channel to treat
		}
		if err = s.refArcVideo(ctx, aid); err != nil {
			errArc("refArcVideo", aid, err)
			continue
		}
		time.Sleep(10 * time.Millisecond)
	}
	return
}

func (s *Service) transFailTreat(ctx context.Context, aid int64) (arcOk bool, failCids []int64, err error) {
	arcOk = true
	if failCids, err = s.dao.TransFailVideos(ctx, aid); err != nil { // delete transcoding failed cids
		errArc("TransFailVideos", aid, err) // db error, stop this archive here
		return
	}
	if len(failCids) == 0 {
		// infoArcVideos("TransFailVideos", aid, failCids, "No Fail Cids")
		return
	}
	if arcOk, err = s.dao.DelVideoArc(ctx, &ugc.DelVideos{
		AID:  aid,
		CIDs: failCids,
	}); err != nil {
		return
	}
	if !arcOk {
		infoArcVideos("TransFailVideos", aid, failCids, " Delete Videos & Arc succ")
		return
	}
	infoArcVideos("TransFailVideos", aid, failCids, " Delete Videos succ")
	return
}

func (s *Service) refArcVideo(ctx context.Context, cAid int64) (err error) {
	var (
		proName  = "videoProducer-video"
		pagesize = s.c.UgcSync.Batch.ProducerPS
		videoCnt int
		maxID    = 0
	)
	if videoCnt, err = s.dao.ArcVideoCnt(ctx, cAid); err != nil {
		log.Error("[%s] CountArcs Aid %d, error [%v]", proName, cAid, err)
		return
	}
	if videoCnt == 0 {
		return
	}
	nbPiece := app.NumPce(videoCnt, pagesize)
	log.Info("[%s] NumPiece %d, Pagesize %d", proName, nbPiece, pagesize)
	for i := 0; i < nbPiece; i++ {
		videos, newMaxID, errR := s.dao.PickArcVideo(ctx, cAid, maxID, pagesize)
		if errR != nil {
			log.Error("[%s] Pick Piece %d Error, Ignore it", proName, i)
			continue
		}
		if newMaxID <= maxID {
			log.Error("[%s] MaxID is not increasing! [%d,%d]", proName, newMaxID, maxID)
			return
		}
		maxID = newMaxID
		for _, v := range videos {
			s.dao.SetVideoCMS(ctx, v)
		}
		time.Sleep(500 * time.Millisecond)
	}
	return
}
