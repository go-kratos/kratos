package service

import (
	"context"
	"sync"

	"go-common/app/interface/main/push/conf"
	"go-common/app/interface/main/push/dao"
	pushrpc "go-common/app/service/main/push/api/grpc/v1"
	pushmdl "go-common/app/service/main/push/model"
	"go-common/library/cache"
	httpx "go-common/library/net/http/blademaster"
)

// Service push service.
type Service struct {
	c          *conf.Config
	dao        *dao.Dao
	cache      *cache.Cache
	pushRPC    pushrpc.PushClient
	callbackCh chan *pushmdl.Callback
	httpClient *httpx.Client
	waiter     sync.WaitGroup
	closed     bool
}

// New creates a push service instance.
func New(c *conf.Config) *Service {
	s := &Service{
		c:          c,
		dao:        dao.New(c),
		cache:      cache.New(1, 10240),
		callbackCh: make(chan *pushmdl.Callback, c.Push.CallbackChanLen),
		httpClient: httpx.NewClient(c.HTTPClient),
	}
	var err error
	if s.pushRPC, err = pushrpc.NewClient(c.PushRPC); err != nil {
		panic(err)
	}
	for i := 0; i < s.c.Push.CallbackGoroutines; i++ {
		s.waiter.Add(1)
		go s.callbackproc()
	}
	return s
}

// Close closes service.
func (s *Service) Close() {
	s.closed = true
	close(s.callbackCh)
	s.waiter.Wait()
	s.dao.Close()
}

// Ping checks service.
func (s *Service) Ping(c context.Context) (err error) {
	err = s.dao.Ping(c)
	return
}
