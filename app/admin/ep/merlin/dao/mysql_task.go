package dao

import (
	"go-common/app/admin/ep/merlin/model"

	pkgerr "github.com/pkg/errors"
)

// FindDeleteMachineTasks find delete machine tasks.
func (d *Dao) FindDeleteMachineTasks() (tasks []*model.Task, err error) {
	tasks = []*model.Task{}
	err = pkgerr.WithStack(d.db.Where("status=? AND type=? AND DATEDIFF(NOW(),execute_time)=0 AND execute_time < NOW()", model.TaskInit, model.DeleteMachine).Find(&tasks).Error)
	return
}

// UpdateTaskStatusByMachines update task status by machines.
func (d *Dao) UpdateTaskStatusByMachines(machineIDs []int64, status int) (err error) {
	return pkgerr.WithStack(d.db.Model(&model.Task{}).Where("machine_id IN (?)", machineIDs).Update("status", status).Error)
}

// UpdateTaskStatusByTaskID update task status by taskId.
func (d *Dao) UpdateTaskStatusByTaskID(taskID int64, status int) (err error) {
	return pkgerr.WithStack(d.db.Model(&model.Task{}).Where("id IN (?)", taskID).Update("status", status).Error)
}

// InsertDeleteMachinesTasks insert delete machines tasks.
// TODO:这块逻辑不是很好，如果一个写入失败希望继续写入失败的话最好返回失败list，后续联调时优化
func (d *Dao) InsertDeleteMachinesTasks(ms []*model.Machine) (err error) {
	for _, machine := range ms {
		mrTask := &model.Task{
			TYPE:        model.DeleteMachine,
			ExecuteTime: machine.EndTime,
			MachineID:   machine.ID,
			Status:      model.TaskInit,
		}
		if tempErr := d.db.Create(mrTask).Error; tempErr != nil {
			err = pkgerr.WithStack(tempErr)
		}
	}
	return
}
