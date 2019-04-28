package service

import (
	"context"
	"sync"

	"go-common/app/admin/ep/marthe/conf"
	"go-common/app/admin/ep/marthe/dao"
	"go-common/app/admin/ep/marthe/model"
	"go-common/library/sync/pipeline/fanout"

	"github.com/robfig/cron"
)

// Service struct
type Service struct {
	c                      *conf.Config
	dao                    *dao.Dao
	batchRunCache          *fanout.Fanout
	tapdBugCache           *fanout.Fanout
	taskCache              *fanout.Fanout
	cron                   *cron.Cron
	syncWechatContactsLock sync.Mutex

	syncTapdBugInsertLock sync.Mutex
	mapTapdBugInsertLocks map[int64]*sync.Mutex

	syncBatchRunLock sync.Mutex
	mapBatchRunLocks map[int64]*sync.Mutex
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:                     c,
		dao:                   dao.New(c),
		batchRunCache:         fanout.New("batchRunCache", fanout.Worker(5), fanout.Buffer(10240)),
		tapdBugCache:          fanout.New("tapdBugCache", fanout.Worker(5), fanout.Buffer(10240)),
		taskCache:             fanout.New("taskCache", fanout.Worker(5), fanout.Buffer(10240)),
		mapTapdBugInsertLocks: make(map[int64]*sync.Mutex),
		mapBatchRunLocks:      make(map[int64]*sync.Mutex),
	}

	if c.Scheduler.Active {
		s.cron = cron.New()

		// 定时批量 跑enable version 抓bugly数据
		if err := s.cron.AddFunc(c.Scheduler.BatchRunEnableVersion, func() { s.BatchRunTask(model.TaskBatchRunVersions, s.BatchRunVersions) }); err != nil {
			panic(err)
		}

		// 定时把超过三小时为执行完毕的任务修改为失败
		if err := s.cron.AddFunc(c.Scheduler.DisableBatchRunOverTime, func() { s.BatchRunTask(model.TaskDisableBatchRunOverTime, s.DisableBatchRunOverTime) }); err != nil {
			panic(err)
		}

		// 定时更新tapd bug
		if err := s.cron.AddFunc(c.Scheduler.BatchRunUpdateTapdBug, func() { s.BatchRunTask(model.TaskBatchRunUpdateBugInTapd, s.BatchRunUpdateBugInTapd) }); err != nil {
			panic(err)
		}

		// 定时更新SyncWechatContact
		if err := s.cron.AddFunc(c.Scheduler.SyncWechatContact, func() { s.BatchRunTask(model.TaskSyncWechatContact, s.SyncWechatContacts) }); err != nil {
			panic(err)
		}

		s.cron.Start()
	}
	return s
}

// Ping Service
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close Service
func (s *Service) Close() {
	s.dao.Close()
}

// BatchRunTask Batch Run Task.
func (s *Service) BatchRunTask(taskName string, task func() error) {
	var err error
	scheduleTask := &model.ScheduleTask{
		Name:   taskName,
		Status: model.TaskStatusRunning,
	}
	if err = s.dao.InsertScheduleTask(scheduleTask); err != nil {
		return
	}

	err = task()

	defer func() {
		if err != nil {
			scheduleTask.Status = model.TaskStatusFailed
		} else {
			scheduleTask.Status = model.TaskStatusDone
		}

		if err = s.dao.UpdateScheduleTask(scheduleTask); err != nil {
			return
		}
	}()
}
