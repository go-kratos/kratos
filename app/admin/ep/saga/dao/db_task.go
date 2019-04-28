package dao

import (
	"go-common/app/admin/ep/saga/model"

	pkgerr "github.com/pkg/errors"
)

// Tasks get all the tasks for the specified project
func (d *Dao) Tasks(projID, status int) (tasks []*model.Task, err error) {
	err = pkgerr.WithStack(d.db.Where(&model.Task{ProjID: projID, Status: status}).Find(&tasks).Error)
	return
}
