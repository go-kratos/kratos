package service

import (
	"context"
	"time"

	"go-common/app/service/main/antispam/conf"
	"go-common/library/log"
)

// Scheduler .
type Scheduler interface {
	BuildTrie()
	RefreshTrie()
	ExpireKeyword()
	RefreshRules()
	RefreshRegexps()
	RunTimeDebugProb()
}

// SchedulerImpl .
type SchedulerImpl struct {
	service                Service
	trieMgr                *TrieMgr
	refreshTrieInterval    time.Duration
	refreshRulesInterval   time.Duration
	refreshRegexpsInterval time.Duration

	buildTrieInterval        time.Duration
	buildTrieMaxRowsPerQuery int64

	expireKeywordInterval        time.Duration
	expireKeywordMaxRowsPerQuery int64
}

// NewScheduler .
func NewScheduler(service Service, trieMgr *TrieMgr, opt *SchedulerOption) Scheduler {
	return &SchedulerImpl{
		service:                      service,
		trieMgr:                      trieMgr,
		refreshTrieInterval:          time.Second * opt.RefreshTrieIntervalSec,
		refreshRulesInterval:         time.Second * opt.RefreshRulesIntervalSec,
		refreshRegexpsInterval:       time.Second * opt.RefreshRegexpsIntervalSec,
		buildTrieInterval:            time.Minute * opt.BuildTrieIntervalMinute,
		buildTrieMaxRowsPerQuery:     opt.BuildTrieMaxRowsPerQuery,
		expireKeywordInterval:        time.Second * opt.GcInterval,
		expireKeywordMaxRowsPerQuery: opt.GcMaxRowsPerQuery,
	}
}

func schedule(name string, dur time.Duration, op func()) {
	for {
		select {
		case <-time.After(dur):
			log.Info("start %s...", name)
			op()
		case <-done:
			log.Info("%s exit ...", name)
			return
		}
	}
}

// RunTimeDebugProb .
func (s *SchedulerImpl) RunTimeDebugProb() {
	schedule("runtime debug prob",
		time.Second*300,
		func() {
			for _, r := range regexps {
				log.Info("regexps:%+v", r)
			}
			for _, r := range rules {
				log.Info("rules:%+v", r)
			}
			log.Info("autowhite config: %+v", conf.Conf.AutoWhite)
		})
}

// RefreshRules .
func (s *SchedulerImpl) RefreshRules() {
	s.service.RefreshRules(context.TODO())
	schedule("refresh rule",
		s.refreshRulesInterval,
		func() { s.service.RefreshRules(context.TODO()) })
}

// RefreshRegexps .
func (s *SchedulerImpl) RefreshRegexps() {
	s.service.RefreshRegexps(context.TODO())
	schedule("refresh regexp",
		s.refreshRegexpsInterval,
		func() { s.service.RefreshRegexps(context.TODO()) })
}

// RefreshTrie .
func (s *SchedulerImpl) RefreshTrie() {
	schedule("refresh trie",
		s.refreshTrieInterval,
		s.trieMgr.Refresh)
}

// BuildTrie .
func (s *SchedulerImpl) BuildTrie() {
	schedule("build trie",
		s.buildTrieInterval,
		func() { s.trieMgr.Build(s.buildTrieMaxRowsPerQuery) })
}

// ExpireKeyword .
func (s *SchedulerImpl) ExpireKeyword() {
	schedule("expire keyword",
		s.expireKeywordInterval,
		func() { s.service.ExpireKeyword(context.TODO(), s.expireKeywordMaxRowsPerQuery) })
}
