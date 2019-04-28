package service

import (
	"context"
	"time"

	"go-common/app/job/main/aegis/model"
	"go-common/library/log"
	"go-common/library/queue/databus/report"
)

// send to log service
func (s *Service) sendTaskLog(c context.Context, task *model.Task, tp int, action string, uid int64, uname string) (err error) {
	logData := &report.ManagerInfo{
		UID:      uid,
		Uname:    uname,
		Business: model.LogBusinessTask,
		Type:     tp,
		Oid:      task.ID,
		Action:   action,
		Ctime:    time.Now(),
		Index:    []interface{}{task.BusinessID, task.FlowID, task.State},
		Content: map[string]interface{}{
			"task": task,
		},
	}
	err = report.Manager(logData)
	log.Info("sendTaskLog logData(%+v) errmsg(%v)", logData, err)
	return
}

func (s *Service) sendWeightLog(c context.Context, task *model.Task, wl *model.WeightLog) (err error) {
	logData := &report.ManagerInfo{
		UID:      399,
		Uname:    "aegis-job",
		Business: model.LogBusinessTask,
		Type:     model.LogTYpeTaskWeight,
		Oid:      task.ID,
		Action:   "weight",
		Ctime:    time.Now(),
		Index:    []interface{}{task.BusinessID, task.FlowID, task.State},
		Content: map[string]interface{}{
			"weightlog": wl,
		},
	}
	err = report.Manager(logData)
	log.Info("sendWeightLog logData(%+v) errmsg(%v)", logData, err)
	return
}
