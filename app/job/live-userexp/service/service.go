package service

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"go-common/app/job/live-userexp/conf"
	"go-common/app/job/live-userexp/dao"
	"go-common/app/job/live-userexp/model"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

// Service http service
type Service struct {
	c       *conf.Config
	keys    map[string]string
	dao     *dao.Dao
	missch  chan func()
	expSub  *databus.Databus
	waiter  *sync.WaitGroup
	expUpMo int64
}

// New for new service obj
func New(c *conf.Config) *Service {
	s := &Service{
		c:      c,
		keys:   map[string]string{},
		dao:    dao.New(c),
		missch: make(chan func(), 1024),
		expSub: databus.New(c.ExpSub),
		waiter: new(sync.WaitGroup),
	}
	s.waiter.Add(1)
	go s.expCanalConsumeproc()
	go s.checkExpCanalConsumeproc()
	return s
}

// Ping check server ok
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close dao
func (s *Service) Close() {
	defer s.waiter.Wait()
	s.expSub.Close()
	s.dao.Close()
}

// expCanalConsumeproc consumer archive
func (s *Service) expCanalConsumeproc() {
	var (
		msgs = s.expSub.Messages()
		err  error
	)
	defer s.waiter.Done()
	for {
		msg, ok := <-msgs
		if !ok {
			log.Info("expCanal databus Consumer exit")
			return
		}
		s.expUpMo++
		msg.Commit()
		m := &model.Message{}
		if err = json.Unmarshal(msg.Value, m); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
			continue
		}
		if !strings.HasPrefix(m.Table, "user_exp_") || m.Action != "update" {
			continue
		}
		s.levelCacheUpdate(m.New, m.Old)
	}
}

// checkConsumeproc check consumer stat
func (s *Service) checkExpCanalConsumeproc() {
	if s.c.Env != "pro" {
		return
	}
	var expMo int64
	for {
		time.Sleep(1 * time.Minute)
		if s.expUpMo-expMo == 0 {
			msg := "live-userexp-job expCanal did not consume within a minute"
			//s.dao.SendSMS(msg)
			log.Warn(msg)
		}

		expMo = s.expUpMo
	}
}

// Wait goroutinue to close
func (s *Service) Wait() {
	s.waiter.Wait()
}
