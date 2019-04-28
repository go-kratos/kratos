package Service

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"go-common/app/job/live/wallet/conf"
	"go-common/app/job/live/wallet/dao"
	"go-common/app/job/live/wallet/model"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

// Service struct
type Service struct {
	c        *conf.Config
	dao      *dao.Dao
	userSub  *databus.Databus
	waiter   *sync.WaitGroup
	userUpMo int64
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:       c,
		dao:     dao.New(c),
		userSub: databus.New(c.UserSub),
		waiter:  new(sync.WaitGroup),
	}
	s.waiter.Add(1)
	go s.userCanalConsumeproc()
	go s.checkUserCanalConsumeproc()
	return s
}

// Ping Service
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close Service
func (s *Service) Close() {
	defer s.waiter.Wait()
	s.userSub.Close()
	s.dao.Close()
}

// Wait goroutinue to close
func (s *Service) Wait() {
	s.waiter.Wait()
}

// expCanalConsumeproc consumer archive
func (s *Service) userCanalConsumeproc() {
	var (
		msgs = s.userSub.Messages()
		err  error
	)
	defer s.waiter.Done()
	for {
		msg, ok := <-msgs
		if !ok {
			log.Info("userCanal databus Consumer exit")
			return
		}
		s.userUpMo++
		msg.Commit()
		m := &model.Message{}
		//log.Info("canal message %s", msg.Value)
		if err = json.Unmarshal(msg.Value, m); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
			continue
		}
		if !strings.HasPrefix(m.Table, "user_") || (m.Action != "update" && m.Action != "insert") {
			continue
		}
		s.mergeData(m.New, m.Old, m.Action)
	}
}

// checkConsumeproc check consumer stat
func (s *Service) checkUserCanalConsumeproc() {
	if s.c.Env != "pro" {
		return
	}
	var userMo int64
	for {
		time.Sleep(1 * time.Minute)
		if s.userUpMo-userMo == 0 {
			msg := "live-wallet-job userCanal did not consume within a minute"
			//s.dao.SendSMS(msg)
			log.Warn(msg)
		}

		userMo = s.userUpMo
	}
}
