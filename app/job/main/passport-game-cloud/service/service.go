package service

import (
	"context"
	"sync"
	"time"

	"go-common/app/job/main/passport-game-cloud/conf"
	"go-common/app/job/main/passport-game-cloud/dao"
	"go-common/library/queue/databus"
)

const (
	// cache retry count and duration
	_accountCacheRetryCount    = 3
	_accountCacheRetryDuration = time.Second

	// table and duration
	_asoAccountTable   = "aso_account"
	_memberTablePrefix = "dede_member"

	_changePwd      = "changePwd"
	_updateUserInfo = "updateUserInfo"

	_gameAppID = int32(876)
)

// Service service.
type Service struct {
	c      *conf.Config
	d      *dao.Dao
	missch chan func()

	// prom
	tokenInterval      *Interval
	memberInterval     *Interval
	asoAccountInterval *Interval
	transInterval      *Interval

	gameAppIDs []int32
	// bin log proc
	binLogDataBus    *databus.Databus
	binLogMergeChans []chan *message
	binLogDoneChan   chan []*message
	head, last       *message
	mu               sync.Mutex
	// aso acc proc
	encryptTransDataBus    *databus.Databus
	encryptTransMergeChans []chan *message
	encryptTransDoneChan   chan []*message
	asoHead, asoLast       *message
	asoMu                  sync.Mutex
}

type message struct {
	next   *message
	data   *databus.Message
	object interface{}
	done   bool
}

// New new a service instance.
func New(c *conf.Config) (s *Service) {
	gameAppIDs := make([]int32, 0)
	gameAppIDs = append(gameAppIDs, _gameAppID)
	for _, id := range c.Game.AppIDs {
		if id == _gameAppID {
			continue
		}
		gameAppIDs = append(gameAppIDs, id)
	}

	s = &Service{
		c:      c,
		d:      dao.New(c),
		missch: make(chan func(), 10240),

		// prom
		tokenInterval: NewInterval(&IntervalConfig{
			Name: "interval_token",
			Rate: 1000,
		}),
		memberInterval: NewInterval(&IntervalConfig{
			Name: "interval_member",
			Rate: 1000,
		}),
		asoAccountInterval: NewInterval(&IntervalConfig{
			Name: "interval_aso_account",
			Rate: 1000,
		}),
		transInterval: NewInterval(&IntervalConfig{
			Name: "interval_trans",
			Rate: 1000,
		}),
		gameAppIDs:       gameAppIDs,
		binLogDataBus:    databus.New(c.DataBus.BinLogSub),
		binLogMergeChans: make([]chan *message, c.Group.BinLog.Num),
		binLogDoneChan:   make(chan []*message, c.Group.BinLog.Chan),
		// aso acc proc
		encryptTransDataBus:    databus.New(c.DataBus.EncryptTransSub),
		encryptTransMergeChans: make([]chan *message, c.Group.EncryptTrans.Num),
		encryptTransDoneChan:   make(chan []*message, c.Group.EncryptTrans.Chan),
	}
	// start bin log proc
	go s.binlogcommitproc()
	for i := 0; i < c.Group.BinLog.Num; i++ {
		ch := make(chan *message, c.Group.BinLog.Chan)
		s.binLogMergeChans[i] = ch
		go s.binlogmergeproc(ch)
	}
	go s.binlogconsumeproc()
	// start encrypt trans proc
	go s.encrypttranscommitproc()
	for i := 0; i < c.Group.EncryptTrans.Num; i++ {
		ch := make(chan *message, c.Group.EncryptTrans.Chan)
		s.encryptTransMergeChans[i] = ch
		go s.encrypttransmergeproc(ch)
	}
	go s.encrypttransconsumeproc()
	// go s.cacheproc()
	return
}

//func (s *Service) addCache(f func()) {
//	select {
//	case s.missch <- f:
//	default:
//		log.Warn("cache chan full")
//	}
//}

// cacheproc is a routine for executing closure.
//func (s *Service) cacheproc() {
//	for {
//		f := <-s.missch
//		f()
//	}
//}

// Ping check server ok.
func (s *Service) Ping(c context.Context) (err error) {
	err = s.d.Ping(c)
	return
}

// Close close service, including closing dao.
func (s *Service) Close() (err error) {
	s.binLogDataBus.Close()
	s.encryptTransDataBus.Close()
	s.d.Close()
	return
}
