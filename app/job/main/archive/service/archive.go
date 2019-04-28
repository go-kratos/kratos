package service

import (
	"context"

	"go-common/app/job/main/archive/model/archive"
	"go-common/app/job/main/archive/model/result"
	"go-common/app/job/main/archive/model/retry"
	arcmdl "go-common/app/service/main/archive/model/archive"
	"go-common/library/log"
)

func (s *Service) isPGC(aid int64) bool {
	if addit, _ := s.archiveDao.Addit(context.TODO(), aid); addit != nil && (addit.UpFrom == archive.UpFromPGC || addit.UpFrom == archive.UpFromPGCSecret) {
		return true
	}
	return false
}

func (s *Service) consumerVideoup(i int) {
	defer s.waiter.Done()
	for {
		var (
			aid int64
			ok  bool
		)
		if aid, ok = <-s.videoupAids[i]; !ok {
			log.Error("s.videoupAids chan closed")
			return
		}
		arc, _ := s.arcServices[0].Archive3(context.TODO(), &arcmdl.ArgAid2{Aid: aid})
		if arc != nil && (arc.AttrVal(arcmdl.AttrBitIsPGC) == arcmdl.AttrYes || arc.AttrVal(arcmdl.AttrBitIsBangumi) == arcmdl.AttrYes) {
			if s.c.PGCAsync == 1 {
				rt := &retry.Info{Action: retry.FailResultAdd}
				rt.Data.Aid = aid
				s.PushFail(context.TODO(), rt)
				log.Warn("async PGC archive(%d)", aid)
				continue
			}
			s.pgcAids <- aid
			log.Info("aid(%d) title(%s) is PGC", aid, arc.Title)
			continue
		}
		if s.c.UGCAsync == 1 {
			rt := &retry.Info{Action: retry.FailResultAdd}
			rt.Data.Aid = aid
			s.PushFail(context.TODO(), rt)
			log.Warn("async UGC archive(%d)", aid)
			continue
		}
		s.arcUpdate(aid)
	}
}

func (s *Service) pgcConsumer() {
	defer s.waiter.Done()
	for {
		var (
			aid int64
			ok  bool
		)
		if aid, ok = <-s.pgcAids; !ok {
			log.Error("s.pgcAids closed")
			return
		}
		s.arcUpdate(aid)
	}
}

func (s *Service) arcUpdate(aid int64) {
	var (
		oldResult *result.Archive
		newResult *result.Archive
		c         = context.TODO()
		upCids    []int64
		delCids   []int64
		err       error
		changed   bool
	)
	log.Info("sync resultDB archive(%d) start", aid)
	defer func() {
		if err != nil {
			if oldResult != nil && (oldResult.AttrVal(result.AttrBitIsBangumi) == result.AttrYes || oldResult.AttrVal(result.AttrBitIsPGC) == result.AttrYes) {
				s.pgcAids <- aid
			} else {
				s.videoupAids[aid%int64(s.c.ChanSize)] <- aid
			}
			log.Error("s.arcUpdate(%d) error(%v)", aid, err)
		}
	}()
	if oldResult, err = s.resultDao.Archive(c, aid); err != nil {
		log.Error("s.resultDao.Archive(%d) error(%v)", aid, err)
	}
	if changed, upCids, delCids, err = s.tranResult(c, aid); err != nil || !changed {
		log.Error("aid(%d) nothing changed err(%+v)", aid, err)
		err = nil
		return
	}
	s.upVideoCache(aid, upCids)
	s.delVideoCache(aid, delCids)
	if newResult, err = s.resultDao.Archive(c, aid); err != nil {
		log.Error("s.resultDao.Archive(%d) error(%v)", aid, err)
		return
	}
	err = s.updateResultCache(newResult, oldResult)
	if oldResult != nil {
		s.updateResultField(newResult, oldResult)
		s.updateSubjectMid(newResult, oldResult)
		s.sendMail(newResult, oldResult)
	}
	action := "update"
	if oldResult == nil {
		action = "insert"
	}
	s.sendNotify(&result.ArchiveUpInfo{Table: "archive", Action: action, Nw: newResult, Old: oldResult})
	if oldResult != nil {
		log.Info("sync resultDB archive(%d) sync old(%+v) new(%+v) updated", aid, oldResult, newResult)
		return
	}
	log.Info("sync resultDB archive(%d) new(%+v) inserted", aid, newResult)
}

func (s *Service) hadPassed(c context.Context, aid int64) (had bool) {
	id, err := s.archiveDao.GetFirstPassByAID(c, aid)
	if err != nil {
		log.Error("hadPassed s.arc.GetFirstPassByAID error(%v) aid(%d)", err, aid)
		return
	}
	had = id > 0
	return
}
