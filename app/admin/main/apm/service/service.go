package service

import (
	"context"
	"go-common/app/admin/main/apm/conf"
	"go-common/app/admin/main/apm/dao"
	"go-common/app/admin/main/apm/model/tree"
	"go-common/app/admin/main/apm/model/ut"
	"go-common/app/tool/saga/service/gitlab"
	bm "go-common/library/net/http/blademaster"
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/robfig/cron"
)

// Service is a service.
type Service struct {
	c         *conf.Config
	dao       *dao.Dao
	DB        *gorm.DB
	DBDatabus *gorm.DB
	DBCanal   *gorm.DB
	client    *bm.Client
	// tree cache
	treeCache map[string][]*tree.Node
	treeLock  sync.RWMutex
	// cron cron
	cron *cron.Cron
	// discoveryID cache
	discoveryIDCache map[string]*tree.Resd
	discoveryIDLock  sync.RWMutex
	ranksCache       *ut.RanksCache
	appsCache        *ut.AppsCache
	// dapper proxy
	dapperProxy *dapperProxy
	// gitlab api conf
	gitlab *gitlab.Gitlab
}

// New new a service
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:      c,
		dao:    dao.New(c),
		client: bm.NewClient(c.HTTPClient),
		// tree cache
		treeCache: map[string][]*tree.Node{},
		// discoveryID cache
		discoveryIDCache: map[string]*tree.Resd{},
		// ranks cache
		ranksCache: &ut.RanksCache{},
		appsCache:  &ut.AppsCache{},
		// cron cron
		cron: cron.New(),
	}
	s.gitlab = gitlab.New(conf.Conf.Gitlab.API, conf.Conf.Gitlab.Token)
	s.DB = s.dao.DB
	s.DBDatabus = s.dao.DBDatabus
	s.DBCanal = s.dao.DBCanal
	if err := s.cron.AddFunc(s.c.Cron.Crontab, s.taskAddMonitor); err != nil {
		panic(err)
	}
	if err := s.cron.AddFunc(s.c.Cron.Crontab, s.taskAddCache); err != nil {
		panic(err)
	}
	if err := s.cron.AddFunc(s.c.Cron.CrontabRepo, s.taskRankWechatReport); err != nil {
		panic(err)
	}
	if err := s.cron.AddFunc(s.c.Cron.CrontabRepo, s.taskWeeklyWechatReport); err != nil {
		panic(err)
	}
	dp, err := newDapperProxy(c.Host.DapperCo)
	if err != nil {
		panic(err)
	}
	s.dapperProxy = dp
	s.cron.Start()
	go s.taskAddCache()
	return
}

// Ping ping db,
func (s *Service) Ping(c context.Context) (err error) {
	return
}

// Close close resource.
func (s *Service) Close() {
	s.dao.Close()
}
