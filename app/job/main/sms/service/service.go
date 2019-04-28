package service

import (
	"context"
	"sync"

	"go-common/app/job/main/sms/conf"
	"go-common/app/job/main/sms/dao"
	"go-common/app/job/main/sms/dao/chuanglan"
	"go-common/app/job/main/sms/dao/mengwang"
	"go-common/app/job/main/sms/model"
	smsgrpc "go-common/app/service/main/sms/api"
	smsmdl "go-common/app/service/main/sms/model"
	"go-common/library/log"
	"go-common/library/queue/databus"
	"go-common/library/sync/pipeline/fanout"
)

// Service struct of service.
type Service struct {
	c         *conf.Config
	dao       *dao.Dao
	databus   *databus.Databus
	smsgrpc   smsgrpc.SmsClient
	waiter    *sync.WaitGroup
	sms       chan *smsmdl.ModelSend // 验证码
	actSms    chan *smsmdl.ModelSend // 营销
	batchSms  chan *smsmdl.ModelSend // 批量
	smsp      *model.ConcurrentRing
	intep     *model.ConcurrentRing
	actp      *model.ConcurrentRing
	batchp    *model.ConcurrentRing
	cache     *fanout.Fanout
	providers int
	closed    bool
	smsCount  int64
	blacklist map[string]struct{}
}

// New create service instance and return.
func New(c *conf.Config) (s *Service) {
	proLen := len(c.Provider.Providers)
	if c.Speedup.Switch {
		proLen *= 2 // 每个短信服务商再提供一个云线路 http client
	}
	s = &Service{
		c:         c,
		dao:       dao.New(c),
		databus:   databus.New(c.Databus),
		waiter:    new(sync.WaitGroup),
		sms:       make(chan *smsmdl.ModelSend, 10240),
		actSms:    make(chan *smsmdl.ModelSend, 10240),
		batchSms:  make(chan *smsmdl.ModelSend, 10240),
		smsp:      model.NewConcurrentRing(proLen),
		intep:     model.NewConcurrentRing(proLen),
		actp:      model.NewConcurrentRing(proLen),
		batchp:    model.NewConcurrentRing(proLen),
		cache:     fanout.New("async-task", fanout.Worker(1), fanout.Buffer(10240)),
		providers: proLen,
	}
	s.initBlacklist()
	var err error
	if s.smsgrpc, err = smsgrpc.NewClient(c.SmsGRPC); err != nil {
		panic(err)
	}
	s.initProviders()
	s.waiter.Add(1)
	go s.subproc()
	go s.monitorproc()
	for i := 0; i < s.c.Sms.SingleSendProc; i++ {
		s.waiter.Add(1)
		go s.smsproc()
		s.waiter.Add(1)
		go s.actsmsproc()
	}
	for i := 0; i < s.c.Sms.BatchSendProc; i++ {
		s.waiter.Add(1)
		go s.actbatchproc()
	}
	return
}

func (s *Service) initBlacklist() {
	s.blacklist = make(map[string]struct{})
	for _, v := range s.c.Sms.Blacklist {
		s.blacklist[v] = struct{}{}
	}
}

func (s *Service) initProviders() {
	// 创建本地网络 http client
	s.newProviders(s.c)
	if !s.c.Speedup.Switch {
		return
	}
	// 替换成 云加速线路 URL 配置
	s.c.Provider.MengWangSmsURL = s.c.Speedup.MengWangSmsURL
	s.c.Provider.MengWangActURL = s.c.Speedup.MengWangActURL
	s.c.Provider.MengWangBatchURL = s.c.Speedup.MengWangBatchURL
	s.c.Provider.MengWangInternationURL = s.c.Speedup.MengWangInternationURL
	s.c.Provider.ChuangLanSmsURL = s.c.Speedup.ChuangLanSmsURL
	s.c.Provider.ChuangLanActURL = s.c.Speedup.ChuangLanActURL
	s.c.Provider.ChuangLanInternationURL = s.c.Speedup.ChuangLanInternationURL
	s.c.Provider.MengWangSmsCallbackURL = s.c.Speedup.MengWangSmsCallbackURL
	s.c.Provider.MengWangActCallbackURL = s.c.Speedup.MengWangActCallbackURL
	s.c.Provider.MengWangInternationalCallbackURL = s.c.Speedup.MengWangInternationalCallbackURL
	s.c.Provider.ChuangLanSmsCallbackURL = s.c.Speedup.ChuangLanSmsCallbackURL
	s.c.Provider.ChuangLanActCallbackURL = s.c.Speedup.ChuangLanActCallbackURL
	s.c.Provider.ChuangLanInternationalCallbackURL = s.c.Speedup.ChuangLanInternationalCallbackURL
	// 创建云加速线路 http client
	s.newProviders(s.c)
}

func (s *Service) newProviders(c *conf.Config) {
	var cli model.Provider
	for _, p := range c.Provider.Providers {
		switch p {
		case smsmdl.ProviderMengWang:
			cli = mengwang.NewClient(c)
		case smsmdl.ProviderChuangLan:
			cli = chuanglan.NewClient(c)
		default:
			log.Error("invalid provider(%d)", p)
			continue
		}
		s.smsp.Value = cli
		s.smsp.Ring = s.smsp.Next()
		s.intep.Value = cli
		s.intep.Ring = s.intep.Next()
		s.actp.Value = cli
		s.actp.Ring = s.actp.Next()
		s.batchp.Value = cli
		s.batchp.Ring = s.batchp.Next()
		for i := 0; i < s.c.Sms.CallbackProc; i++ {
			s.dispatchCallback(p)
		}
	}
}

// Ping check service health.
func (s *Service) Ping(ctx context.Context) error {
	return s.dao.Ping(ctx)
}

// Close kafka consumer close.
func (s *Service) Close() {
	s.closed = true
	s.databus.Close()
	s.waiter.Wait()
}
