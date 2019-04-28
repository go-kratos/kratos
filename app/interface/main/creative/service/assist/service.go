package assist

import (
	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/dao/account"
	"go-common/app/interface/main/creative/dao/assist"
	"go-common/app/interface/main/creative/dao/danmu"
	"go-common/app/interface/main/creative/dao/reply"
	"go-common/app/interface/main/creative/service"
)

// Service assist.
type Service struct {
	c      *conf.Config
	assist *assist.Dao
	reply  *reply.Dao
	dm     *danmu.Dao
	acc    *account.Dao
}

// New get assist service.
func New(c *conf.Config, rpcdaos *service.RPCDaos) *Service {
	s := &Service{
		c:      c,
		assist: assist.New(c),
		reply:  reply.New(c),
		dm:     danmu.New(c),
		acc:    rpcdaos.Acc,
	}
	return s
}
