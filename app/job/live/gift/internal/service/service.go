package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"go-common/app/job/live/gift/internal/model"
	"go-common/library/log"
	"go-common/library/queue/databus"

	pb "go-common/app/job/live/gift/api"
	"go-common/app/job/live/gift/internal/conf"
	"go-common/app/job/live/gift/internal/dao"

	"github.com/golang/protobuf/ptypes/empty"
)

// Service struct
type Service struct {
	c           *conf.Config
	dao         *dao.Dao
	addFreeGift *databus.Databus
	waiter      sync.WaitGroup
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:           c,
		dao:         dao.New(c),
		addFreeGift: databus.New(c.Databus.AddGift),
	}
	go s.infocproc()
	for i := 0; i < c.Consumer.AddGift.Num; i++ {
		s.waiter.Add(1)
		go s.addGiftConsumeProc()
	}
	return s
}

func (s *Service) addGiftConsumeProc() {
	defer s.waiter.Done()
	var err error
	for {
		msg, ok := <-s.addFreeGift.Messages()
		if !ok {
			log.Error("s.addFreeGift.Messages channel closed")
			return
		}
		m := &model.AddFreeGift{}
		if err = json.Unmarshal(msg.Value, &m); err != nil {
			log.Error("json.Unmarshal(%v) error(%v)", msg, err)
			continue
		}
		ctx := context.Background()
		// 消息幂等
		if m.MsgID != "" {
			key := m.MsgID + m.Source
			gotLock, _, errLock := s.dao.Lock(ctx, key, 3600000, 0, 0)
			if errLock != nil {
				continue
			}
			if !gotLock {
				log.Error("msg has been processed,%v", m)
				continue
			}
		}
		s.AddGift(ctx, m)
		// 打点上报
		s.giftActionInfoc(m.UID, 0, m.GiftID, 0, m.GiftNum, m.Source, "")
		log.Info("consume addFreeGift topic:%s, Key:%s, Value:%s ", msg.Topic, msg.Key, msg.Value)
		if err = msg.Commit(); err != nil {
			log.Error("commit msg(%v) error(%v)", msg, err)
		}
	}
}

// SayHello grpc demo func
func (s *Service) SayHello(ctx context.Context, req *pb.HelloReq) (reply *empty.Empty, err error) {
	reply = new(empty.Empty)
	fmt.Printf("hello %s", req.Name)
	return
}

// Ping Service
func (s *Service) Ping(ctx context.Context) (err error) {
	return s.dao.Ping(ctx)
}

// Close Service
func (s *Service) Close() {
	s.dao.Close()
}
