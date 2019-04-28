package service

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"go-common/app/job/main/answer/conf"
	"go-common/app/job/main/answer/dao"
	"go-common/app/job/main/answer/model"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

const (
	_insertAction = "insert"
	_updateAction = "update"

	_labourTable = "blocked_labour_question"
)

// Service service def.
type Service struct {
	c              *conf.Config
	dao            *dao.Dao
	labourDatabus  *databus.Databus
	accountFormal  *databus.Databus
	waiter         sync.WaitGroup
	uploadInterval time.Duration
	closed         bool
}

// New create a instance of Service and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:              c,
		dao:            dao.New(c),
		uploadInterval: time.Duration(c.Properties.UploadInterval),
	}
	if c.Databus.Labour != nil {
		s.labourDatabus = databus.New(c.Databus.Labour)
		s.waiter.Add(1)
		go s.labourproc()
	}
	if c.Databus.Account != nil {
		s.accountFormal = databus.New(c.Databus.Account)
		go s.formalproc()
	}
	s.waiter.Add(1)
	go s.loadextarqueproc()
	return s
}

func (s *Service) labourproc() {
	defer s.waiter.Done()
	var (
		err     error
		msg     *databus.Message
		msgChan = s.labourDatabus.Messages()
		ok      bool
	)
	for {
		msg, ok = <-msgChan
		if !ok {
			log.Info("labour msgChan closed")
		}
		if s.closed {
			return
		}
		if err = msg.Commit(); err != nil {
			log.Error("msg.Commit err(%v)", err)
		}
		v := &model.MsgCanal{}
		if err = json.Unmarshal([]byte(msg.Value), v); err != nil {
			log.Error("json.Unmarshal(%v) err(%v)", v, err)
			continue
		}
		if v.Table == _labourTable {
			switch v.Action {
			case _insertAction:
				s.AddLabourQuestion(context.Background(), v)
			case _updateAction:
				s.ModifyLabourQuestion(context.Background(), v)
			}
		}
	}
}

func (s *Service) formalproc() {
	var (
		ok      bool
		err     error
		msg     *databus.Message
		msgChan = s.accountFormal.Messages()
	)
	for {
		msg, ok = <-msgChan
		if !ok {
			log.Info("account formal msgChan closed")
		}
		if s.closed {
			return
		}
		v := &model.Formal{}
		if err = json.Unmarshal([]byte(msg.Value), v); err != nil {
			log.Error("json.Unmarshal(%v) err(%v)", v, err)
			continue
		}
		for retries := 0; retries < s.c.Properties.MaxRetries; retries++ {
			if err = s.dao.BeFormal(context.Background(), v.Mid, v.IP); err != nil {
				sleep := s.c.Backoff.Backoff(retries)
				log.Error("s.dao.BeFormal(%+v) sleep(%d) err(%+v)", v, sleep, err)
				time.Sleep(sleep * time.Second)
				continue
			}
			break
		}
		if err = msg.Commit(); err != nil {
			log.Error("msg.Commit err(%v)", err)
		}
	}
}

// Close all resource.
func (s *Service) Close() (err error) {
	defer s.waiter.Wait()
	s.closed = true
	s.dao.Close()
	if err = s.labourDatabus.Close(); err != nil {
		log.Error("s.labourDatabus.Close() error(%v)", err)
		return
	}
	return
}

// Ping check dao health.
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

func (s *Service) loadextarqueproc() {
	defer s.waiter.Done()
	for {
		time.Sleep(s.uploadInterval)
		if s.closed {
			return
		}
		res, err := s.dao.QidsExtraByState(context.Background(), model.LimitSize)
		if err != nil {
			log.Error("s.dao.QidsExtraByState() error(%v)", err)
			continue
		}
		if len(res) == 0 {
			continue
		}
		for _, q := range res {
			s.UploadQueImg(context.Background(), q)
		}
	}
}
