package service

import (
	"context"
	"fmt"
	"sync"

	"go-common/app/job/main/relation-cache/conf"
	"go-common/app/job/main/relation-cache/dao"
	"go-common/library/queue/databus"
)

// Service struct
type Service struct {
	c              *conf.Config
	dao            *dao.Dao
	relationBinLog *databus.Databus
	waiter         *sync.WaitGroup
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:              c,
		dao:            dao.New(c),
		relationBinLog: databus.New(c.RelationBinLog),
		waiter:         &sync.WaitGroup{},
	}
	s.Start(context.Background())
	return s
}

// Start to handle requests
func (s *Service) Start(ctx context.Context) {
	for i := 0; i < 50; i++ {
		go s.relationBinLogproc(context.Background())
	}
}

// Ping Service
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close Service
func (s *Service) Close() {
	s.dao.Close()
}

// BeautifyMessage is
func BeautifyMessage(msg *databus.Message) string {
	pmsg := struct {
		Key       string `json:"key"`
		Value     string `json:"value"`
		Topic     string `json:"topic"`
		Partition int32  `json:"partition"`
		Offset    int64  `json:"offset"`
		Timestamp int64  `json:"timestamp"`
	}{
		Key:       msg.Key,
		Value:     string(msg.Value),
		Topic:     msg.Topic,
		Partition: msg.Partition,
		Offset:    msg.Offset,
		Timestamp: msg.Timestamp,
	}
	return fmt.Sprintf("%+v", pmsg)
}
