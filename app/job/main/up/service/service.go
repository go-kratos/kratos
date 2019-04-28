package service

import (
	"context"
	"sync"
	"time"

	"go-common/app/interface/main/mcn/tool/worker"
	"go-common/app/job/main/up/conf"
	"go-common/app/job/main/up/dao/account"
	"go-common/app/job/main/up/dao/email"
	"go-common/app/job/main/up/dao/upcrm"
	archive "go-common/app/service/main/archive/api"
	upGRPCv1 "go-common/app/service/main/up/api/v1"
	"go-common/library/queue/databus"

	"github.com/robfig/cron"
	"go-common/app/admin/main/up/util/databusutil"
)

// Service struct
type Service struct {
	c              *conf.Config
	maildao        *email.Dao
	crmdb          *upcrm.Dao
	acc            *account.Dao
	arcRPC         archive.ArchiveClient
	cron           *cron.Cron
	worker         *worker.Pool
	wg             sync.WaitGroup
	archiveNotifyT *databus.Databus
	archiveT       *databus.Databus
	closeCh        chan struct{}

	upRPC upGRPCv1.UpClient

	databusHandler *databusutil.DatabusHandler
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:       c,
		cron:    cron.New(),
		crmdb:   upcrm.New(c),
		acc:     account.New(c),
		maildao: email.New(c),
		worker: worker.New(&worker.Conf{
			WorkerProcMax: 10,
			QueueSize:     1024,
			WorkerNumber:  4}),
		archiveNotifyT: databus.New(c.DatabusConf.ArchiveNotify),
		archiveT:       databus.New(c.DatabusConf.Archive),
		closeCh:        make(chan struct{}),
		databusHandler: databusutil.NewDatabusHandler(),
	}
	var err error
	s.arcRPC, err = archive.NewClient(c.GRPCClient.Archive)
	if err != nil {
		panic(err)
	}
	if err = s.initEmailTemplate(); err != nil {
		panic(err)
	}

	if s.upRPC, err = upGRPCv1.NewClient(c.GRPCClient.Up); err != nil {
		panic(err)
	}

	s.createJobs()
	s.databusHandler.GoWatch(s.archiveNotifyT, s.handleArchiveNotifyT)
	s.databusHandler.GoWatch(s.archiveT, s.handleArchiveT)
	return s
}

func (s *Service) createJobs() {
	s.cron.AddFunc(conf.Conf.Job.UpCheckDateDueTaskTime, cronWrap(s.CheckDateDueJob))
	s.cron.AddFunc(conf.Conf.Job.TaskScheduleTime, cronWrap(s.CheckTaskJob))
	s.cron.AddFunc(conf.Conf.Job.CheckStateJobTime, cronWrap(s.CheckStateJob))
	s.cron.AddFunc(conf.Conf.Job.UpdateUpTidJobTime, cronWrap(s.UpdateUpTidJob))
	s.cron.Start()
}

func cronWrap(f func(tm time.Time)) func() {
	return func() {
		f(time.Now())
	}
}

// Ping Service
func (s *Service) Ping(c context.Context) (err error) {
	return s.crmdb.Ping(c)
}

// Close Service
func (s *Service) Close() {
	s.databusHandler.Close()
	s.wg.Wait()
	s.crmdb.Close()
}
