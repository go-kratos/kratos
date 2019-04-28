package service

import (
	"context"
	"math/rand"
	"time"

	"go-common/app/admin/main/coupon/conf"
	"go-common/app/admin/main/coupon/dao"
	"go-common/app/admin/main/coupon/model"
	courpc "go-common/app/service/main/coupon/rpc/client"
	"go-common/library/log"
	"go-common/library/sync/pipeline/fanout"
)

const (
	_maxSalaryCount = 100000
	_notLimitSalary = -1
	_maxretry       = 3
	_lockseconds    = 604800
)

// Service struct
type Service struct {
	c          *conf.Config
	dao        *dao.Dao
	r          *rand.Rand
	allAppInfo map[int64]string
	couRPC     *courpc.Service
	// cache async del
	cache *fanout.Fanout
	// msg async send
	msgchan *fanout.Fanout
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:          c,
		dao:        dao.New(c),
		r:          rand.New(rand.NewSource(time.Now().Unix())),
		allAppInfo: make(map[int64]string),
		couRPC:     courpc.New(c.RPCClient2.Coupon),
		// cache chan
		cache: fanout.New("cache", fanout.Worker(5), fanout.Buffer(1024)),
		// msg chan
		msgchan: fanout.New("cache", fanout.Worker(5), fanout.Buffer(10240)),
	}
	if err := s.loadappinfo(); err != nil {
		panic(err)
	}
	go s.loadappinfoproc()
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

func (s *Service) loadappinfo() (err error) {
	var (
		c  = context.Background()
		as []*model.AppInfo
	)
	if as, err = s.dao.AllAppInfo(c); err != nil {
		log.Error("loadappinfo allappinfo error(%v)", err)
		return
	}
	tmp := make(map[int64]string, len(as))
	for _, v := range as {
		tmp[v.ID] = v.Name
	}
	s.allAppInfo = tmp
	log.Info("loadappinfo (%v) load success", tmp)
	return
}

func (s *Service) loadappinfoproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.loadappinfoproc panic(%v)", x)
			go s.loadappinfoproc()
			log.Info("service.loadappinfoproc recover")
		}
	}()
	for {
		time.Sleep(time.Minute * 2)
		s.loadappinfo()
	}
}
