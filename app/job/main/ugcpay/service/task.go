package service

import (
	"context"
	"runtime/debug"

	"go-common/app/job/main/ugcpay/dao"
	"go-common/app/job/main/ugcpay/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"

	"github.com/pkg/errors"
)

func (s *Service) wrapDisProc(tp TaskProcess) func() {
	return func() {
		defer func() {
			if x := recover(); x != nil {
				log.Error("task : %s, panic(%+v): %s", tp.Name(), x, debug.Stack())
			}
		}()
		var (
			ok  bool
			err error
		)
		if ok, err = s.taskCreate(tp.Name(), tp.TTL()); err != nil {
			log.Info("s.taskCreate err: %+v", err)
			return
		}
		if !ok {
			log.Info("task : %s end, other task is running", tp.Name())
			return
		}
		defer func() {
			if err = s.taskDone(tp.Name()); err != nil {
				log.Error("task : %s, taskDone error: %+v", tp.Name(), err)
			}
		}()
		log.Info("task : %s, task start", tp.Name())
		if err = tp.Run(); err != nil {
			log.Error("task : %s end, error: %+v", tp.Name(), err)
		}
	}
}

// TaskProcess .
type TaskProcess interface {
	Run() error   // 运行任务
	TTL() int32   // 任务的最长生命周期
	Name() string // 任务名称
}

func (s *Service) taskCreate(task string, ttl int32) (ok bool, err error) {
	log.Info("task create: %s, ttl: %d", task, ttl)
	return s.dao.AddCacheTask(context.Background(), task, ttl)
}

func (s *Service) taskDone(task string) (err error) {
	// return s.dao.DelCacheTask(context.Background(), task)
	return
}

func checkOrCreateTaskFromLog(ctx context.Context, task TaskProcess, tl *taskLog, expectFN func(context.Context) (int64, error)) (finished bool, err error) {
	var (
		taskCreated bool
		expect      int64
	)
	if taskCreated, finished = tl.checkTask(task); finished {
		log.Info("%s already finished", task.Name())
		return
	}
	if !taskCreated {
		if expect, err = expectFN(ctx); err != nil {
			return
		}
		if _, err = tl.createTask(ctx, task, expect); err != nil {
			return
		}
	}
	return
}

func runTXCASTaskWithLog(ctx context.Context, task TaskProcess, tl *taskLog, biz func(context.Context, *xsql.Tx) (bool, error)) (err error) {
	fn := func(ctx context.Context) (affected bool, err error) {
		affected = true
		tx, err := tl.d.BeginTran(ctx)
		if err != nil {
			return
		}
		if affected, err = biz(ctx, tx); err != nil {
			// 业务报错，不主动rollback
			return
		}
		if err = tl.recordTaskSuccess(ctx, tx, task); err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
		return
	}
	if err = runCAS(ctx, fn); err != nil {
		tl.recordTaskFailure(ctx, task)
	}
	return
}

type taskLog struct {
	d *dao.Dao
}

func (t *taskLog) createTask(ctx context.Context, task TaskProcess, expect int64) (logTask *model.LogTask, err error) {
	logTask = &model.LogTask{
		Name:   task.Name(),
		Expect: expect,
		State:  "created",
	}
	logTask.ID, err = t.d.InsertLogTask(ctx, logTask)
	return
}

func (t *taskLog) recordTaskSuccess(ctx context.Context, tx *xsql.Tx, task TaskProcess) (err error) {
	_, err = t.d.TXIncrLogTaskSuccess(ctx, tx, task.Name())
	if err != nil {
		err = errors.Wrapf(err, "taskLog recordTaskSuccess: %s", task.Name())
	}
	return
}

func (t *taskLog) recordTaskFailure(ctx context.Context, task TaskProcess) {
	_, err := t.d.IncrLogTaskFailure(ctx, task.Name())
	if err != nil {
		err = errors.Wrapf(err, "taskLog recordTaskFailure: %s", task.Name())
		log.Error("%+v", err)
	}
}

func (t *taskLog) checkTask(task TaskProcess) (created, finished bool) {
	data, err := t.d.LogTask(ctx, task.Name())
	if err != nil {
		return
	}
	if data == nil {
		return
	}
	log.Info("checkTask: %s, data: %+v", task.Name(), data)
	created = true
	if data.State == "success" {
		finished = true
		return
	}
	if data.Expect == data.Success {
		finished = true
	}
	return
}
