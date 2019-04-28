package service

import (
	"context"

	"go-common/app/admin/main/tag/conf"
	"go-common/app/admin/main/tag/dao"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/sync/pipeline/fanout"
)

// Service service struct.
type Service struct {
	conf    *conf.Config
	dao     *dao.Dao
	client  *bm.Client
	cacheCh *fanout.Fanout
}

// New new a service and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		conf:    c,
		dao:     dao.New(c),
		client:  bm.NewClient(c.HTTPClient),
		cacheCh: fanout.New("cache", fanout.Worker(1), fanout.Buffer(1024)),
	}
	return s
}

//Ping check the service health.
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close when program have recived SIGQUIT, SIGTERM, SIGSTOP or SIGINT signal, stop and relese all resource.
func (s *Service) Close() (err error) {
	return s.dao.Close(context.TODO())
}
