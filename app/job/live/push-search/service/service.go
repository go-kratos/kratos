package service

import (
	"context"

	"go-common/app/job/live/push-search/conf"
	"go-common/app/job/live/push-search/dao"
	accountApi "go-common/app/service/main/account/api"

	"go-common/library/queue/databus"
	"sync"
)

const (
	_tableArchive = "ap_room"
)

// Service struct
type Service struct {
	c                  *conf.Config
	dao                *dao.Dao
	binLogMergeChan    []chan *message
	attentionMergeChan []chan *message
	unameMergeChan     []chan *message
	waiter             *sync.WaitGroup
	waiterChan         *sync.WaitGroup
	AccountClient      accountApi.AccountClient
}

type message struct {
	next   *message
	data   *databus.Message
	object interface{}
	done   bool
}

// New init
func New(c *conf.Config) (s *Service) {
	dao.InitAPI()
	s = &Service{
		c:                  c,
		dao:                dao.New(c),
		binLogMergeChan:    make([]chan *message, c.Group.RoomInfo.Num),
		attentionMergeChan: make([]chan *message, c.Group.Attention.Num),
		unameMergeChan:     make([]chan *message, c.Group.UserInfo.Num),
		waiter:             new(sync.WaitGroup),
		waiterChan:         new(sync.WaitGroup),
	}
	accountClient, err := accountApi.NewClient(nil)
	if err != nil {
		panic(err)
	}
	s.AccountClient = accountClient

	//ap room 表 binlog qps 高, hash roomId 并行
	for i := 0; i < c.Group.RoomInfo.Num; i++ {
		ch := make(chan *message, c.Group.RoomInfo.Chan)
		s.binLogMergeChan[i] = ch
		s.waiterChan.Add(1)
		go s.roomInfoNotifyHandleProc(ch)
	}

	for i := 0; i < c.Group.Attention.Num; i++ {
		ch := make(chan *message, c.Group.Attention.Chan)
		s.attentionMergeChan[i] = ch
		s.waiterChan.Add(1)
		go s.attentionNotifyHandleProc(ch)
	}

	for i := 0; i < c.Group.UserInfo.Num; i++ {
		ch := make(chan *message, c.Group.UserInfo.Chan)
		s.unameMergeChan[i] = ch
		s.waiterChan.Add(1)
		go s.unameNotifyHandleProc(ch)
	}
	s.waiter.Add(1)
	go s.roomInfoNotifyConsumeProc()
	s.waiter.Add(1)
	go s.attentionNotifyConsumeProc()
	s.waiter.Add(1)
	go s.unameNotifyConsumeProc()
	return s
}

// Ping Service
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close Service
func (s *Service) Close() {
	//databus chan close
	s.dao.Close()
	s.waiter.Wait()
	//task goroutine close
	for _, ch := range s.binLogMergeChan {
		close(ch)
	}

	for _, ch := range s.attentionMergeChan {
		close(ch)
	}

	for _, ch := range s.unameMergeChan {
		close(ch)
	}
	s.waiterChan.Wait()
	s.dao.PushSearchDataBus.Close()
}
