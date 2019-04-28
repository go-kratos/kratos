package service

import (
	"context"
	"runtime"
	"sync"
	"time"

	"go-common/app/job/main/app/conf"
	monitordao "go-common/app/job/main/app/dao/monitor"
	spacedao "go-common/app/job/main/app/dao/space"
	viewdao "go-common/app/job/main/app/dao/view"
	"go-common/app/job/main/app/model"
	accapi "go-common/app/service/main/account/api"
	arcrpc "go-common/app/service/main/archive/api/gorpc"
	resmdl "go-common/app/service/main/resource/model"
	resrpc "go-common/app/service/main/resource/rpc/client"
	"go-common/library/queue/databus"
)

// Service is service.
type Service struct {
	c *conf.Config
	// vdao
	vdao       *viewdao.Dao
	arcRPC     *arcrpc.Service2
	accAPI     accapi.AccountClient
	spdao      *spacedao.Dao
	monitorDao *monitordao.Dao
	// sub
	archiveNotifySub *databus.Databus
	accountNotifySub *databus.Databus
	contributeSub    *databus.Databus
	waiter           sync.WaitGroup
	// space
	contributeChan chan *model.ContributeMsg
	closed         bool
	// stat sub
	statSub  map[string]*databus.Databus
	statChan []chan *model.StatMsg
	// archive
	aidsChan     chan []int64
	notifyMidMap *SyncMap
	resourceRPC  *resrpc.Service
	sideBars     map[int8]map[int][]*resmdl.SideBar
	sliceCnt     int64
}

type SyncMap struct {
	sync.Mutex
	Map map[int64]string
}

// New new a service.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:                c,
		vdao:             viewdao.New(c),
		spdao:            spacedao.New(c),
		monitorDao:       monitordao.New(c),
		arcRPC:           arcrpc.New2(c.ArchiveRPC),
		archiveNotifySub: databus.New(c.ArchiveNotifySub),
		accountNotifySub: databus.New(c.AccountNotifySub),
		resourceRPC:      resrpc.New(nil),
		aidsChan:         make(chan []int64, 10240),
		notifyMidMap:     &SyncMap{Map: map[int64]string{}},
		closed:           false,
		sliceCnt:         10,
	}
	var err error
	if s.accAPI, err = accapi.NewClient(nil); err != nil {
		panic("account GRPC not found!!!!!!!!!!!!!!!!!!!")
	}
	// stat sub
	s.statSub = make(map[string]*databus.Databus)
	s.statSub[model.TypeForView] = databus.New(c.StatViewSub)
	s.statSub[model.TypeForDm] = databus.New(c.StatDMSub)
	s.statSub[model.TypeForReply] = databus.New(c.StatReplySub)
	s.statSub[model.TypeForFav] = databus.New(c.StatFavSub)
	s.statSub[model.TypeForCoin] = databus.New(c.StatCoinSub)
	s.statSub[model.TypeForShare] = databus.New(c.StatShareSub)
	s.statSub[model.TypeForLike] = databus.New(c.StatLikeSub)
	s.statSub[model.TypeForRank] = databus.New(c.StatRankSub)
	// arc consumer
	s.waiter.Add(1)
	go s.arcConsumeproc()
	// stat consumer
	for bus := range s.statSub {
		s.waiter.Add(1)
		go s.statConsumeproc(bus)
	}
	for i := int64(0); i < s.sliceCnt; i++ {
		s.waiter.Add(1)
		s.statChan = append(s.statChan, make(chan *model.StatMsg, 1024))
		go s.statproc(i)
	}
	// contribute consumer
	if model.EnvRun() {
		s.contributeChan = make(chan *model.ContributeMsg, 10240)
		s.contributeSub = databus.New(c.ContributeSub)
		s.waiter.Add(1)
		go s.contributeConsumeproc()
		for i := 0; i < runtime.NumCPU(); i++ {
			s.waiter.Add(1)
			go s.contributeroc()
		}
	}
	// retry consumer
	for i := 0; i < 4; i++ {
		s.waiter.Add(1)
		go s.retryproc()
	}
	// account consumer
	s.waiter.Add(1)
	go s.accConsumeproc()
	s.waiter.Add(1)
	go s.notifyConsumeproc()
	// flush consumer
	if s.c.View.Flush {
		s.waiter.Add(1)
		go s.flushConsumeproc()
		for i := 0; i < 4; i++ {
			s.waiter.Add(1)
			go s.flushproc()
		}
	}
	go s.monitorproc()
	return
}

// Close Databus consumer close.
func (s *Service) Close() {
	s.closed = true
	s.archiveNotifySub.Close()
	s.accountNotifySub.Close()
	for _, sub := range s.statSub {
		sub.Close()
	}
	time.Sleep(10 * time.Second)
	for i := int64(0); i < s.sliceCnt; i++ {
		close(s.statChan[i])
	}
	s.waiter.Wait()
}

func (s *Service) Ping(c context.Context) (err error) {
	if err = s.vdao.PingMc(c); err != nil {
		return
	}
	err = s.vdao.PingRedis(c)
	return
}
