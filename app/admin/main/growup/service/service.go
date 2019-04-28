package service

import (
	"context"

	"go-common/app/admin/main/growup/conf"
	"go-common/app/admin/main/growup/dao"
	"go-common/app/admin/main/growup/dao/message"
	"go-common/app/admin/main/growup/dao/resource"
	"go-common/app/admin/main/growup/dao/shell"
	"go-common/app/admin/main/growup/model/offlineactivity"
	"go-common/library/net/http/blademaster"
)

// Service struct
type Service struct {
	conf                *conf.Config
	dao                 *dao.Dao
	msg                 *message.Dao
	chanCheckDb         chan int
	chanCheckShellOrder chan *offlineactivity.OfflineActivityResult
	chanCheckActivity   chan int64 // it's result id in this channel
	shellClient         *shell.Client
}

// New fn
func New(c *conf.Config) (s *Service) {
	s = &Service{
		conf:                c,
		dao:                 dao.New(c),
		msg:                 message.New(c),
		chanCheckDb:         make(chan int, 1),
		chanCheckShellOrder: make(chan *offlineactivity.OfflineActivityResult, 10240),
		chanCheckActivity:   make(chan int64, 1000),
		shellClient:         shell.New(c.ShellConf, blademaster.NewClient(c.HTTPClient)),
	}
	resource.Init(c)
	if c.OtherConf.OfflineOrderConsume {
		go s.offlineactivityCheckSendDbProc()
	}
	go s.offlineactivityCheckShellOrderProc()
	return s
}

// Ping fn
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close dao
func (s *Service) Close() {
	s.dao.Close()
}
