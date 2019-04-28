package block

import (
	"context"
	"runtime/debug"

	"go-common/app/admin/main/member/conf"
	"go-common/app/admin/main/member/dao/block"
	account "go-common/app/service/main/account/api"
	rpcfigure "go-common/app/service/main/figure/rpc/client"
	rpcspy "go-common/app/service/main/spy/rpc/client"
	"go-common/library/log"
	"go-common/library/queue/databus"
	"go-common/library/sync/pipeline/fanout"
)

// Service struct
type Service struct {
	conf             *conf.Config
	dao              *block.Dao
	cache            *fanout.Fanout
	spyRPC           *rpcspy.Service
	figureRPC        *rpcfigure.Service
	accountClient    account.AccountClient
	missch           chan func()
	accountNotifyPub *databus.Databus
}

// New init
func New(conf *conf.Config, dao *block.Dao, spyRPC *rpcspy.Service, figureRPC *rpcfigure.Service,
	accountClient account.AccountClient, accountNotifyPub *databus.Databus) (s *Service) {
	s = &Service{
		conf:             conf,
		dao:              dao,
		cache:            fanout.New("memberAdminCache", fanout.Worker(1), fanout.Buffer(10240)),
		missch:           make(chan func(), 10240),
		accountNotifyPub: accountNotifyPub,
		spyRPC:           spyRPC,
		figureRPC:        figureRPC,
		accountClient:    accountClient,
	}
	go s.missproc()
	return s
}

func (s *Service) missproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.missproc panic(%+v) : %s", x, debug.Stack())
			go s.missproc()
		}
	}()
	for {
		f := <-s.missch
		f()
	}
}

func (s *Service) mission(f func()) {
	select {
	case s.missch <- f:
	default:
		log.Error("s.missch full")
	}
}

// Ping Service
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close Service
func (s *Service) Close() {
	s.dao.Close()
}
