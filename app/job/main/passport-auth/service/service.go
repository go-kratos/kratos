package service

import (
	"context"

	"go-common/app/job/main/passport-auth/conf"
	"go-common/app/job/main/passport-auth/dao"
	auth "go-common/app/service/main/passport-auth/rpc/client"
	"go-common/library/queue/databus"
	"go-common/library/queue/databus/databusutil"
)

// Service struct
type Service struct {
	c               *conf.Config
	dao             *dao.Dao
	g               *databusutil.Group
	oldAuthConsumer *databus.Databus
	authRPC         *auth.Service
	authConsumer    *databus.Databus
	authGroup       *databusutil.Group
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:               c,
		dao:             dao.New(c),
		oldAuthConsumer: databus.New(c.Databus),
		authRPC:         auth.New(c.AuthRPC),
		authConsumer:    databus.New(c.AuthDataBus),
	}
	// new a group
	s.g = databusutil.NewGroup(
		c.DatabusUtil,
		s.oldAuthConsumer.Messages(),
	)
	s.authGroup = databusutil.NewGroup(
		c.DatabusUtil,
		s.authConsumer.Messages(),
	)
	s.consumeproc()
	s.authConsumeProc()
	// go s.syncCookie()
	// for i := c.IDXFrom; i < c.IDXTo; i ++ {
	// 	go s.syncCookie(int64(i))
	// }

	// go s.syncToken("201804", 0, 50000000)
	// go s.syncToken("201804", 50000001, 100000000)
	// go s.syncToken("201804", 100000001, 150000000)
	// go s.syncToken("201804", 150000001, 0)

	return s
}

// Ping Service
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close Service
func (s *Service) Close() {
	s.g.Close()
	s.authGroup.Close()
	s.dao.Close()
}
