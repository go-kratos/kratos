package kfc

import (
	"context"

	"go-common/app/interface/main/activity/conf"
	"go-common/app/interface/main/activity/dao/kfc"
	kfcmdl "go-common/app/interface/main/activity/model/kfc"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/pipeline/fanout"
)

// Service struct
type Service struct {
	c     *conf.Config
	dao   *kfc.Dao
	cache *fanout.Fanout
}

// Close service
func (s *Service) Close() {
	s.dao.Close()
	s.cache.Close()
}

// New Service
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:     c,
		dao:   kfc.New(c),
		cache: fanout.New("cache", fanout.Worker(5), fanout.Buffer(10240)),
	}
	return
}

// KfcInfo .
func (s *Service) KfcInfo(c context.Context, id, mid int64) (kfcInfo *kfcmdl.BnjKfcCoupon, err error) {
	var (
		res *kfcmdl.BnjKfcCoupon
	)
	if res, err = s.dao.KfcCoupon(c, id); err != nil {
		log.Error("s.dao.KfcCoupon(%d) error(%+v)", id, err)
		return
	}
	if res.ID == 0 {
		err = ecode.NothingFound
		return
	}
	if res.Mid == 0 {
		var uid int64
		if uid, err = s.kfcRecall(c, id); err == nil && uid > 0 {
			if uid == mid {
				res.Mid = uid
			}
		}
	}
	if res.Mid != 0 && res.Mid == mid {
		kfcInfo = res
	} else {
		err = ecode.NothingFound
	}
	return
}

// kfcRecall .
func (s *Service) kfcRecall(c context.Context, id int64) (uid int64, err error) {
	if uid, err = s.dao.KfcWinner(c, id); err != nil {
		log.Error("s.dao.KfcWinner(%d) error(%+v)", id, err)
		return
	}
	if uid > 0 {
		s.DeliverKfc(c, id, uid)
	}
	return
}

// KfcUse .
func (s *Service) KfcUse(c context.Context, code string) (kfcID int64, err error) {
	var (
		kfcInfo    *kfcmdl.BnjKfcCoupon
		effectRows int64
	)
	if kfcID, err = s.dao.CacheKfcCode(c, code); err != nil {
		log.Error("s.dao.CacheKfcCode(%s) error(%+v)", code, err)
		return
	}
	if kfcID > 0 {
		err = ecode.ActivityKfcHasUsed
		return
	}
	if kfcInfo, err = s.dao.RawKfcCode(c, code); err != nil {
		log.Error("s.dao.RawKfcCode(%s) error(%+v)", code, err)
		return
	}
	if kfcInfo.ID == 0 {
		err = ecode.ActivityKfcNotExist
		return
	}
	if kfcInfo.Mid == 0 {
		err = ecode.ActivityKfcNotGiveOut
		return
	}
	if kfcInfo.State == int64(kfc.KfcCodeUsed) {
		s.cache.Do(c, func(c context.Context) {
			s.dao.AddCacheKfcCode(c, code, kfcInfo.ID)
		})
		err = ecode.ActivityKfcHasUsed
		return
	}
	if effectRows, err = s.dao.KfcCodeGiveOut(c, kfcInfo.ID); err != nil {
		log.Error("s.dao.KfcCodeGiveOut(%d) error(%+v)", kfcInfo.ID, err)
		return
	}
	if effectRows == 0 {
		err = ecode.ActivityKfcSqlError
	}
	kfcID = kfcInfo.ID
	s.cache.Do(c, func(c context.Context) {
		s.dao.AddCacheKfcCode(c, code, kfcInfo.ID)
		s.dao.DelCacheKfcCoupon(c, kfcInfo.ID)
	})
	return
}

// DeliverKfc .
func (s *Service) DeliverKfc(c context.Context, id, mid int64) (err error) {
	effectID, err := s.dao.KfcDeliver(c, id, mid)
	if err != nil {
		log.Error("s.dao.KfcDeliver(%d,%d) error(%+v)", id, mid, err)
		return
	}
	if effectID > 0 {
		s.cache.Do(c, func(c context.Context) {
			if e := s.dao.DelCacheKfcCoupon(c, id); e == nil {
				s.dao.KfcCoupon(c, id)
			}
		})
	} else {
		err = ecode.ActivityKfcSqlError
		log.Error("DeliverKfc mysql effect o rows (%d,%d)", id, mid)
	}
	return
}
