package service

import (
	"context"
	"sync"

	"go-common/app/service/main/thumbup/conf"
	"go-common/app/service/main/thumbup/dao"
	"go-common/library/sync/pipeline/fanout"
)

// Service service
type Service struct {
	c      *conf.Config
	dao    *dao.Dao
	cache  *fanout.Fanout
	dbus   *fanout.Fanout
	waiter *sync.WaitGroup
	close  bool
}

// New new
func New(c *conf.Config) *Service {
	s := &Service{
		c:      c,
		dao:    dao.New(c),
		cache:  fanout.New("cache", fanout.Buffer(10240)),
		dbus:   fanout.New("dbus"),
		waiter: new(sync.WaitGroup),
	}
	return s
}

// Close close dao.
func (s *Service) Close() {
	s.cache.Close()
	s.dbus.Close()
	s.dao.Close()
	s.close = true
	s.waiter.Wait()
}

// Ping check connection success.
func (s *Service) Ping(c context.Context) (err error) {
	err = s.dao.Ping(c)
	return
}
