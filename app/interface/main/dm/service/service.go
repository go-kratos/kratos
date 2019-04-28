package service

import (
	"context"
	"time"

	"go-common/app/interface/main/dm/conf"
	"go-common/app/interface/main/dm/dao"
	"go-common/app/interface/main/dm/model"
	dmCli "go-common/app/interface/main/dm2/rpc/client"
	accoutCli "go-common/app/service/main/account/api"
	arcCli "go-common/app/service/main/archive/api/gorpc"
	assCli "go-common/app/service/main/assist/rpc/client"
)

// Service define Service struct
type Service struct {
	c *conf.Config
	// dao
	dao *dao.Dao
	// rpc
	acvSvc     *arcCli.Service2
	accountSvc accoutCli.AccountClient
	astSvc     *assCli.Service
	dmRPC      *dmCli.Service
	//proc
	delDMReportChan  chan *model.Report
	hideDMReportChan chan *model.Report
}

// New new a Service and return.cdfg
func New(c *conf.Config) *Service {
	s := &Service{
		c: c,
		// dmDao
		dao: dao.New(c),
		// archive rpc service
		acvSvc: arcCli.New2(c.ArchiveRPC),
		astSvc: assCli.New(c.AssistRPC),
		dmRPC:  dmCli.New(c.DMRPC),
		//proc
		delDMReportChan:  make(chan *model.Report, 1024),
		hideDMReportChan: make(chan *model.Report, 1024),
	}
	accountSvc, err := accoutCli.NewClient(c.AccountRPC)
	if err != nil {
		panic(err)
	}
	s.accountSvc = accountSvc
	go s.dmReportProc()
	go s.cronproc()
	return s
}

// Ping check server ok
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// cronproc 一分钟执行一次
func (s *Service) cronproc() {
	for {
		<-time.After(time.Minute)
		go s.sendProtectNotifyToUp(context.Background())
		go s.sendProtectNotifyToUser(context.Background())
	}
}
