package service

import (
	"context"

	"go-common/app/service/main/coin/conf"
	"go-common/app/service/main/coin/dao"
	"go-common/app/service/main/coin/model"
	memrpc "go-common/app/service/main/member/api"
	"go-common/library/ecode"
	"go-common/library/sync/pipeline/fanout"
)

// Service define service.
type Service struct {
	c             *conf.Config
	coinDao       *dao.Dao
	memRPC        memrpc.MemberClient
	cache         *fanout.Fanout
	job           *fanout.Fanout
	businesses    map[int64]*model.Business
	businessNames map[string]*model.Business
	statMerge     *statMerge
}

type statMerge struct {
	Business string
	Target   int64
	Sources  map[int64]bool
}

// New new a Service and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:       c,
		coinDao: dao.New(c),
		cache:   fanout.New("service-cache", fanout.Buffer(10240)),
		job:     fanout.New("job", fanout.Worker(10), fanout.Buffer(10240)),
	}
	var err error
	if s.memRPC, err = memrpc.NewClient(c.MemberRPC); err != nil {
		panic(err)
	}
	s.businesses = s.coinDao.Businesses
	s.businessNames = s.coinDao.BusinessNames
	if c.StatMerge != nil {
		s.statMerge = &statMerge{
			Business: c.StatMerge.Business,
			Target:   c.StatMerge.Target,
			Sources:  make(map[int64]bool),
		}
		for _, id := range c.StatMerge.Sources {
			s.statMerge.Sources[id] = true
		}
	}
	return
}

// Ping check service health.
func (s *Service) Ping(c context.Context) (err error) {
	err = s.coinDao.Ping(c)
	return
}

// Round 保留两位小数
func Round(x float64) float64 {
	n := -0.5
	if x > 0 {
		n = 0.5
	}
	return float64(int64(x/0.01+n)) / 100.0
}

// CheckBusiness .
func (s *Service) CheckBusiness(bs string) (id int64, err error) {
	if bs == "" {
		return
	}
	b := s.businessNames[bs]
	if b == nil {
		err = ecode.AppDenied
		return
	}
	id = b.ID
	return
}

// MustCheckBusiness .
// +wd:ignore
func (s *Service) MustCheckBusiness(bs string) (id int64, err error) {
	if bs == "" {
		err = ecode.AppDenied
		return
	}
	return s.CheckBusiness(bs)
}

// GetBusinessName .
func (s *Service) GetBusinessName(id int64) (res string, err error) {
	b := s.businesses[id]
	if b == nil {
		err = ecode.AppDenied
		return
	}
	return b.Name, nil
}

// Close .
func (s *Service) Close() (err error) {
	s.job.Close()
	s.cache.Close()
	return
}
