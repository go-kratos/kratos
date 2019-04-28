package service

import (
	"context"

	"go-common/app/admin/ep/saga/conf"
	"go-common/app/admin/ep/saga/dao"
	"go-common/app/admin/ep/saga/service/gitlab"
	"go-common/app/admin/ep/saga/service/wechat"

	"github.com/robfig/cron"
)

// Service struct
type Service struct {
	dao    *dao.Dao
	gitlab *gitlab.Gitlab
	git    *gitlab.Gitlab
	cron   *cron.Cron
	wechat *wechat.Wechat
}

// New a DirService and return.
func New() (s *Service) {
	var (
		err error
	)
	s = &Service{
		dao:  dao.New(),
		cron: cron.New(),
	}
	if err = s.cron.AddFunc(conf.Conf.Property.SyncProject.CheckCron, s.collectprojectproc); err != nil {
		panic(err)
	}
	if err = s.cron.AddFunc(conf.Conf.Property.Git.CheckCron, s.alertProjectPipelineProc); err != nil {
		panic(err)
	}

	if err = s.cron.AddFunc(conf.Conf.Property.SyncData.CheckCron, s.syncdataproc); err != nil {
		panic(err)
	}
	if err = s.cron.AddFunc(conf.Conf.Property.SyncData.CheckCronAll, s.syncalldataproc); err != nil {
		panic(err)
	}
	if err = s.cron.AddFunc(conf.Conf.Property.SyncData.CheckCronWeek, s.syncweekdataproc); err != nil {
		panic(err)
	}
	s.cron.Start()

	// init gitlab client
	s.gitlab = gitlab.New(conf.Conf.Property.Gitlab.API, conf.Conf.Property.Gitlab.Token)
	// init online gitlab client
	s.git = gitlab.New(conf.Conf.Property.Git.API, conf.Conf.Property.Git.Token)
	// init wechat client
	s.wechat = wechat.New(s.dao)

	return
}

// Ping check dao health.
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Wait wait all closed.
func (s *Service) Wait() {
}

// Close close all dao.
func (s *Service) Close() {
	s.dao.Close()
}
