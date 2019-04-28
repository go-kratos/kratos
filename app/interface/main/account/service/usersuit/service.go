package usersuit

import (
	"context"

	"go-common/app/interface/main/account/conf"
	"go-common/app/interface/main/account/dao/usersuit"
	"go-common/app/interface/main/account/dao/vip"
	accrpc "go-common/app/service/main/account/rpc/client"
	coinrpc "go-common/app/service/main/coin/api/gorpc"
	memrpc "go-common/app/service/main/member/api/gorpc"
	usrpc "go-common/app/service/main/usersuit/rpc/client"
)

// Service struct.
type Service struct {
	c       *conf.Config
	dao     *usersuit.Dao
	vipDao  *vip.Dao
	usRPC   *usrpc.Service2
	accRPC  *accrpc.Service3
	coinRPC *coinrpc.Service
	memRPC  *memrpc.Service
}

// New a pendant service
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:       c,
		dao:     usersuit.New(c),
		vipDao:  vip.New(c),
		usRPC:   usrpc.New(c.RPCClient2.Usersuit),
		memRPC:  memrpc.New(c.RPCClient2.Member),
		accRPC:  accrpc.New3(c.RPCClient2.Account),
		coinRPC: coinrpc.New(c.RPCClient2.Coin),
	}
	return
}

// Ping check server ok.
func (s *Service) Ping(c context.Context) (err error) {
	return
}

// Close dao.
func (s *Service) Close() {}
