package service

import (
	"context"
	"sync"
	"time"

	"go-common/app/job/main/reply-feed/conf"
	"go-common/app/job/main/reply-feed/dao"
	"go-common/app/job/main/reply-feed/model"
	"go-common/library/log"
	"go-common/library/net/netutil"
	"go-common/library/queue/databus"
	"go-common/library/sync/pipeline/fanout"

	"github.com/ivpusic/grpool"
	"github.com/robfig/cron"
)

// Service struct
type Service struct {
	c   *conf.Config
	dao *dao.Dao
	// 定时任务
	cron *cron.Cron
	// backoff
	bc            netutil.BackoffConfig
	statsConsumer *databus.Databus
	// eventConsumer  *databus.Databus
	taskQ      *fanout.Fanout
	uvQ        *fanout.Fanout
	statQ      *fanout.Fanout
	replyListQ *fanout.Fanout
	waiter     sync.WaitGroup

	// 专门计算热评分数的goroutine pool
	calculator *grpool.Pool

	statisticsStats [model.SlotsNum]model.StatisticsStat

	algorithmsLock sync.RWMutex
	statisticsLock sync.RWMutex
	algorithms     []model.Algorithm
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:    c,
		dao:  dao.New(c),
		cron: cron.New(),
		bc: netutil.BackoffConfig{
			MaxDelay:  1 * time.Second,
			BaseDelay: 100 * time.Millisecond,
			Factor:    1.6,
			Jitter:    0.2,
		},
		statsConsumer: databus.New(c.Databus.Stats),
		// eventConsumer: databus.New(c.Databus.Event),
		// 处理异步写任务的goroutine
		taskQ:      fanout.New("task"),
		uvQ:        fanout.New("uv-task", fanout.Worker(4), fanout.Buffer(2048)),
		statQ:      fanout.New("memcache", fanout.Worker(4), fanout.Buffer(2048)),
		replyListQ: fanout.New("redis", fanout.Worker(4), fanout.Buffer(2048)),
		calculator: grpool.NewPool(4, 2048),
	}
	var err error
	if err = s.loadAlgorithm(); err != nil {
		panic(err)
	}
	if err = s.loadSlots(); err != nil {
		panic(err)
	}
	go s.loadproc()

	// 消费databus
	s.waiter.Add(1)
	go s.statsproc()
	// s.waiter.Add(1)
	// go s.eventproc()

	// 每整小时执行一次将统计数据写入DB
	s.cron.AddFunc("@hourly", func() {
		s.persistStatistics()
	})
	s.cron.Start()
	return s
}

func (s *Service) loadproc() {
	for {
		time.Sleep(time.Minute)
		s.loadAlgorithm()
		s.loadSlots()
	}
}

// Ping Service
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close Service
func (s *Service) Close() {
	s.statsConsumer.Close()
	// s.eventConsumer.Close()
	log.Warn("consumer closed")
	s.waiter.Wait()
	s.persistStatistics()
	s.dao.Close()
}
