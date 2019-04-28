package service

import (
	"context"
	"errors"
	"time"

	"go-common/app/admin/main/aegis/model/common"
	taskmod "go-common/app/admin/main/aegis/model/task"
	"go-common/library/cache/redis"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
)

//ERROR
var (
	ErrTaskMiss = errors.New("task miss")
)

// NextTask 下一个任务
func (s *Service) NextTask(c context.Context, opt *taskmod.NextOptions) (tasks []*taskmod.Task, count int64, err error) {
	log.Info("task-NextTask opt(%+v)", opt)
	if count, err = s.countPersonalTask(c, &opt.BaseOptions, opt.NoCache); err != nil {
		return
	}

	if count < opt.DispatchCount {
		if count, err = s.syncSeize(c, opt); err != nil {
			return
		}
	}
	/* 去掉异步抢占
	else if count < opt.SeizeCount {
		s.asyncSeize(opt)
	}
	*/

	if count == 0 {
		return
	}

	return s.dispatch(c, opt)
}

// ListTasks 实时列表，停滞列表，延迟列表
func (s *Service) ListTasks(c context.Context, opt *taskmod.ListOptions) (tasks []*taskmod.Task, count int64, err error) {
	switch opt.State {
	case 1: // 实时任务,从redis取出,在数据库校验
		tasks, count, err = s.listUnseized(c, opt)
	case 2: // 停滞任务,组员的直接从redis取。组长的从数据库取id，redis取任务
		tasks, count, err = s.listMyTasks(c, "seized", opt)
	case 3: // 延迟任务,组员的直接从redis取。组长的从数据库取id，redis取任务
		tasks, count, err = s.listMyTasks(c, "delayd", opt)
	case 4: // 指派停滞任务,从数据库取id，redis取任务
		tasks, count, err = s.listMyTasks(c, "assignd", opt)
	default: //  所有未完成任务
	}

	if err != nil {
		tasks, count, err = s.mysql.ListTasks(c, opt)
	}

	opt.Total = int(count)
	return
}

// Task 直接读取某个任务
func (s *Service) Task(c context.Context, tid int64) (task *taskmod.Task, err error) {
	return s.mysql.TaskFromDB(c, tid)
}

// TxSubmitTask 提交任务
func (s *Service) TxSubmitTask(c context.Context, ormTx *gorm.DB, opt *common.BaseOptions, state int8) (ostate int8, otaskid, ouid int64, err error) {
	var (
		t    *taskmod.Task
		rows int64
	)

	// 根据rid,flowid检索最新的未完成任务
	if t, err = s.gorm.TaskByRID(c, opt.RID, opt.FlowID); err != nil || t == nil || t.ID == 0 {
		log.Warn("TaskByRID(%d,%d) miss(%v)", opt.RID, opt.FlowID, err)
		t, err = s.gorm.TaskByRID(c, opt.RID, 0)
	}
	// TODO 先默认一个资源只在一个flow下分发，解决目前存在flow状态与task状态不同步
	if err != nil || t == nil || t.ID == 0 {
		log.Warn("s.gorm.TaskByRID(%d,%d) miss(%v)", opt.RID, 0, err)
		err = nil
		return
	}

	ostate = t.State
	ouid = t.UID
	otaskid = t.ID
	var utime uint64
	if !t.Gtime.Time().IsZero() {
		utime = uint64(time.Since(t.Gtime.Time()).Seconds())
	}

	subopt := &taskmod.SubmitOptions{
		BaseOptions: *opt,
		TaskID:      t.ID,
		OldUID:      t.UID,
		Utime:       utime,
		OldState:    t.State,
	}

	// 1. 改数据库
	if rows, err = s.gorm.TxSubmit(ormTx, subopt, state); err != nil {
		return
	}
	if rows != 1 {
		err = ecode.NothingFound
		log.Error("Submit (%v) error(%v)", opt, err)
		return
	}
	return
}

func (s *Service) submitTaskCache(c context.Context, opt *common.BaseOptions, ostate int8, taskid, ouid int64) {
	log.Info("SubmitTaskCache opt(%+v) ostate(%d) taskid(%d) ouid(%d)", opt, ostate, taskid, ouid)
	optc := &common.BaseOptions{
		BusinessID: opt.BusinessID,
		FlowID:     opt.FlowID,
		UID:        ouid,
	}
	if ostate == taskmod.TaskStateDelay {
		s.redis.RemoveDelayTask(c, optc, taskid)
		return
	}
	s.redis.RemovePersonalTask(c, optc, taskid)
}

// Delay 延迟任务
func (s *Service) Delay(c context.Context, opt *taskmod.DelayOptions) (err error) {
	var (
		taskmod *taskmod.Task
		rows    int64
	)

	if taskmod, err = s.mysql.TaskFromDB(c, opt.TaskID); err != nil || taskmod == nil {
		return
	}
	if !s.checkDelayOption(c, opt, taskmod) {
		log.Error("checkDelayOption error opt(%+v) taskmod(%+v)", opt, taskmod)
		return ecode.AegisTaskFinish
	}

	if rows, err = s.mysql.Delay(c, opt); err != nil {
		return
	}
	if rows != 1 {
		err = ecode.AegisTaskFinish
		log.Error("Submit (%v) error(%v)", opt, err)
		return
	}
	if err = s.redis.RemovePersonalTask(c, &opt.BaseOptions, opt.TaskID); err != nil {
		return
	}
	s.redis.PushDelayTask(c, &opt.BaseOptions, opt.TaskID)

	return
}

