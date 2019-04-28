package service

import (
	"context"

	"go-common/library/queue/databus"

	"go-common/app/service/live/xanchor/conf"
	"go-common/app/service/live/xanchor/dao"
	v1 "go-common/app/service/live/xanchor/service/v1"
)

// Service struct
type Service struct {
	c            *conf.Config
	dao          *dao.Dao
	v1svc        *v1.XAnchorService
	liveDanmuSub *databus.Databus
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:            c,
		dao:          dao.New(c),
		v1svc:        v1.NewXAnchorService(c),
		liveDanmuSub: databus.New(c.LiveDanmuSub),
	}
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

// V1Svc return v1 service
func (s *Service) V1Svc() *v1.XAnchorService {
	return s.v1svc
}
