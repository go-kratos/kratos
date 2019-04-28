package service

import (
	"context"
	"math/rand"
	"time"

	"github.com/pkg/errors"

	"go-common/app/service/main/coupon/conf"
	"go-common/app/service/main/coupon/dao"
	"go-common/app/service/main/coupon/model"
	v1 "go-common/app/service/main/vipinfo/api"
	"go-common/library/log"
	"go-common/library/sync/pipeline/fanout"
)

var (
	_emptyCoupons          = make([]*model.CouponInfo, 0)
	_emptyBlance           = make([]*model.CouponBalanceInfo, 0)
	_emptyAllowance        = make([]*model.CouponAllowanceInfo, 0)
	_defps                 = 20
	_defpn                 = 1
	_deftitle              = "观影劵"
	_defCartoonTitle       = "漫画体验券"
	_maxCount        int64 = -1
)

// Service struct
type Service struct {
	c                    *conf.Config
	dao                  *dao.Dao
	r                    *rand.Rand
	messageChan          chan func()
	cache                *fanout.Fanout
	allBranchInfo        map[string]*model.CouponBatchInfo
	vipinfoClient        v1.VipInfoClient //vipinfo grpc client
	MapNoVipBatchToken   map[int8]string
	MapMonthBatchToken   map[int8]string
	MapMore180BatchToken map[int8]string
	MapLess180BatchToken map[int8]string
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:             c,
		dao:           dao.New(c),
		r:             rand.New(rand.NewSource(time.Now().Unix())),
		allBranchInfo: make(map[string]*model.CouponBatchInfo),
		messageChan:   make(chan func(), 10240),
		// cache chan
		cache: fanout.New("cache", fanout.Buffer(10240)),
		MapNoVipBatchToken: map[int8]string{
			model.CardType1:  c.NewYearConf.NoVipBatchToken1,
			model.CardType3:  c.NewYearConf.NoVipBatchToken3,
			model.CardType12: c.NewYearConf.NoVipBatchToken12,
		},
		MapMonthBatchToken: map[int8]string{
			model.CardType1:  c.NewYearConf.MonthBatchToken1,
			model.CardType3:  c.NewYearConf.MonthBatchToken3,
			model.CardType12: c.NewYearConf.MonthBatchToken12,
		},
		MapMore180BatchToken: map[int8]string{
			model.CardType1:  c.NewYearConf.More180BatchToken1,
			model.CardType3:  c.NewYearConf.More180BatchToken3,
			model.CardType12: c.NewYearConf.More180BatchToken12,
		},
		MapLess180BatchToken: map[int8]string{
			model.CardType1:  c.NewYearConf.Less180BatchToken1,
			model.CardType3:  c.NewYearConf.Less180BatchToken3,
			model.CardType12: c.NewYearConf.Less180BatchToken12,
		},
	}
	var err error
	if s.vipinfoClient, err = v1.NewClient(c.VipinfoRPC); err != nil {
		panic(errors.Wrap(err, "v1.NewClient failed"))
	}
	if err := s.loadbatchinfo(); err != nil {
		panic(err)
	}
	go s.loadbatchinfoproc()
	go s.handlermessageproc()
	return s
}

// Ping Service
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close Service
func (s *Service) Close() {
	s.dao.Close()
	s.cache.Close()
}

func (s *Service) loadbatchinfo() (err error) {
	var (
		c  = context.Background()
		bs []*model.CouponBatchInfo
	)
	if bs, err = s.dao.AllBranchInfo(c); err != nil {
		log.Error("loadbatchinfo allevent error(%v)", err)
		return
	}
	tmp := make(map[string]*model.CouponBatchInfo, len(bs))
	for _, v := range bs {
		tmp[v.BatchToken] = v
	}
	s.allBranchInfo = tmp
	log.Info("loadbatchinfo (%v) load success", tmp)
	return
}

func (s *Service) loadbatchinfoproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.loadbatchinfoproc panic(%v)", x)
			go s.loadbatchinfoproc()
			log.Info("service.loadbatchinfoproc recover")
		}
	}()
	for {
		time.Sleep(time.Minute * 1)
		s.loadbatchinfo()
	}
}

func (s *Service) sendMessage(f func()) {
	select {
	case s.messageChan <- f:
	default:
		log.Warn("message chan full")
	}
}

func (s *Service) handlermessageproc() {
	for {
		f := <-s.messageChan
		f()
	}
}
