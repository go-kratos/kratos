package service

import (
	"context"
	"math/rand"
	"time"

	"go-common/app/service/main/ugcpay/conf"
	"go-common/app/service/main/ugcpay/dao"
	"go-common/app/service/main/ugcpay/service/pay"
	"go-common/library/cache"
	"go-common/library/log"
)

// Service struct
type Service struct {
	c     *conf.Config
	dao   *dao.Dao
	pay   *pay.Pay
	cache *cache.Cache
	rnd   *rand.Rand
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:   c,
		dao: dao.New(c),
		pay: &pay.Pay{
			ID:              conf.Conf.Biz.Pay.ID,
			Token:           conf.Conf.Biz.Pay.Token,
			OrderTTL:        1800,
			NotifyURL:       conf.Conf.Biz.Pay.URLPayCallback,
			RefundNotifyURL: conf.Conf.Biz.Pay.URLRefundCallback,
		},
		cache: cache.New(10, 10240),
		rnd:   rand.New(rand.NewSource(time.Now().Unix())),
	}
	return s
}

// Ping Service
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close Service
func (s *Service) Close() {
	s.dao.Close()
}

func runCAS(ctx context.Context, fn func(ctx context.Context) (effected bool, err error)) (err error) {
	times := conf.Conf.Biz.RunCASTimes
	if times <= 0 {
		times = 3
	}
	effected := false
	for times > 0 {
		times--
		if effected, err = fn(ctx); err != nil {
			return
		}
		if effected {
			return
		}
	}
	if times <= 0 {
		log.Error("runCAS failed!!!")
	}
	return
}
