package danmu

import (
	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/dao/account"
	"go-common/app/interface/main/creative/dao/archive"
	"go-common/app/interface/main/creative/dao/danmu"
	"go-common/app/interface/main/creative/dao/elec"
	"go-common/app/interface/main/creative/dao/subtitle"
	"go-common/app/interface/main/creative/service"
)

// Service danmu.
type Service struct {
	c    *conf.Config
	dm   *danmu.Dao
	arc  *archive.Dao
	acc  *account.Dao
	sub  *subtitle.Dao
	elec *elec.Dao
}

// New get danmu service.
func New(c *conf.Config, rpcdaos *service.RPCDaos) *Service {
	s := &Service{
		c:    c,
		dm:   danmu.New(c),
		acc:  rpcdaos.Acc,
		arc:  rpcdaos.Arc,
		sub:  rpcdaos.Sub,
		elec: elec.New(c),
	}
	return s
}
