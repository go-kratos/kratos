package service

import (
	"context"
	"sync/atomic"
	"time"

	arcrpc "go-common/app/service/main/archive/api/gorpc"
	"go-common/app/service/main/tag/conf"
	"go-common/app/service/main/tag/dao"
	"go-common/app/service/main/tag/model"
	"go-common/library/cache"
)

// Service service.
type Service struct {
	conf    *conf.Config
	dao     *dao.Dao
	cache   *cache.Cache
	cacheCh *cache.Cache
	arcRPC  *arcrpc.Service2
	// memery cache
	limitRes     map[int64]int8     // limitRes
	whiteUser    map[int64]struct{} // superUser
	whiteUserMap atomic.Value
}

// New new a service and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		conf:      c,
		dao:       dao.New(c),
		cache:     cache.New(1, 1024),
		limitRes:  make(map[int64]int8),
		whiteUser: make(map[int64]struct{}),
		arcRPC:    arcrpc.New2(c.RPC.Archive),
		cacheCh:   cache.New(1, 1024),
	}
	go s.whiteUserproc()
	return s
}

// Ping ping dao .
func (s *Service) Ping(c context.Context) error {
	return s.dao.Ping(c)
}

// Close close all dao.
func (s *Service) Close() {
	s.dao.Close()
}

// LimitResource .
func (s *Service) LimitResource(c context.Context) ([]*model.ResourceLimit, error) {
	return s.dao.LimitRes(c, model.ResTypeArchive)
}

func (s *Service) whiteUserproc() {
	for {
		userMap, err := s.WhiteUser(context.TODO())
		if err == nil {
			s.whiteUserMap.Store(userMap)
		}
		time.Sleep(time.Minute * 10)
	}
}

// WhiteUser .
func (s *Service) WhiteUser(c context.Context) (map[int64]struct{}, error) {
	return s.dao.WhiteUser(c)
}

// RankingHot .
func (s *Service) RankingHot(c context.Context) ([]*model.Tag, error) {
	return s.dao.RankHots(c)
}

// RankingBangumi .
func (s *Service) RankingBangumi(c context.Context) ([]int64, []*model.RankingBangumi, error) {
	return s.dao.Bangumis(c)
}

// RankingRegion .
func (s *Service) RankingRegion(c context.Context, rid int64) ([]*model.RankingRegion, error) {
	return s.dao.Regions(c, rid)
}
