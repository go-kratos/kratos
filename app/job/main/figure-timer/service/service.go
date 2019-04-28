package service

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"go-common/app/job/main/figure-timer/conf"
	"go-common/app/job/main/figure-timer/dao"
	"go-common/app/job/main/figure-timer/model"
	"go-common/library/log"

	"github.com/robfig/cron"
)

// Service struct of service.
type Service struct {
	c      *conf.Config
	dao    dao.Int
	missch chan func()
	curVer int64
	cron   *cron.Cron
}

// New create service instance and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:      c,
		dao:    dao.New(c),
		missch: make(chan func(), 1024),
	}
	s.cron = cron.New()
	s.cron.AddFunc(s.c.Property.CycleCron, s.cycleproc)
	if c.Property.CycleAll {
		s.cron.AddFunc(s.c.Property.CycleAllCron, s.cycleallproc)
	}
	go s.missproc()
	if c.Property.FixRecord {
		go s.fixproc()
	}
	s.cron.Start()
	return
}

func (s *Service) missproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("s.missproc panic(%v)", x)
			go s.missproc()
			log.Info("s.missproc recover")
		}
	}()
	for {
		for fn := range s.missch {
			fn()
		}
	}
}

func (s *Service) cycleproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("s.cycleproc panic(%v)", x)
			go s.cycleproc()
			log.Info("s.cycleproc recover")
		}
	}()
	var (
		err    error
		mids   []int64
		c      = context.TODO()
		wg     sync.WaitGroup
		newVer = weekVersion(time.Now().AddDate(0, 0, int(-s.c.Property.CalcWeekOffset*7)))
	)
	// Refresh Version
	atomic.StoreInt64(&s.curVer, newVer)
	log.Info("Calc active users start ver [%d]", s.curVer)
	rank.Init()
	// Calc figure concurrently
	for i := s.c.Property.PendingMidStart; i < s.c.Property.PendingMidShard; i++ {
		if mids, err = s.PendingMids(c, s.curVer, i, s.c.Property.PendingMidRetry); err != nil {
			log.Error("%+v", err)
		}
		if len(mids) == 0 {
			continue
		}
		smids := splitMids(mids, s.c.Property.ConcurrencySize)
		for c := range smids {
			csmids := smids[c]
			wg.Add(1)
			go func() {
				defer func() {
					wg.Done()
				}()
				for _, mid := range csmids {
					log.Info("Start handle mid [%d] figure ver [%d]", mid, s.curVer)
					if err = s.HandleFigure(context.TODO(), mid, s.curVer); err != nil {
						log.Error("%+v", err)
					}
				}
			}()
		}
		wg.Wait()
	}
	log.Info("Calc rank info start [%d]", s.curVer)
	s.calcRank(c, s.curVer)
	log.Info("Calc rank info finished [%d]", s.curVer)
	log.Info("Calc active users finished ver [%d]", s.curVer)
}

func splitMids(mids []int64, concurrencySize int64) (smids [][]int64) {
	if len(mids) == 0 {
		return
	}
	if concurrencySize == 0 {
		concurrencySize = 1
	}
	step := int64(len(mids))/concurrencySize + 1
	for c := int64(0); c < concurrencySize; c++ {
		var cMids []int64
		indexFrom := c * step
		indexTo := (c + 1) * step
		if indexFrom >= int64(len(mids)) {
			break
		}
		if indexTo >= int64(len(mids)) {
			cMids = mids[indexFrom:]
		} else {
			cMids = mids[indexFrom:indexTo]
		}
		smids = append(smids, cMids)
	}
	return
}

// PendingMids get pending mid list with retry
func (s *Service) PendingMids(c context.Context, version int64, shard int64, retry int64) (mids []int64, err error) {
	var (
		maxDo   = retry + 1
		doTimes int64
	)
	for doTimes < maxDo {
		if mids, err = s.dao.PendingMidsCache(c, s.curVer, shard); err != nil {
			doTimes++
			log.Info("s.dao.PendingMidsCache(%d,%d) retry (%d) error (%+v)", version, shard, doTimes, err)
		} else {
			doTimes = maxDo
		}
	}
	return
}

func (s *Service) cycleallproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("cycleallproc panic(%+v)", x)
		}
	}()
	var (
		ctx    = context.TODO()
		err    error
		newVer = weekVersion(time.Now().AddDate(0, 0, int(-s.c.Property.CalcWeekOffset*7)))
	)
	// Refresh Version
	atomic.StoreInt64(&s.curVer, newVer)
	log.Info("cycleallproc active users start ver [%d]", s.curVer)
	rank.Init()

	for shard := s.c.Property.PendingMidStart; shard < 100; shard++ {
		log.Info("cycleallproc start run: %d", shard)
		var (
			figures []*model.Figure
			fromMid = int64(shard)
			end     bool
		)
		for !end {
			if figures, end, err = s.dao.Figures(ctx, fromMid, 100); err != nil {
				log.Error("%+v", err)
				break
			}
			if len(figures) == 0 {
				continue
			}
			for _, figure := range figures {
				if fromMid < figure.Mid {
					fromMid = figure.Mid
				}
				log.Info("Start handle mid [%d] figure ver [%d]", figure.Mid, s.curVer)
				if err = s.HandleFigure(ctx, figure.Mid, s.curVer); err != nil {
					log.Error("%+v", err)
					continue
				}
			}
		}
		log.Info("cycleallproc rank info start [%d]", s.curVer)
		s.calcRank(ctx, s.curVer)
		log.Info("cycleallproc rank info finished [%d]", s.curVer)
		log.Info("cycleallproc active users finished ver [%d]", s.curVer)
	}
}

// Close kafka consumer close.
func (s *Service) Close() (err error) {
	s.dao.Close()
	return
}

// Wait wait service end.
func (s *Service) Wait() {
}

// Ping check service health.
func (s *Service) Ping(c context.Context) error {
	return s.dao.Ping(c)
}
