package template

import (
	"context"

	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/dao/template"
	"go-common/app/interface/main/creative/service"
	"go-common/library/log"
)

//Service struct
type Service struct {
	c   *conf.Config
	tpl *template.Dao
}

//New get service
func New(c *conf.Config, rpcdaos *service.RPCDaos) *Service {
	s := &Service{
		c:   c,
		tpl: template.New(c),
	}
	return s
}

// Ping service
func (s *Service) Ping(c context.Context) (err error) {
	if err = s.tpl.Ping(c); err != nil {
		log.Error("s.template.Dao.PingDb err(%v)", err)
	}
	return
}

// Close dao
func (s *Service) Close() {
	s.tpl.Close()
}
