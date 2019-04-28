package service

import (
	"context"
	"math/rand"
	"strconv"
	"sync"
	"time"

	filterrpc "go-common/app/service/main/filter/rpc/client"
	"go-common/app/service/main/push/conf"
	"go-common/app/service/main/push/dao"
	"go-common/app/service/main/push/model"
	"go-common/library/cache"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// Service push service.
type Service struct {
	c              *conf.Config
	dao            *dao.Dao
	cache          *cache.Cache
	waiter         sync.WaitGroup
	filterRPC      *filterrpc.Service
	reportCache    *cache.Cache
	progress       map[string]*model.Progress
	businesses     map[int64]*model.Business
	apnsCh         chan *model.PushChanItem
	miCh, huaweiCh chan *model.PushChanItems
	jpushCh        chan *model.PushChanItems
	fcmCh          chan *model.PushChanItems
	oppoCh         chan *model.PushChanItem
	chCounter      map[string]int64
	pmu, ppmu      sync.RWMutex // progress mutex / progress proc mutex
	closed         bool
	httpclient     *bm.Client
	progressCh     chan func()
}

// New creates a push service instance.
func New(c *conf.Config) *Service {
	rand.Seed(time.Now().UnixNano())
	s := &Service{
		c:           c,
		dao:         dao.New(c),
		filterRPC:   filterrpc.New(c.FilterRPC),
		cache:       cache.New(1, 1024000),
		reportCache: cache.New(1, 1024000),
		progress:    make(map[string]*model.Progress),
		businesses:  make(map[int64]*model.Business),
		apnsCh:      make(chan *model.PushChanItem, c.Push.PushChanSizeAPNS),
		miCh:        make(chan *model.PushChanItems, c.Push.PushChanSizeMi),
		huaweiCh:    make(chan *model.PushChanItems, c.Push.PushChanSizeHuawei),
		jpushCh:     make(chan *model.PushChanItems, c.Push.PushChanSizeJpush),
		oppoCh:      make(chan *model.PushChanItem, c.Push.PushChanSizeOppo),
		fcmCh:       make(chan *model.PushChanItems, c.Push.PushChanSizeFCM),
		chCounter:   make(map[string]int64),
		httpclient:  bm.NewClient(c.HTTPClient),
		progressCh:  make(chan func(), c.Push.UpdateTaskProgressProc),
	}
	s.loadBusiness()
	go s.loadBusinessproc()
	go s.loadTaskproc()
	s.waiter.Add(1)
	go s.updateTaskProgressproc()
	for i := 0; i < s.c.Push.PushGoroutinesAPNS; i++ {
		go s.pushAPNSproc()
	}
	for i := 0; i < s.c.Push.PushGoroutinesMi; i++ {
		go s.pushMiproc()
	}
	for i := 0; i < s.c.Push.PushGoroutinesHuawei; i++ {
		go s.pushHuaweiproc()
	}
	for i := 0; i < s.c.Push.PushGoroutinesOppo; i++ {
		go s.pushOppoproc()
	}
	for i := 0; i < s.c.Push.PushGoroutinesJpush; i++ {
		go s.pushJpushproc()
	}
	for i := 0; i < s.c.Push.PushGoroutinesFCM; i++ {
		go s.pushFCMproc()
	}
	for i := 0; i < s.c.Push.UpdateTaskProgressProc; i++ {
		s.waiter.Add(1)
		go s.updateProgressproc()
	}
	return s
}

func (s *Service) loadTaskproc() {

	if !s.c.Push.PickUpTask {
		log.Warn("service do not pick up new tasks from database")
		return
	}

	for _, v := range model.Platforms {
		s.waiter.Add(1)
		go func(platform int) {
			defer s.waiter.Done()
			for {
				if s.closed {
					return
				}
				task, err := s.pickNewTask(platform)
				if err != nil {
					time.Sleep(5 * time.Second)
					continue
				}
				if task != nil {
					s.handleTask(task)
				}
				time.Sleep(time.Duration(s.c.Push.LoadTaskInteval))
			}
		}(v)
	}
}

func (s *Service) updateTaskProgressproc() {
	defer s.waiter.Done()
	for {
		if s.closed {
			close(s.progressCh)
			return
		}
		time.Sleep(time.Duration(s.c.Push.UpdateTaskProgressInteval))
		dao.PromChanLen("progress len", int64(len(s.progress)))
		s.updateTaskProgress()
	}
}

func (s *Service) updateTaskProgress() {
	progress := make(map[string]*model.Progress)
	s.ppmu.Lock()
	for id, p := range s.progress {
		brands := make(map[int]int64)
		for key, data := range p.Brands {
			brands[key] = data
		}
		pNew := *p
		pNew.Brands = brands
		progress[id] = &pNew
	}
	s.ppmu.Unlock()
	for id, p := range progress {
		id := id
		p := p
		if i, _ := strconv.ParseInt(id, 10, 64); i == 0 {
			continue
		}
		s.progressCh <- func() { s.dao.UpdateTaskProgress(context.Background(), id, p) }
		if p.Status == model.TaskStatusDone || p.Status == model.TaskStatusFailed || p.Status == model.TaskStatusExpired {
			s.ppmu.Lock()
			delete(s.progress, id)
			s.ppmu.Unlock()
		}
	}
}

// Close closes service.
func (s *Service) Close() {
	s.closed = true
	s.waiter.Wait()
	close(s.apnsCh)
	close(s.miCh)
	close(s.huaweiCh)
	close(s.oppoCh)
	close(s.jpushCh)
	s.dao.Close()
}

// Ping checks service.
func (s *Service) Ping(c context.Context) (err error) {
	err = s.dao.Ping(c)
	return
}
