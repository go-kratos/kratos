package service

import (
	"context"

	"github.com/robfig/cron"
	location "go-common/app/service/main/location/api"
	"go-common/app/service/video/stream-mng/conf"
	"go-common/app/service/video/stream-mng/dao"
	"go-common/library/net/rpc/warden"
	"go-common/library/sync/pipeline/fanout"
	xtime "go-common/library/time"
	"time"
)

// Service struct
type Service struct {
	c           *conf.Config
	dao         *dao.Dao
	locationCli location.LocationClient
	cron        *cron.Cron
	liveAside   *fanout.Fanout
}

// New init
func New(c *conf.Config) (s *Service) {
	cfg := &warden.ClientConfig{
		Dial:    xtime.Duration(time.Second * 1),
		Timeout: xtime.Duration(time.Second * 3),
	}
	locConn, err := warden.NewClient(cfg).Dial(context.Background(), "discovery://default/location.service")
	if err != nil {
		panic(err)
	}

	s = &Service{
		c:           c,
		dao:         dao.New(c),
		locationCli: location.NewLocationClient(locConn),
		cron:        cron.New(),
		liveAside:   fanout.New("stream-service", fanout.Worker(2), fanout.Buffer(1024)),
	}

	//if err := s.cron.AddFunc("0 */1 * * * *", s.refreshLiveStreamList); err != nil {
	//	panic(err)
	//}
	//
	//s.cron.Start()
	return s
}

// Ping Service
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close Service
func (s *Service) Close() {
	s.cron.Stop()
	s.dao.Close()
}
