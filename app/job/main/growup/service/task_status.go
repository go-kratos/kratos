package service

import (
	"context"
	"fmt"
	"time"

	"go-common/app/job/main/growup/dao"
	"go-common/app/job/main/growup/dao/dataplatform"
)

const (
	// TaskAvCharge .
	TaskAvCharge = iota + 1
	// TaskCmCharge .
	TaskCmCharge
	// TaskTagRatio .
	TaskTagRatio
	// TaskBubbleMeta .
	TaskBubbleMeta
	// TaskBlacklist .
	TaskBlacklist
	// TaskCreativeIncome .
	TaskCreativeIncome
	// TaskCreativeStatis .
	TaskCreativeStatis
	// TaskBgmSync .
	TaskBgmSync
	// TaskTagIncome .
	TaskTagIncome
	// TaskCreativeCharge .
	TaskCreativeCharge
	// TaskBudget .
	TaskBudget
	// TaskSnapshotBubbleIncome .
	TaskSnapshotBubbleIncome
)

const (
	_taskSuccess = 1
	_taskFail    = 2
)

var taskSvr *taskService

type taskService struct {
	dao *dao.Dao
	dp  *dataplatform.Dao
}

// GetTaskService get task service
func GetTaskService() *taskService {
	return taskSvr
}

// TaskReady is task ready
func (s *taskService) TaskReady(c context.Context, date string, typs ...int) (err error) {
	t, err := time.ParseInLocation("2006-01-02", date, time.Local)
	if err != nil {
		return
	}
	ok, err := s.checkBasicStatus(c, t)
	if err != nil {
		return
	}
	if !ok {
		err = fmt.Errorf("basic task not ready yet: %s", date)
		return
	}
	for _, typ := range typs {
		var status int
		status, err = s.dao.TaskStatus(c, date, typ)
		if err != nil {
			return
		}
		if status != 1 {
			err = fmt.Errorf("task(%d) not ready yet: %s", typ, date)
			return
		}
	}
	return
}

func (s *taskService) setTaskSuccess(c context.Context, typ int, date, message string) (rows int64, err error) {
	return s.dao.InsertTaskStatus(c, typ, _taskSuccess, date, message)
}

func (s *taskService) setTaskFail(c context.Context, typ int, date, message string) (rows int64, err error) {
	return s.dao.InsertTaskStatus(c, typ, _taskFail, date, message)
}

// SetTaskStatus set task status by error
func (s *taskService) SetTaskStatus(c context.Context, typ int, date string, err error) (int64, error) {
	if err != nil {
		return s.setTaskFail(c, typ, date, err.Error())
	}
	return s.setTaskSuccess(c, typ, date, "success")
}

// checkBasicStatus check basic date status
func (s *taskService) checkBasicStatus(c context.Context, date time.Time) (ok bool, err error) {
	return s.dp.SendBasicDataRequest(c, fmt.Sprintf("{\"select\": [], \"where\":{\"job_name\":{\"in\":[\"ucs_%s\"]}}}", date.Format("20060102")))
}

// UpdateTaskStatus update task status
func (s *Service) UpdateTaskStatus(c context.Context, date string, typ int, status int) (err error) {
	_, err = s.dao.UpdateTaskStatus(c, date, typ, status)
	return
}
