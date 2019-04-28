package service

import (
	"context"

	"go-common/app/job/main/mcn/model"
	"go-common/library/log"

	"github.com/pkg/errors"
)

// McnRecommendCron .
func (s *Service) McnRecommendCron() {
	defer func() {
		if r := recover(); r != nil {
			r = errors.WithStack(r.(error))
			log.Error("recover panic  error(%+v)", r)
		}
	}()
	var (
		err   error
		limit = 100
		c     = context.TODO()
		rps   []*model.McnUpRecommendPool
	)
	for {
		if rps, err = s.dao.McnUpRecommendSources(c, limit); err != nil {
			log.Error("s.dao.McnUpRecommendSources(%d) error(%+v)", limit, err)
			return
		}
		if len(rps) == 0 {
			log.Warn("big data recommend up data is empty!")
			return
		}
		for _, v := range rps {
			rp := new(model.McnUpRecommendPool)
			*rp = *v
			rp.GenerateTime = v.Mtime
			if _, err = s.dao.AddMcnUpRecommend(c, rp); err != nil {
				log.Error("s.dao.AddMcnUpRecommend(%+v) error(%+v)", rp, err)
				continue
			}
			if _, err = s.dao.DelMcnUpRecommendSource(c, v.ID); err != nil {
				log.Error("s.dao.DelMcnUpRecommendSource(%d) error(%+v)", v.ID, err)
				continue
			}
			log.Info("source id(%d) sync to recommend poll(%+v)", v.ID, rp)
		}
	}
}

// DealFailRecommendCron .
func (s *Service) DealFailRecommendCron() {
	defer func() {
		if r := recover(); r != nil {
			r = errors.WithStack(r.(error))
			log.Error("recover panic  error(%+v)", r)
		}
	}()
	var (
		err error
		c   = context.TODO()
	)
	if _, err = s.dao.DelMcnUpRecommendPool(c); err != nil {
		log.Error("s.dao.DelMcnUpRecommendPool error(%+v)", err)
	}
}
