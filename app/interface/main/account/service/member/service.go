package member

import (
	"context"
	"net"

	"go-common/app/interface/main/account/conf"
	"go-common/app/interface/main/account/dao/account"
	"go-common/app/interface/main/account/dao/passport"
	"go-common/app/interface/main/account/dao/reply"
	artRPC "go-common/app/interface/openplatform/article/rpc/client"
	accrpc "go-common/app/service/main/account/rpc/client"
	arcRPC "go-common/app/service/main/archive/api/gorpc"
	coinrpc "go-common/app/service/main/coin/api/gorpc"
	filterrpc "go-common/app/service/main/filter/rpc/client"
	locrpc "go-common/app/service/main/location/rpc/client"
	memrpc "go-common/app/service/main/member/api/gorpc"
	passRPC "go-common/app/service/main/passport/rpc/client"
	securerpc "go-common/app/service/main/secure/rpc/client"
	upRPC "go-common/app/service/main/up/api/gorpc"
	usrpc "go-common/app/service/main/usersuit/rpc/client"
	"go-common/library/queue/databus"

	"github.com/pkg/errors"
)

// Service struct of service.
type Service struct {
	// conf
	c *conf.Config

	accRPC             *accrpc.Service3
	memRPC             *memrpc.Service
	usRPC              *usrpc.Service2
	arcRPC             *arcRPC.Service2
	upRPC              *upRPC.Service
	artRPC             *artRPC.Service
	passRPC            *passRPC.Client2
	coinRPC            *coinrpc.Service
	locRPC             *locrpc.Service
	secureRPC          *securerpc.Service
	filterRPC          *filterrpc.Service
	accDao             *account.Dao
	replyDao           *reply.Dao
	passDao            *passport.Dao
	nickFreeAppKeys    map[string]string
	accountNotify      *databus.Databus
	removeLoginLogCIDR []*net.IPNet
}

// New create service instance and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:               c,
		memRPC:          memrpc.New(c.RPCClient2.Member),
		accRPC:          accrpc.New3(c.RPCClient2.Account),
		usRPC:           usrpc.New(c.RPCClient2.Usersuit),
		arcRPC:          arcRPC.New2(c.RPCClient2.Archive),
		upRPC:           upRPC.New(c.RPCClient2.UP),
		artRPC:          artRPC.New(c.RPCClient2.Article),
		passRPC:         passRPC.New(c.RPCClient2.PassPort),
		coinRPC:         coinrpc.New(c.RPCClient2.Coin),
		locRPC:          locrpc.New(c.RPCClient2.Location),
		secureRPC:       securerpc.New(c.RPCClient2.Secure),
		filterRPC:       filterrpc.New(c.RPCClient2.Filter),
		accDao:          account.New(c),
		passDao:         passport.New(c),
		replyDao:        reply.New(c),
		nickFreeAppKeys: c.NickFreeAppKeys,
		accountNotify:   databus.New(c.AccountNotify),
	}
	cidrs := make([]*net.IPNet, 0, len(c.Account.RemoveLoginLogCIDR))
	for _, raw := range c.Account.RemoveLoginLogCIDR {
		_, inet, err := net.ParseCIDR(raw)
		if err != nil {
			panic(errors.Wrapf(err, "Invalid CIDR: %s", raw))
		}
		cidrs = append(cidrs, inet)
	}
	s.removeLoginLogCIDR = cidrs
	return
}

// Ping check server ok.
func (s *Service) Ping(c context.Context) (err error) {
	return
}

// Close dao.
func (s *Service) Close() {}
