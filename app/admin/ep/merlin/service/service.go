package service

import (
	"context"

	"go-common/app/admin/ep/merlin/conf"
	"go-common/app/admin/ep/merlin/dao"
	"go-common/library/sync/pipeline/fanout"

	"github.com/robfig/cron"
)

// Service struct
type Service struct {
	c          *conf.Config
	dao        *dao.Dao
	cron       *cron.Cron
	deviceChan *fanout.Fanout
}

// New init.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:          c,
		dao:        dao.New(c),
		deviceChan: fanout.New("deviceChan", fanout.Worker(1), fanout.Buffer(1024)),
	}
	scheduler := c.Scheduler
	if scheduler.Active {
		s.cron = cron.New()
		if err := s.cron.AddFunc(scheduler.GetExpiredMachinesTime, s.taskGetExpiredMachinesIntoTask); err != nil {
			panic(err)
		}
		if err := s.cron.AddFunc(scheduler.SendTaskMailMachinesWillExpiredTime, s.taskSendTaskMailMachinesWillExpired); err != nil {
			panic(err)
		}
		if err := s.cron.AddFunc(scheduler.DeleteExpiredMachinesInTask, s.taskDeleteExpiredMachines); err != nil {
			panic(err)
		}
		if err := s.cron.AddFunc(scheduler.CheckMachinesStatusInTask, s.taskMachineStatus); err != nil {
			panic(err)
		}
		if err := s.cron.AddFunc(scheduler.UpdateMobileDeviceInTask, s.taskSyncMobileDeviceList); err != nil {
			panic(err)
		}
		if err := s.cron.AddFunc(scheduler.UpdateSnapshotStatusInDoing, s.taskUpdateSnapshotStatusInDoing); err != nil {
			panic(err)
		}
		s.cron.Start()
	}
	return
}

// Close Service.
func (s *Service) Close() {
	s.dao.Close()
}

// Ping check server ok.
func (s *Service) Ping(c context.Context) (err error) {
	err = s.dao.Ping(c)
	return
}

// ConfVersion Conf Version.
func (s *Service) ConfVersion(c context.Context) string {
	return conf.Conf.Version
}
