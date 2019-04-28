package service

import (
	"context"
	"fmt"

	"go-common/app/job/main/member-cache/conf"
	"go-common/app/job/main/member-cache/dao"
	"go-common/library/queue/databus"
)

// Service struct
type Service struct {
	c             *conf.Config
	dao           *dao.Dao
	memberBinLog  *databus.Databus
	blockBinLog   *databus.Databus
	accountNotify *databus.Databus
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:             c,
		dao:           dao.New(c),
		memberBinLog:  databus.New(c.MemberBinLog),
		blockBinLog:   databus.New(c.BlockBinLog),
		accountNotify: databus.New(c.AccountNotify),
	}
	s.Start(context.Background())
	return s
}

// Start to handle requests
func (s *Service) Start(ctx context.Context) {
	for i := 0; i < 10; i++ {
		go s.memberBinLogproc(context.Background())
	}
	for i := 0; i < 10; i++ {
		go s.blockBinLogproc(context.Background())
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
