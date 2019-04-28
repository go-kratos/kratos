package service

import (
	"context"
	"time"

	"go-common/app/service/main/account/conf"
	"go-common/app/service/main/account/dao"
	"go-common/app/service/main/account/model/queue"
	"go-common/app/service/main/coin/api/gorpc"
	"go-common/app/service/main/relation/rpc/client"
	mc "go-common/library/cache/memcache"
	"go-common/library/log"

	"github.com/pkg/errors"
)

// Service http service
type Service struct {
	c       *conf.Config
	dao     *dao.Dao
	relRPC  *relation.Service
	coinRPC *coin.Service
	cachepq *queue.PriorityQueue
}

// New for new service obj
func New(c *conf.Config) *Service {
	s := &Service{
		c:       c,
		dao:     dao.New(c),
		relRPC:  relation.New(c.RelationRPC),
		coinRPC: coin.New(c.CoinRPC),
		cachepq: queue.NewPriorityQueue(2048, false),
	}
	go s.cachedelayproc(context.Background())
	return s
}

// Ping check server ok
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close dao
func (s *Service) Close() {
	s.dao.Close()
}

// DelCache for del user and access cahe
func (s *Service) DelCache(c context.Context, mid int64, action string) error {
	defer func() {
		item := &Item{
			Mid:    mid,
			Time:   time.Now(),
			Action: action,
		}
		switch action {
		case "updateVip":
			if err := s.cachepq.Put(item); err != nil {
				log.Warn("Failed to enqueue cache delay item: %+v: %+v", item, err)
			}
		}
	}()

	errs := s.dao.DelCache(context.TODO(), mid)
	for _, e := range errs {
		if errors.Cause(e) == mc.ErrNotFound {
			continue
		}
		return e
	}
	return nil
}
