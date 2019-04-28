package kfc

import (
	"context"

	"go-common/app/admin/main/activity/conf"
	kfcDao "go-common/app/admin/main/activity/dao/kfc"
	kfcmdl "go-common/app/admin/main/activity/model/kfc"
	"go-common/library/log"
)

// Service struct
type Service struct {
	c   *conf.Config
	dao *kfcDao.Dao
}

// Close service
func (s *Service) Close() {
	if s.dao != nil {
		s.dao.Close()
	}
}

// New Service
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:   c,
		dao: kfcDao.New(c),
	}
	return
}

// List .
func (s *Service) List(c context.Context, arg *kfcmdl.ListParams) (list []*kfcmdl.BnjKfcCoupon, err error) {
	if list, err = s.dao.SearchList(c, arg.CouponCode, arg.Mid, arg.Pn, arg.Ps); err != nil {
		log.Error("s.dao.SearchList(%v) error(%+v)", arg, err)
	}
	return
}
