package service

import (
	"context"
	"go-common/app/job/main/videoup-report/model/archive"
	"go-common/library/log"
	"sort"
	"time"
)

func (s *Service) loadTask() {
	var (
		err   error
		took  *archive.TaskTook
		tooks []*archive.TaskTook
		tasks []*archive.Task
	)
	s.taskCache.Lock()
	defer s.taskCache.Unlock()
	if len(s.taskCache.Took) == 0 && len(s.taskCache.Task) == 0 {
		if took, err = s.arc.TaskTookByHalfHour(context.TODO()); err != nil {
			log.Error("s.arc.TaskTookByHalfHour error(%v)", err)
			return
		}
		if took != nil {
			if tooks, err = s.arc.TaskTooks(context.TODO(), took.Ctime); err != nil {
				log.Error("s.arc.TaskTooks(%v) error(%v)", took.Ctime, err)
				return
			}
			s.taskCache.Took = tooks
		}
		if tasks, err = s.arc.TaskByUntreated(context.TODO()); err != nil {
			log.Error("s.arc.TaskByUntreated() error(%v)", err)
			return
		}
	} else {
		var tasksOrig, tasksDone []*archive.Task
		if tasksOrig, err = s.arc.TaskByMtime(context.TODO(), s.taskCache.Mtime.Add(-time.Minute*1)); err != nil {
			log.Error("s.arc.TaskByMtime(%v) error(%v)", s.taskCache.Mtime, err)
			return
		}
		if tasksDone, err = s.arc.TaskDoneByMtime(context.TODO(), s.taskCache.Mtime.Add(-time.Minute*1)); err != nil {
			log.Error("s.arc.TaskDoneByMtime(%v) error(%v)", s.taskCache.Mtime, err)
			return
		}
		tasks = make([]*archive.Task, len(tasksOrig)+len(tasksDone))
		copy(tasks, tasksOrig)
		copy(tasks[len(tasksOrig):], tasksDone)
	}
	for _, task := range tasks {
		_, ok := s.taskCache.Task[task.ID]
		if ok && (task.State != archive.TaskStateUnclaimed && task.State != archive.TaskStateUntreated) {
			delete(s.taskCache.Task, task.ID)
		} else if task.State == archive.TaskStateUnclaimed || task.State == archive.TaskStateUntreated {
			s.taskCache.Task[task.ID] = task
		}
	}
}

func (s *Service) loadTaskTookSort() {
	var (
		took         int
		tooks        []int
		taskMinCtime *archive.Task
	)
	s.taskCache.Lock()
	defer s.taskCache.Unlock()
	for _, task := range s.taskCache.Task {
		if (s.taskCache.Mtime == time.Time{} || s.taskCache.Mtime.Unix() < task.Mtime.Unix()) {
			s.taskCache.Mtime = task.Mtime
		}
		if taskMinCtime == nil || taskMinCtime.Ctime.Unix() > task.Ctime.Unix() {
			taskMinCtime = task
		}
		took = int(time.Now().Unix() - task.Ctime.Unix())
		tooks = append(tooks, took)
	}
	if len(tooks) == 0 {
		return
	}
	sort.Ints(tooks)
	s.taskCache.Sort = tooks
	log.Info("s.loadTaskTookSort() 本轮统计: 耗时最久id(%d) ctime(%v)", taskMinCtime.ID, taskMinCtime.Ctime)
}

func (s *Service) hdlTaskTook() (lastID int64, err error) {
	s.taskCache.Lock()
	defer s.taskCache.Unlock()
	var (
		spacing         float32
		m50             float32
		m50Index        float32
		m50IndexPoint   float32
		m50Value        int
		m60             float32
		m60Index        float32
		m60IndexPoint   float32
		m60Value        int
		m80             float32
		m80Index        float32
		m80IndexPoint   float32
		m80Value        int
		m90             float32
		m90Index        float32
		m90IndexPoint   float32
		m90Value        int
		took            *archive.TaskTook
		taskTookSortLen = len(s.taskCache.Sort)
	)
	if taskTookSortLen > 1 {
		spacing = float32(taskTookSortLen-1) / 10
		m50Index = 1 + spacing*5
		m50IndexPoint = m50Index - float32(int(m50Index))
		m50Value = s.taskCache.Sort[int(m50Index)-1]
		m50 = float32(s.taskCache.Sort[int(m50Index)]-m50Value)*m50IndexPoint + float32(m50Value)
		m60Index = 1 + spacing*6
		m60IndexPoint = m60Index - float32(int(m60Index))
		m60Value = s.taskCache.Sort[int(m60Index)-1]
		m60 = float32(s.taskCache.Sort[int(m60Index)]-m60Value)*m60IndexPoint + float32(m60Value)
		m80Index = 1 + spacing*8
		m80IndexPoint = m80Index - float32(int(m80Index))
		m80Value = s.taskCache.Sort[int(m80Index)-1]
		m80 = float32(s.taskCache.Sort[int(m80Index)]-m80Value)*m80IndexPoint + float32(m80Value)
		m90Index = 1 + spacing*9
		m90IndexPoint = m90Index - float32(int(m90Index))
		m90Value = s.taskCache.Sort[int(m90Index)-1]
		m90 = float32(s.taskCache.Sort[int(m90Index)]-m90Value)*m90IndexPoint + float32(m90Value)
		took = &archive.TaskTook{}
		took.M50 = int(m50 + 0.5)
		took.M60 = int(m60 + 0.5)
		took.M80 = int(m80 + 0.5)
		took.M90 = int(m90 + 0.5)
		took.TypeID = archive.TookTypeMinute
		took.Ctime = time.Now()
		took.Mtime = took.Ctime
		s.taskCache.Took = append(s.taskCache.Took, took)
		lastID, err = s.arc.AddTaskTook(context.TODO(), took)
	}
	return
}

func (s *Service) hdlTaskTookByHourHalf() (lastID int64, err error) {
	s.taskCache.Lock()
	defer s.taskCache.Unlock()
	var (
		m50     int
		m60     int
		m80     int
		m90     int
		took    *archive.TaskTook
		tookLen = len(s.taskCache.Took)
	)
	for _, v := range s.taskCache.Took {
		m50 += v.M50
		m60 += v.M60
		m80 += v.M80
		m90 += v.M90
	}
	if tookLen >= 30 {
		took = &archive.TaskTook{}
		m50 /= int(float32(tookLen) + 0.5)
		m60 /= int(float32(tookLen) + 0.5)
		m80 /= int(float32(tookLen) + 0.5)
		m90 /= int(float32(tookLen) + 0.5)
		took.M50 = m50
		took.M60 = m60
		took.M80 = m80
		took.M90 = m90
		took.Ctime = time.Now()
		took.Mtime = took.Ctime
		took.TypeID = archive.TookTypeHalfHour
		lastID, err = s.arc.AddTaskTook(context.TODO(), took)
		s.taskCache.Took = nil
	}
	return
}

// TaskTooksByHalfHour get task books by ctime
func (s *Service) TaskTooksByHalfHour(c context.Context, stime, etime time.Time) (tooks []*archive.TaskTook, err error) {
	if tooks, err = s.arc.TaskTooksByHalfHour(c, stime, etime); err != nil {
		log.Error("s.arc.TaskTooksByHalfHour(%v,%v)", stime, etime)
		return
	}
	return
}
