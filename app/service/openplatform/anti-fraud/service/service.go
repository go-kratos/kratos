package service

import (
	"context"
	"math/rand"
	"time"

	"go-common/app/service/openplatform/anti-fraud/conf"
	"go-common/app/service/openplatform/anti-fraud/dao"
	"go-common/library/cache/redis"
)

// Service struct of service.
type Service struct {
	d     *dao.Dao
	c     *conf.Config // conf
	redis *redis.Pool
}

func init() {
	rand.Seed(time.Now().Unix())
}

// New create service instance and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:     c,
		d:     dao.New(c),
		redis: redis.NewPool(c.Redis.Config),
	}
	return
}

// Close dao.
func (s *Service) Close() {
	s.d.Close()
}

// Ping check server ok.
func (s *Service) Ping(c context.Context) (err error) {
	if err = s.d.Ping(c); err != nil {
		return
	}
	return
}
