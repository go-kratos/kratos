package service

import (
	"context"
	"errors"
	"sync"

	"go-common/app/service/main/antispam/conf"
	"go-common/app/service/main/antispam/dao"
	"go-common/app/service/main/antispam/extern"
	"go-common/library/cache"

	"go-common/library/log"
)

var done = make(chan struct{})

// Ping .
func (s *SvcImpl) Ping(ctx context.Context) error {
	return s.antiDao.Ping(ctx)
}

// Spawns spawn goroutines with waitGroup
func (s *SvcImpl) Spawns(fns ...func()) {
	for _, fn := range fns {
		s.wg.Add(1)
		go func(f func()) {
			defer s.wg.Done()
			f()
		}(fn)
	}
}

// Close close service and all the resources it opens
func (s *SvcImpl) Close() {
	close(done)
	close(s.UserGeneratedContentChan)
	close(s.AsyncTaskChan)

	s.wg.Wait()
	dao.Close()
}

// AddTask add async task
func (s *SvcImpl) AddTask(fn func()) {
	select {
	case s.AsyncTaskChan <- fn:
	default:
		log.Warn("task chan full, will discard operation")
	}
}

// HandleTask perform task
func (s *SvcImpl) HandleTask() {
	for fn := range s.AsyncTaskChan {
		if fn != nil {
			log.Info("receive task ...")
			fn()
		}
	}
	log.Info("async task chan closed ...")
}

// New .
func New(config *conf.Config) *SvcImpl {
	if ok := dao.Init(config); !ok {
		panic(errors.New("init dao fail"))
	}
	s := &SvcImpl{
		wg:     new(sync.WaitGroup),
		Option: NewOption(config),

		antiDao: dao.New(config),

		RegexpDao:  dao.NewRegexpDao(),
		KeywordDao: dao.NewKeywordDao(),
		RuleDao:    dao.NewRuleDao(),
	}
	s.TrieMgr = NewTrieMgr(s)
	s.Scheduler = NewScheduler(s, s.TrieMgr, s.Option.Scheduler)
	s.AsyncTaskChan = make(chan func(), s.Option.AsyncTaskChanSize)
	s.tokens = make(chan struct{}, s.Option.MaxSpawnGoroutines)
	for i := 0; i < cap(s.tokens); i++ {
		s.tokens <- struct{}{}
	}
	s.UserGeneratedContentChan = make(chan UserGeneratedContent, s.Option.DefaultChanSize)
	if config.ServiceOption.GcOpt.Open {
		s.Spawns(s.Scheduler.ExpireKeyword)
	}
	s.TrieMgr.Build(s.Option.Scheduler.BuildTrieMaxRowsPerQuery)
	s.Spawns(
		s.Digest,
		s.HandleTask,
		s.Scheduler.BuildTrie,
		s.Scheduler.RefreshTrie,
		s.Scheduler.RefreshRules,
		s.Scheduler.RefreshRegexps,
		s.Scheduler.RunTimeDebugProb,
	)
	return s
}

// SvcImpl .
type SvcImpl struct {
	tokens chan struct{}
	sync.RWMutex
	wg            *sync.WaitGroup
	Option        *Option
	AsyncTaskChan chan func()

	TrieMgr       *TrieMgr
	Cache         cache.Cache
	ExternHandler extern.Handler
	Scheduler     Scheduler

	antiDao *dao.Dao

	RuleDao                  dao.RuleDao
	RegexpDao                dao.RegexpDao
	KeywordDao               dao.KeywordDao
	UserGeneratedContentChan chan UserGeneratedContent
}
