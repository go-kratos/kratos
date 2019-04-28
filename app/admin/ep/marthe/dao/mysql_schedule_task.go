package dao

import (
	"go-common/app/admin/ep/marthe/model"

	pkgerr "github.com/pkg/errors"
)

// InsertScheduleTask Insert Schedule Task.
func (d *Dao) InsertScheduleTask(scheduleTask *model.ScheduleTask) error {
	return pkgerr.WithStack(d.db.Create(scheduleTask).Error)
}

// UpdateScheduleTask Update Schedule Task.
func (d *Dao) UpdateScheduleTask(scheduleTask *model.ScheduleTask) error {
	return pkgerr.WithStack(d.db.Model(&model.ScheduleTask{}).Updates(scheduleTask).Error)
}
