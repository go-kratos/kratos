package up

import (
	"context"

	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/dao/account"
	"go-common/app/interface/main/creative/dao/up"
	"go-common/app/interface/main/creative/service"
	"go-common/library/log"
)

//Service struct
type Service struct {
	c   *conf.Config
	up  *up.Dao
	acc *account.Dao
}

//New get service
func New(c *conf.Config, rpcdaos *service.RPCDaos) *Service {
	s := &Service{
		c:   c,
		up:  up.New(c),
		acc: rpcdaos.Acc,
	}
	return s
}

// Ping service
func (s *Service) Ping(c context.Context) (err error) {
	if err = s.up.Ping(c); err != nil {
		log.Error("s.up.Ping err(%v)", err)
	}
	return
}
