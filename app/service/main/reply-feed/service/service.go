package service

import (
	"context"
	"strconv"
	"sync"
	"time"

	"go-common/app/service/main/reply-feed/conf"
	"go-common/app/service/main/reply-feed/dao"
	"go-common/app/service/main/reply-feed/model"
	"go-common/library/log"
	"go-common/library/net/netutil"

	"github.com/robfig/cron"
)

// Service struct
type Service struct {
	c   *conf.Config
	dao *dao.Dao
	// backoff
	bc              netutil.BackoffConfig
	cron            *cron.Cron
	statisticsStats []model.StatisticsStat
	statisticsLock  sync.RWMutex
	// eventProducer   *databus.Databus
	midMapping   map[int64]int
	oidWhiteList map[int64]int
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:   c,
		dao: dao.New(c),
		bc: netutil.BackoffConfig{
			MaxDelay:  1 * time.Second,
			BaseDelay: 100 * time.Millisecond,
			Factor:    1.6,
			Jitter:    0.2,
		},
		cron:            cron.New(),
		statisticsStats: make([]model.StatisticsStat, model.SlotsNum),
		// eventProducer:   databus.New(c.Databus.Event),
		midMapping:   make(map[int64]int),
		oidWhiteList: make(map[int64]int),
	}
	// toml不支持int为key
	for k, v := range s.c.MidMapping {
		mid, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			continue
		}
		s.midMapping[mid] = v
	}
	for k, v := range s.c.OidWhiteList {
		oid, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			continue
		}
		s.oidWhiteList[oid] = v
	}
	// 初始化各个实验组所占流量槽位
	var err error
	if err = s.loadSlots(); err != nil {
		panic(err)
	}
	go s.loadproc()
	// 每整小时执行一次将统计数据写入DB
	s.cron.AddFunc("@hourly", func() {
		s.persistStatistics()
	})
	s.cron.Start()
	return s
}

// namePercent ...
// func (s *Service) namePercent() map[string]int64 {
// 	p := make(map[string]int64)
// 	s.statisticsLock.RLock()
// 	for _, slot := range s.statisticsStats {
// 		if _, ok := p[slot.Name]; ok {
// 			p[slot.Name]++
// 		} else {
// 			p[slot.Name] = 1
// 		}
// 	}
// 	s.statisticsLock.RUnlock()
// 	return p
// }

func (s *Service) loadproc() {
	for {
		time.Sleep(time.Minute)
		s.loadSlots()
	}
}

func (s *Service) loadSlots() (err error) {
	slotsMapping, err := s.dao.SlotsMapping(context.Background())
	if err != nil {
		return
	}
	s.statisticsLock.Lock()
	for _, mapping := range slotsMapping {
		for _, slot := range mapping.Slots {
			s.statisticsStats[slot].Name = mapping.Name
			s.statisticsStats[slot].Slot = slot
			s.statisticsStats[slot].State = mapping.State
		}
	}
	s.statisticsLock.Unlock()
	log.Warn("statistics stat (%v)", s.statisticsStats)
	return
}

func (s *Service) persistStatistics() {
	ctx := context.Background()
	statisticsMap := make(map[string]*model.StatisticsStat)
	s.statisticsLock.RLock()
	for _, stat := range s.statisticsStats {
		s, ok := statisticsMap[stat.Name]
		if ok {
			statisticsMap[stat.Name] = s.Merge(&stat)
		} else {
			statisticsMap[stat.Name] = &stat
		}
	}
	s.statisticsLock.RUnlock()
	now := time.Now()
	year, month, day := now.Date()
	date := year*10000 + int(month)*100 + day
	hour := now.Hour()
	for name, stat := range statisticsMap {
		err := s.dao.UpsertStatistics(ctx, name, date, hour, stat)
		var (
			retryTimes    = 0
			maxRetryTimes = 5
		)
		for err != nil && retryTimes < maxRetryTimes {
			time.Sleep(s.bc.Backoff(retryTimes))
			err = s.dao.UpsertStatistics(ctx, name, date, hour, stat)
			retryTimes++
		}
		if retryTimes >= maxRetryTimes {
			log.Error("upsert statistics error retry reached max retry times.")
		}
	}
	for i := range s.statisticsStats {
		reset(&s.statisticsStats[i])
	}
}

func reset(stat *model.StatisticsStat) {
	stat.HotView = 0
	stat.View = 0
	stat.HotClick = 0
	stat.TotalView = 0
}

// Ping Service
func (s *Service) Ping(ctx context.Context) (err error) {
	return s.dao.Ping(ctx)
}

// Close Service
func (s *Service) Close() {
	s.persistStatistics()
	s.dao.Close()
}
