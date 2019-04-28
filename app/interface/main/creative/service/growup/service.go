package growup

import (
	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/dao/archive"
	"go-common/app/interface/main/creative/dao/growup"
	"go-common/app/interface/main/creative/service"
)

//Service struct.
type Service struct {
	c      *conf.Config
	arc    *archive.Dao
	growup *growup.Dao
	p      *service.Public
}

//New get service.
func New(c *conf.Config, rpcdaos *service.RPCDaos, p *service.Public) *Service {
	s := &Service{
		c:      c,
		arc:    rpcdaos.Arc,
		growup: growup.New(c),
		p:      p,
	}
	return s
}
