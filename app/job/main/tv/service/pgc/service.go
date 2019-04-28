package pgc

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"go-common/app/job/main/tv/conf"
	"go-common/app/job/main/tv/dao/app"
	"go-common/app/job/main/tv/dao/cms"
	"go-common/app/job/main/tv/dao/lic"
	playdao "go-common/app/job/main/tv/dao/playurl"
	model "go-common/app/job/main/tv/model/pgc"
	"go-common/library/log"
	"go-common/library/queue/databus"

	"go-common/app/job/main/tv/dao/ftp"

	"github.com/robfig/cron"
)

var ctx = context.Background()

// Service struct of service.
type Service struct {
	dao            *app.Dao
	daoClosed      bool // logic close the dao's DB
	playurlDao     *playdao.Dao
	licDao         *lic.Dao
	ftpDao         *ftp.Dao
	cmsDao         *cms.Dao
	c              *conf.Config
	waiter         *sync.WaitGroup // general waiter
	waiterConsumer *sync.WaitGroup
	contentSub     *databus.Databus // consumer for state change
	cron           *cron.Cron
	ResuEps        []*model.Content
	ResuSns        []*model.TVEpSeason
	resuRetry      map[string]int
}

// New create service instance and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:              c,
		dao:            app.New(c),
		playurlDao:     playdao.New(c),
		licDao:         lic.New(c),
		ftpDao:         ftp.New(c),
		cmsDao:         cms.New(c),
		daoClosed:      false,
		waiter:         new(sync.WaitGroup),
		waiterConsumer: new(sync.WaitGroup),
		contentSub:     databus.New(c.ContentSub),
		cron:           cron.New(),
		resuRetry:      make(map[string]int),
	}
	rand.Seed(time.Now().UnixNano())
	// flush Redis - zone list
	go s.ZoneIdx()
	if err := s.cron.AddFunc(s.c.Redis.CronPGC, s.ZoneIdx); err != nil {
		panic(err)
	}
	if err := s.cron.AddFunc(s.c.PlayControl.ProducerCron, s.refreshCache); err != nil {
		panic(err)
	}
	if err := s.cron.AddFunc(s.c.Cfg.Merak.Cron, s.cmsShelve); err != nil {
		panic(err)
	}
	s.cron.Start()
	go s.searchSugproc()  // uploads the passed season's list to search sug's FTP
	go s.seaPgcContproc() // uploads pgc search content to sug's FTP
	s.waiter.Add(1)
	go s.syncEPs()
	s.waiter.Add(1)
	go s.resubEps()
	s.waiter.Add(1)
	go s.resubSns()
	s.waiter.Add(1)
	go s.syncSeason()
	s.waiter.Add(1)
	go s.delSeason()
	s.waiter.Add(1)
	go s.delCont()
	// Databus
	s.waiterConsumer.Add(1)
	go s.consumeContent() // consume Databus Message to update MC
	return
}

// Close dao.
func (s *Service) Close() {
	if s.dao != nil {
		s.daoClosed = true
		log.Info("Dao Closed!")
	}
	log.Info("Crontab Closed!")
	s.cron.Stop()
	log.Info("Databus Closed!")
	s.contentSub.Close()
	log.Info("Wait Producer!")
	s.waiter.Wait()
	log.Info("Wait SyncMC Consumers")
	s.waiterConsumer.Wait()
	log.Info("Physical Dao Closed!")
	s.dao.Close()
	log.Info("tv-job has been closed.")
}