// Release 释放任务
func (s *Service) Release(c context.Context, opt *common.BaseOptions, delay bool) (rows int64, err error) {
	if rows, err = s.mysql.Release(c, opt, delay); err != nil {
		return
	}
	//err = s.redis.Release(c, opt, delay)
	return
}

// MaxWeight 当前最高权重
func (s *Service) MaxWeight(c context.Context, opt *common.BaseOptions) (max int64, err error) {
	return s.gorm.MaxWeight(c, opt.BusinessID, opt.FlowID)
}

// UnDoStat undo stat
func (s *Service) UnDoStat(c context.Context, opt *common.BaseOptions) (stat *taskmod.UnDOStat, err error) {
	return s.gorm.UndoStat(c, opt.BusinessID, opt.FlowID, opt.UID)
}

// TaskStat task stat
func (s *Service) TaskStat(c context.Context, opt *common.BaseOptions) (stat *taskmod.Stat, err error) {
	return s.gorm.TaskStat(c, opt.BusinessID, opt.FlowID, opt.UID)
}

func (s *Service) countPersonalTask(c context.Context, opt *common.BaseOptions, nocache bool) (count int64, err error) {
	log.Info("task-countPersonalTask opt(%+v) nocache(%v)", opt, nocache)
	defer func() { log.Info("task-countPersonalTask count(%d) err(%v)", count, err) }()

	if nocache {
		return s.mysql.CountPersonal(c, opt)
	}

	if count, err = s.redis.CountPersonalTask(c, opt); err != nil {
		// redis 挂了
		if count, err = s.mysql.CountPersonal(c, opt); err != nil {
			return
		}
	}
	return
}

func (s *Service) syncSeize(c context.Context, opt *taskmod.NextOptions) (count int64, err error) {
	return s.seize(c, opt)
}

func (s *Service) seize(c context.Context, opt *taskmod.NextOptions) (count int64, err error) {
	log.Info("task-seize opt(%+v)", opt)
	defer func() { log.Info("task-seize count(%d) err(%v)", count, err) }()

	var (
		hitids, missids []int64
		others          map[int64]int64
	)

	// TODO: 抢占任务要根据用户是否在线，处理任务指派
	if opt.NoCache {
		hitids, err = s.mysql.QueryForSeize(c, opt.BusinessID, opt.FlowID, opt.UID, opt.SeizeCount)
	} else {
		hitids, missids, others, err = s.redis.SeizeTask(c, opt.BusinessID, opt.FlowID, opt.UID, opt.SeizeCount)
		if err != nil {
			hitids, err = s.mysql.QueryForSeize(c, opt.BusinessID, opt.FlowID, opt.UID, opt.SeizeCount)
		}
	}
	if err != nil {
		return
	}

	log.Info("seize uid(%d) hitids(%v), missids(%v), others(%v)", opt.UID, hitids, missids, others)

	if !opt.NoCache && len(missids) > 0 {
		log.Error("seize uid(%d) missids(%v)", opt.UID, missids)
		for _, id := range missids {
			if err = s.syncTask(c, id); err != nil {
				s.redis.RemovePublicTask(c, &opt.BaseOptions, id)
			}
		}
	}

	if len(hitids) > 0 {
		log.Info("seize uid(%d) hitids(%v)", opt.UID, hitids)
		mhits := make(map[int64]int64)
		for _, id := range hitids {
			mhits[id] = opt.UID
		}
		if count, err = s.mysql.Seize(c, mhits); err != nil || count == 0 {
			return
		}
		return
	}
	return
}

func (s *Service) dispatch(c context.Context, opt *taskmod.NextOptions) (tasks []*taskmod.Task, count int64, err error) {
	log.Info("task-dispatch opt(%+v)", opt)
	defer func() { log.Info("task-dispatch tasks(%+v) count(%d) err(%v)", tasks, count, err) }()
	listopt := &taskmod.ListOptions{
		BaseOptions: opt.BaseOptions,
		Pager: common.Pager{
			Pn: 1,
			Ps: int(opt.DispatchCount),
		}}
	tasks, count, err = s.calibur(c, listopt, s.redis.RangePersonalTask, s.mysql.DispatchByID, s.redis.RemovePersonalTask)
	if err != nil {
		tasks, count, err = s.mysql.DBDispatch(c, opt)
	}
	return
}

