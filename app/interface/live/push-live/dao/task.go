package dao

import (
	"context"
	"go-common/app/interface/live/push-live/model"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_createNewTask = "INSERT INTO ap_push_task(type,target_id,alert_title,alert_body,mid_source,link_type,link_value,expire_time,total) VALUES (?,?,?,?,?,?,?,?,?)"
)

// CreateNewTask 新增推送任务记录
func (d *Dao) CreateNewTask(c context.Context, task *model.ApPushTask) (affected int64, err error) {
	res, err := d.db.Exec(c, _createNewTask, model.LivePushType, task.TargetID, task.AlertTitle,
		task.AlertBody, task.MidSource, task.LinkType, task.LinkValue, task.ExpireTime, task.Total)
	if err != nil {
		err = errors.WithStack(err)
		log.Error("[dao.task|CreateNewTask] db.Exec() error(%v)", err)
		return
	}
	return res.RowsAffected()
}
