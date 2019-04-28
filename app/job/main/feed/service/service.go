package service

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"go-common/app/job/main/feed/conf"
	"go-common/app/job/main/feed/dao"
	"go-common/app/job/main/feed/model"
	feed "go-common/app/service/main/feed/rpc/client"
	"go-common/library/conf/env"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

type Service struct {
	c          *conf.Config
	dao        *dao.Dao
	archiveSub *databus.Databus
	arcUpMo    int64
	feedRPC    *feed.Service
	waiter     *sync.WaitGroup
}

// New is feed service implementation.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:          c,
		dao:        dao.New(c),
		archiveSub: databus.New(c.ArchiveSub),
		feedRPC:    feed.New(c.FeedRPC),
		waiter:     new(sync.WaitGroup),
	}
	// arc databus consumer
	s.waiter.Add(1)
	go s.arcConsumeproc()
	go s.checkConsumeproc()
	return s
}

// arcConsumeproc consumer archive
func (s *Service) arcConsumeproc() {
	var (
		msgs = s.archiveSub.Messages()
		err  error
	)
	defer s.waiter.Done()
	for {
		msg, ok := <-msgs
		if !ok {
			log.Info("arc databus Consumer exit")
			return
		}
		s.arcUpMo++
		dao.PromInfo("消费稿件变更")
		msg.Commit()
		m := &model.Message{}
		if err = json.Unmarshal(msg.Value, m); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
			continue
		}
		if m.Table != "archive" {
			continue
		}
		s.archiveUpdate(m.Action, m.New, m.Old)
	}
}

// checkConsumeproc check consumer stat
func (s *Service) checkConsumeproc() {
	if env.DeployEnv != env.DeployEnvProd {
		return
	}
	var arcMo int64
	for {
		time.Sleep(1 * time.Minute)
		if s.arcUpMo-arcMo == 0 {
			msg := "feed-job arhieve did not consume within a minute"
			s.dao.SendSMS(msg)
			log.Warn(msg)
		}

		arcMo = s.arcUpMo
	}
}

// Close Databus consumer close.
func (s *Service) Close() error {
	return s.archiveSub.Close()
}

// Wait goroutinue to close
func (s *Service) Wait() {
	s.waiter.Wait()
}

// Ping check server ok
func (s *Service) Ping(c context.Context) error {
	return s.dao.Ping(c)
}
