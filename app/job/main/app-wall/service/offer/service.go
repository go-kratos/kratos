package offer

import (
	"context"
	"strings"
	"sync"
	"time"

	"go-common/app/job/main/app-wall/conf"
	offerDao "go-common/app/job/main/app-wall/dao/offer"
	"go-common/app/job/main/app-wall/model/offer"
	"go-common/library/log"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
)

// Service struct
type Service struct {
	c          *conf.Config
	dao        *offerDao.Dao
	consumer   *cluster.Consumer
	activeChan chan *offer.ActiveMsg
	closed     bool
	waiter     sync.WaitGroup
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:          c,
		dao:        offerDao.New(c),
		activeChan: make(chan *offer.ActiveMsg, 10240),
		closed:     false,
	}
	var err error
	if s.consumer, err = s.NewConsumer(); err != nil {
		log.Error("%+v", err)
		panic(err)
	}
	s.waiter.Add(1)
	go s.activeConsumer()
	s.waiter.Add(1)
	go s.activeproc()
	// retry consumer
	for i := 0; i < 4; i++ {
		s.waiter.Add(1)
		go s.retryproc()
	}
	return s
}

// Ping Service
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close Service
func (s *Service) Close() {
	s.closed = true
	s.consumer.Close()
	s.waiter.Wait()
	log.Info("app-wall-job closed.")
}

// NewConsumer new cluster consumer.
func (s *Service) NewConsumer() (*cluster.Consumer, error) {
	// cluster config
	cfg := cluster.NewConfig()
	// NOTE cluster auto commit offset interval
	cfg.Consumer.Offsets.CommitInterval = time.Second * 1
	// NOTE set fetch.wait.max.ms
	cfg.Consumer.MaxWaitTime = time.Millisecond * 100
	// NOTE errors that occur during offset management,if enabled, c.Errors channel must be read
	cfg.Consumer.Return.Errors = true
	// NOTE notifications that occur during consumer, if enabled, c.Notifications channel must be read
	cfg.Group.Return.Notifications = true
	// The initial offset to use if no offset was previously committed.
	// default: OffsetOldest
	if strings.ToLower(s.c.Consumer.Offset) != "new" {
		cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	} else {
		cfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	}
	// new cluster consumer
	log.Info("s.c.Consumer.Brokers: %v", s.c.Consumer.Brokers)
	return cluster.NewConsumer(s.c.Consumer.Brokers, s.c.Consumer.Group, []string{s.c.Consumer.Topic}, cfg)
}
