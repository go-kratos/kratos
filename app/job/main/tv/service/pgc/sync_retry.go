package pgc

import (
	"time"

	"go-common/app/job/main/tv/model/common"
	"go-common/app/job/main/tv/model/pgc"
	"go-common/library/log"
)

func (s *Service) addRetryEp(in *pgc.Content) {
	s.addRetryEps([]*pgc.Content{in})
}

// addRetryEps adds eps into retry list
func (s *Service) addRetryEps(in []*pgc.Content) {
	var (
		epids    []int
		newConts []*pgc.Content
	)
	for _, v := range in {
		if !s.retryLimit(false, int64(v.EPID)) { // filter retried too many times ep
			continue
		}
		newConts = append(newConts, v)
		epids = append(epids, v.EPID)
	}
	if len(newConts) == 0 {
		return
	}
	log.Warn("addRetryEps Add IDs %v", epids)
	s.ResuEps = append(s.ResuEps, newConts...)
}

// pickRetryEp picks the to-retry eps from memory
func (s *Service) pickRetryEp() (res []*pgc.Content) {
	if len(s.ResuEps) == 0 {
		return
	}
	res = append(res, s.ResuEps...)
	log.Info("pickRetry EP Len %d", len(res))
	s.ResuEps = make([]*pgc.Content, 0)
	return
}

// re-submit eps
func (s *Service) resubEps() {
	defer s.waiter.Done()
	for {
		if s.daoClosed {
			log.Info("resubEps DB closed!")
			return
		}
		readyEps := s.pickRetryEp() // pick to-retry eps from memory
		if len(readyEps) == 0 {
			log.Info("resubEps Empty")
			time.Sleep(time.Duration(s.c.Cfg.SyncRetry.RetryFre))
			continue
		}
		againEps := make([]*pgc.Content, 0)
		for _, ep := range readyEps { // retry them
			if err := s.epsSync(int64(ep.SeasonID), []*pgc.Content{ep}); err != nil { // if error, re-add this item into re-sub list
				log.Error("resubEps Sid %d, Epid %v, Err %v", ep.SeasonID, ep.EPID, err)
				againEps = append(againEps, ep)
				continue
			}
			retry := &common.SyncRetry{}
			retry.FromEp(0, int64(ep.EPID))
			s.dao.DelRetry(ctx, retry) // after succ, del it from MC
		}
		if len(againEps) > 0 {
			s.addRetryEps(againEps)
		}
		time.Sleep(time.Duration(s.c.Cfg.SyncRetry.RetryFre))
	}
}

func (s *Service) addRetrySn(in *pgc.TVEpSeason) {
	s.addRetrySns([]*pgc.TVEpSeason{in})
}

// addRetrySns adds sns into retry list
func (s *Service) addRetrySns(in []*pgc.TVEpSeason) {
	var (
		sids     []int64
		newConts []*pgc.TVEpSeason
	)
	for _, v := range in {
		if !s.retryLimit(true, v.ID) { // filter retried too many times ep
			continue
		}
		newConts = append(newConts, v)
		sids = append(sids, v.ID)
	}
	log.Warn("addRetrySns Add IDs %v", sids)
	s.ResuSns = append(s.ResuSns, newConts...)
}

// pickRetryEp picks the to-retry eps from memory
func (s *Service) pickRetrySn() (res []*pgc.TVEpSeason) {
	if len(s.ResuSns) == 0 {
		return
	}
	res = append(res, s.ResuSns...)
	log.Info("pickRetry Sn Len %d", len(res))
	s.ResuSns = make([]*pgc.TVEpSeason, 0)
	return
}

// re-submit eps
func (s *Service) resubSns() {
	defer s.waiter.Done()
	for {
		if s.daoClosed {
			log.Info("resubSns DB closed!")
			return
		}
		readySns := s.pickRetrySn()
		if len(readySns) == 0 {
			log.Info("resubSns Empty")
			time.Sleep(time.Duration(s.c.Cfg.SyncRetry.RetryFre))
			continue
		}
		againSns := make([]*pgc.TVEpSeason, 0)
		for _, sn := range readySns {
			if err := s.snSync(sn); err != nil { // if error, re-add this item into re-sub list
				log.Error("resubSns Sid %d, Err %v", sn.ID, err)
				againSns = append(againSns, sn)
				continue
			}
			retry := &common.SyncRetry{}
			retry.FromSn(0, sn.ID)
			s.dao.DelRetry(ctx, retry)
		}
		if len(againSns) > 0 {
			s.addRetrySns(againSns)
		}
		time.Sleep(time.Duration(s.c.Cfg.SyncRetry.RetryFre))
	}
}

// retryLimit limits the retry times
func (s *Service) retryLimit(isSn bool, id int64) bool {
	var req = &common.SyncRetry{}
	if isSn {
		req.FromSn(0, id)
	} else {
		req.FromEp(0, id)
	}
	retryTms, err := s.dao.GetRetry(ctx, req)
	if err != nil {
		log.Error("GetRetry Req %s, Err %v", req.MCKey(), err)
		return true
	}
	if retryTms > s.c.Cfg.SyncRetry.MaxRetry {
		log.Error("retryLimit Req %s, Retry Already %d times, stop here", req.MCKey(), retryTms)
		return false
	}
	s.dao.SetRetry(ctx, &common.SyncRetry{
		Ctype: req.Ctype,
		CID:   req.CID,
		Retry: retryTms + 1,
	})
	return true
}
