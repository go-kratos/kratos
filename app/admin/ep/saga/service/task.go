package service

import (
	"context"

	"go-common/app/admin/ep/saga/model"
)

// MergeTasks query all tasks for the project.
func (s *Service) MergeTasks(c context.Context, req *model.TasksReq) (resp *model.TasksResp, err error) {
	var (
		tasks  []*model.Task
		status int
	)
	resp = new(model.TasksResp)
	if _, err = s.dao.ProjectInfoByID(req.ProjID); err != nil {
		return
	}

	for _, status = range req.Statuses {
		if tasks, err = s.dao.Tasks(req.ProjID, status); err != nil {
			return
		}
		resp.Tasks = append(resp.Tasks, tasks...)
	}

	return
}
