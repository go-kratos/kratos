package staff

import (
	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/service"
)

type Service struct {
	c *conf.Config
}

func New(c *conf.Config, rpcdaos *service.RPCDaos) *Service {
	s := &Service{
		c: c,
	}
	return s
}
func (s *Service) Config() (conf *conf.StaffConf) {
	conf = s.c.StaffConf
	return
}
func (s *Service) TypeConfig() (conf []*conf.StaffTypeConf) {
	conf = s.c.StaffConf.TypeList
	return
}
