package service

import (
	"context"
	"time"

	filtermdl "go-common/app/service/main/filter/model/rpc"
	"go-common/app/service/main/push/dao"
	"go-common/app/service/main/push/model"
	xsql "go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
)

// pickNewTask get a new task by tx.
func (s *Service) pickNewTask(platform int) (t *model.Task, err error) {
	c := context.Background()
	var tx *xsql.Tx
	if tx, err = s.dao.BeginTx(c); err != nil {
		log.Error("tx.BeginTx() error(%v)", err)
		return
	}
	if t, err = s.dao.TxTaskByPlatform(tx, platform); err != nil {
		if e := tx.Rollback(); e != nil {
			dao.PromError("task:获取新任务")
			log.Error("tx.Rollback() error(%v)", e)
		}
		return
	}
	if t == nil {
		if e := tx.Rollback(); e != nil {
			dao.PromError("task:获取新任务")
			log.Error("tx.Rollback() error(%v)", e)
		}
		return
	}
	if err = s.dao.TxUpdateTaskStatus(tx, t.ID, model.TaskStatusDoing); err != nil {
		if e := tx.Rollback(); e != nil {
			dao.PromError("task:更新任务状态")
			log.Error("tx.Rollback() error(%v)", e)
		}
		return
	}
	if err = tx.Commit(); err != nil {
		dao.PromError("task:获取新任务commit")
		log.Error("tx.Commit() error(%v)", err)
		return
	}
	s.setProgress(t.ID, _pgStatus, int64(model.TaskStatusDoing))
	return
}

// Task gets task info.
func (s *Service) Task(c context.Context, businessID int64, taskID, token string) (task *model.Task, err error) {
	if _, ok := s.businesses[businessID]; !ok {
		err = ecode.RequestErr
		log.Error("no business: %d", businessID)
		return
	}
	business := s.businesses[businessID]
	if token != business.Token {
		err = ecode.RequestErr
		log.Error("mismatch business token: %s, expected: %s", token, s.businesses[businessID])
		return
	}
	if task, err = s.dao.Task(c, taskID); err != nil {
		return
	}
	if task != nil {
		if p := s.fetchProgress(taskID); p != nil {
			task.Progress = p
		}
	}
	return
}

func (s *Service) handleTask(task *model.Task) (err error) {
	if time.Now().Unix() > int64(task.ExpireTime) {
		log.Warn("task(%s) expired", task.ID)
		dao.PromInfo("task:任务过期")
		err = s.updateTaskStatus(task.ID, model.TaskStatusExpired)
		return
	}
	s.setProgress(task.ID, _pgBeginTime, time.Now().Unix())
	if err = s.pushTokens(task); err != nil {
		err = s.updateTaskStatus(task.ID, model.TaskStatusFailed)
		return
	}
	err = s.updateTaskStatus(task.ID, model.TaskStatusDone)
	return
}

func (s *Service) updateTaskStatus(ID string, status int8) (err error) {
	if err = s.dao.UpdateTaskStatus(context.Background(), ID, status); err != nil {
		dao.PromError("task:更新任务状态")
		log.Error("s.updateTaskStatus(%d,%d) error(%v)", ID, status, err)
		return
	}
	if status == model.TaskStatusDone || status == model.TaskStatusFailed || status == model.TaskStatusExpired {
		s.setProgress(ID, _pgEndTime, time.Now().Unix())
	}
	s.setProgress(ID, _pgStatus, int64(status))
	return
}

// Filter filter content.
func (s *Service) Filter(c context.Context, content string) (res string, err error) {
	var (
		filterRes *filtermdl.FilterRes
		arg       = filtermdl.ArgFilter{Area: "common", Message: content}
	)
	if filterRes, err = s.filterRPC.Filter(c, &arg); err != nil {
		dao.PromError("push:过滤服务")
		log.Error("s.filter(%s) error(%v)", content, err)
		return
	}
	if filterRes.Level < 20 {
		return content, nil
	}
	res = filterRes.Result
	return
}
