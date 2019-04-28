package service

import (
	"context"
	"sync"

	"go-common/app/job/main/usersuit/conf"
	medalDao "go-common/app/job/main/usersuit/dao/medal"
	pendantDao "go-common/app/job/main/usersuit/dao/pendant"
	"go-common/library/log"
	"go-common/library/queue/databus"

	"github.com/robfig/cron"
)

// Service struct of service.
type Service struct {
	pendantDao *pendantDao.Dao
	medalDao   *medalDao.Dao
	// conf
	c                *conf.Config
	accountNotifyPub *databus.Databus
	vipBinLogSub     *databus.Databus
	notifych         chan func()
	// wait group
	wg sync.WaitGroup
}

// New create service instance and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:                c,
		pendantDao:       pendantDao.New(c),
		medalDao:         medalDao.New(c),
		accountNotifyPub: databus.New(c.Databus.AccountNotify),
		vipBinLogSub:     databus.New(c.Databus.VipBinLog),
		notifych:         make(chan func(), 10240),
	}
	// this is function
	go s.startexpireproc()
	s.wg.Add(1)
	go s.notifyproc()
	s.wg.Add(1)
	go s.vipconsumerproc()
	t := cron.New()
	if len(s.c.Properties.MedalCron) != 0 {
		t.AddFunc(s.c.Properties.MedalCron, s.cronUpNameplate)
	}
	t.Start()
	return
}

func (s *Service) addNotify(f func()) {
	select {
	case s.notifych <- f:
	default:
		log.Warn("addNotify chan full")
	}
}

// notifyproc nofity clear cache
func (s *Service) notifyproc() {
	defer s.wg.Done()
	for {
		f := <-s.notifych
		f()
	}
}

// Close dao.
func (s *Service) Close() {
	if s.pendantDao != nil {
		s.pendantDao.Close()
		s.medalDao.Close()
	}
	s.wg.Wait()
}

// Ping check server ok.
func (s *Service) Ping(c context.Context) (err error) {
	if err = s.pendantDao.Ping(c); err != nil {
		return
	}
	return
}

// PDTStatStep PDT Stat step
type PDTStatStep struct {
	Start, End, Step int64
}

// PDTGHisStep PDT G his
type PDTGHisStep struct {
	Start, End, Step int64
}

// PDTOHisStep Order his
type PDTOHisStep struct {
	Start, End, Step int64
}
