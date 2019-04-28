package newcomer

import (
	"context"
	"fmt"
	"go-common/app/interface/main/creative/model/newcomer"
	"go-common/library/ecode"
	"go-common/library/log"
)

// H5TaskList for H5 detail task list
func (s *Service) H5TaskList(c context.Context, mid int64, from string) (res *newcomer.H5TaskRewardList, err error) {
	var (
		u            *UserTaskInfo
		tasks        []*newcomer.Task
		taskTypeMap  = make(map[int8][]*newcomer.Task)  // taskType-task
		taskGroupMap = make(map[int64][]*newcomer.Task) // groupID-task
	)

	// get user tasks
	userTasks, err := s.newc.UserTasks(c, mid)
	if err != nil {
		log.Error("TaskList s.newc.UserTasks mid(%d)|error(%v)", mid, err)
		return
	}
	if len(userTasks) == 0 {
		// return：User did not receive the task
		res = &newcomer.H5TaskRewardList{
			TaskReceived: newcomer.NoBindTask,
		}
		return
	}

	// get user info
	u, err = s.getUserTaskInfo(c, mid, userTasks)
	if err != nil {
		return
	}

	// get tasks
	tasks = s.getTasksInfoByType(userTasks, newcomer.DefualtTaskType)
	if len(tasks) == 0 {
		err = ecode.CreativeNewcomerNoTask
		log.Error("TaskList s.GetTaskByType len(tasks) == 0")
		return
	}
	// group by groupID & taskType
	taskGroupMap, taskTypeMap = s.groupByTasks(tasks)

	// add task label & redirect
	s.addLabelRedirect(tasks, from)

	// task_gift
	taskGift, err := s.getTaskGiftData(c, mid, taskTypeMap, newcomer.FromH5)
	if err != nil {
		return
	}

	// task_groups
	tgs, err := s.getTaskGroupData(c, mid, taskGroupMap)
	if err != nil {
		return
	}
	// if userLevel == UserTaskLevel01 , set unlock state
	if u.UserTaskLevel == newcomer.UserTaskLevel01 {
		for _, v := range tgs {
			if v.TaskType == newcomer.AdvancedTaskType {
				v.RewardState = newcomer.RewardUnlock
			}
		}
	}

	// add tips
	s.addTaskGroupTip(tgs)
	s.addGiftTip(taskGift, taskTypeMap)

	res = &newcomer.H5TaskRewardList{
		TaskReceived: newcomer.BindTask,
		TaskGroups:   tgs,
		TaskGift:     taskGift,
	}
	return
}

// addLabelRedirect add label & redirect
func (s *Service) addLabelRedirect(tasks []*newcomer.Task, from string) {
	if len(tasks) == 0 {
		return
	}
	for _, v := range tasks {
		if v == nil {
			continue
		}
		t, ok := s.TaskMapCache[v.ID]
		if !ok {
			continue
		}
		m, ook := newcomer.H5RedirectMap[from][t.TargetType]
		if !ook || len(m) == 0 {
			continue
		}
		v.Label = m[0]
		v.Redirect = m[1]
	}
}

// addGiftTip get gift tip
func (s *Service) addGiftTip(tg []*newcomer.TaskGift, kindMap map[int8][]*newcomer.Task) {
	for _, v := range tg {
		if v == nil {
			continue
		}
		if v.State != newcomer.RewardUnlock {
			if tip, ok := newcomer.GiftTipMap[v.State][v.Type]; ok {
				v.Tip = tip
			} else {
				v.Tip = ""
			}
			continue
		}

		// 判断还需要完成奖几个任务
		if len(kindMap[v.Type]) == 0 {
			v.Tip = ""
			continue
		}
		unfinished := 0
		for _, v := range kindMap[v.Type] {
			if v == nil {
				continue
			}
			if v.CompleteSate == newcomer.TaskIncomplete {
				unfinished++
			}
		}
		v.Tip = fmt.Sprintf("再完成%d个任务就能领取了呢", unfinished)
	}
}

// addTaskGroupTip get taskGroup tip
func (s *Service) addTaskGroupTip(tr []*newcomer.TaskRewardGroup) {
	for _, v := range tr {
		if v == nil {
			continue
		}
		if tip, ok := newcomer.TaskGroupTipMap[v.RewardState][v.GroupID]; ok {
			v.Tip = tip
		} else {
			v.Tip = ""
		}
	}
}