func (s *Service) syncTask(c context.Context, taskID int64) (err error) {
	var task *taskmod.Task

	if task, err = s.mysql.TaskFromDB(c, taskID); err != nil || task == nil {
		return ErrTaskMiss
	}

	var option = &common.BaseOptions{
		BusinessID: task.BusinessID,
		FlowID:     task.FlowID,
		UID:        task.UID,
	}

	s.redis.SetTask(c, task)
	switch task.State {
	case taskmod.TaskStateInit:
		s.redis.PushPublicTask(c, task)
	case taskmod.TaskStateDispatch:
		s.redis.RemovePublicTask(c, option, task.ID)
		s.redis.PushPersonalTask(c, option, task.ID)
	case taskmod.TaskStateDelay:
		s.redis.RemovePublicTask(c, option, task.ID)
		s.redis.PushDelayTask(c, option, task.ID)
	default:
		s.redis.RemovePublicTask(c, option, task.ID)
	}

	return
}

func (s *Service) listUnseized(c context.Context, opt *taskmod.ListOptions) (tasks []*taskmod.Task, count int64, err error) {
	return s.calibur(c, opt, s.redis.RangePublicTask, s.mysql.ListCheckUnSeized, s.redis.RemovePublicTask)
}

func (s *Service) listMyTasks(c context.Context, ltp string, opt *taskmod.ListOptions) (tasks []*taskmod.Task, count int64, err error) {
	if !opt.BisLeader {
		if ltp == "delayd" {
			return s.calibur(c, opt, s.redis.RangeDealyTask, s.mysql.ListCheckDelay, s.redis.RemoveDelayTask)
		}
		if ltp == "seized" {
			return s.calibur(c, opt, s.redis.RangePersonalTask, s.mysql.ListCheckSeized, s.redis.RemovePersonalTask)
		}
	}
	if opt.BisLeader {
		opt.UID = 0
	}
	var ids []int64
	switch ltp {
	case "delayd":
		ids, count, err = s.gorm.TaskListDelayd(c, opt)
	case "seized":
		ids, count, err = s.gorm.TaskListSeized(c, opt)
	case "assignd":
		ids, count, err = s.gorm.TaskListAssignd(c, opt)
	}
	if err != nil || len(ids) == 0 {
		return
	}
	if tasks, err = s.redis.GetTask(c, ids); err != nil {
		err = redis.ErrNil
	}
	return
}

func (s *Service) calibur(c context.Context, opt *taskmod.ListOptions, rfunc taskmod.RangeFunc, lfunc taskmod.ListFuncDB, remove taskmod.RemoveFunc) (taskmods []*taskmod.Task, count int64, err error) {
	var (
		hitids, missids []int64
		missmap         map[int64]struct{}
		mtaskmods       map[int64]*taskmod.Task
	)

	mtaskmods, count, hitids, missids, err = rfunc(c, opt)
	log.Info("calibur(%+v) rfunc count(%d) hitids(%v) missids(%v)", opt, count, hitids, missids)
	if err != nil {
		return
	}

	if len(missids) > 0 {
		for _, id := range missids {
			if err = s.syncTask(c, id); err != nil {
				log.Error("syncTask error(%v)", err)
				remove(c, &opt.BaseOptions, id)
			}
		}
	}
	if len(hitids) > 0 {
		if missmap, err = lfunc(c, mtaskmods, hitids, opt.UID); err != nil {
			log.Error("calibur lfunc error(%v)", err)
			return
		}
		if len(missmap) > 0 {
			log.Info("calibur personal任务移除%v", missmap)
			for id := range missmap {
				remove(c, &opt.BaseOptions, id)
			}
		}
	}

	for _, id := range hitids {
		if _, ok := missmap[id]; ok && opt.Action != "release" {
			delete(mtaskmods, id)
		} else {
			taskmods = append(taskmods, mtaskmods[id])
		}
	}
	return
}

/*
func (s *Service) checkSubmitOption(c context.Context, opt *taskmod.SubmitOptions, task *taskmod.Task) bool {
	opt.OldState = task.State
	opt.OldUID = task.UID
	// 1. 组员只能处理自己的延迟任务
	if task.State == taskmod.TaskStateDelay {
		if opt.BisLeader {
			return true
		}
		if task.UID != opt.UID {
			return false
		}
	}
	if task.State == taskmod.TaskStateDispatch && opt.UID == task.UID {
		opt.Utime = uint64(time.Since(task.Gtime.Time()).Seconds())
		return true
	}

	return false
}
*/

func (s *Service) checkDelayOption(c context.Context, opt *taskmod.DelayOptions, task *taskmod.Task) bool {
	if task.State == taskmod.TaskStateDispatch && task.UID == opt.UID {
		return true
	}
	return false
}

func (s *Service) syncUpCache(c context.Context) (err error) {
	if s.Debug() == "local" {
		return
	}
	upGroup := make(map[int64]*common.Group)
	upgs, err := s.rpc.UpGroups(c)
	if err != nil || upgs == nil {
		return
	}

	for gid, upg := range upgs {
		if _, ok := upGroup[gid]; !ok {
			upGroup[gid] = &common.Group{
				ID:        gid,
				Name:      upg.Name,
				Note:      upg.Note,
				Tag:       upg.Tag,
				FontColor: upg.FontColor,
				BgColor:   upg.BgColor,
			}
			log.Info("groupCache upg(%+v) upGroup(%+v)", upg, upGroup[gid])
		}
	}
	s.groupCache = upGroup
	log.Info("groupCache(%+v)", s.groupCache)

	return
}
