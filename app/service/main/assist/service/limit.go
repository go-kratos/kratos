package service

import (
	"context"
	"go-common/library/ecode"
	"go-common/library/log"
)

func (s *Service) checkMaxAssistCnt(c context.Context, mid int64) (err error) {
	cnt, err := s.ass.AssistCnt(c, mid)
	if err != nil {
		log.Error("s.ass.AssistCnt(%d) error(%v)", mid, err)
		return
	}
	if cnt >= s.c.MaxAssCnt {
		err = ecode.AssistOverMaxLimit
		log.Error("ecode.AssistOverMaxLimit(%d) error(%v)", mid, err)
		return
	}
	return
}

func (s *Service) checkTotalLimit(c context.Context, mid int64) (err error) {
	cnt, err := s.ass.TotalAssCnt(c, mid)
	if err != nil {
		log.Error("s.ass.DailyCntAddAllAss(%d) error(%v)", mid, err)
		return
	}
	// 100
	if cnt >= 100 {
		err = ecode.AssistOverMaxLimitDailyAddAll
		log.Error("ecode.AssistOverMaxLimitDailyAddAll(%d) error(%v)", mid, err)
		return
	}
	return
}

func (s *Service) checkSameLimit(c context.Context, mid, assistMid int64) (err error) {
	cnt, err := s.ass.SameAssCnt(c, mid, assistMid)
	if err != nil {
		log.Error("s.ass.DailyCntAddSameAss(%d),(%d) error(%v)", mid, assistMid, err)
		return
	}
	// 2
	if cnt >= 2 {
		err = ecode.AssistOverMaxLimitDailyAddSame
		log.Error("ecode.AssistOverMaxLimitDailyAddSame(%d) error(%v)", mid, err)
		return
	}
	return
}
