package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"go-common/app/job/main/aegis/model"
	moniMdl "go-common/app/job/main/aegis/model/monitor"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

var (
	_taskTable = "task"
)

func (s *Service) taskconsumeproc() {
	defer func() {
		log.Warn("taskconsumeproc exited.")
		s.wg.Done()
	}()

	var (
		binLogMsgs = s.binLogDataBus.Messages()
	)

	for {
		select {
		case msg, ok := <-binLogMsgs:
			if !ok {
				log.Warn("binLogDataBus has been closed.")
				return
			}
			log.Info("binLogDataBus key(%s) offset(%d) message(%s)",
				msg.Key, msg.Offset, msg.Value)
			s.handleBinLog(msg)
		case rpi := <-s.chanReport:
			s.reportResource(context.Background(), rpi.BizID, rpi.FlowID, rpi.RID, rpi.UID)
		default:
			time.Sleep(time.Second)
		}
	}
}

func (s *Service) archiveConsumeProc() {
	defer func() {
		log.Warn("archiveConsumeProc exited.")
		s.wg.Done()
	}()
	var (
		msgs = s.archiveDataBus.Messages()
	)
	for {
		var (
			msg *databus.Message
			ok  bool
			err error
		)
		if msg, ok = <-msgs; !ok {
			log.Error("s.archiveDataBus.Messages() closed.")
			return
		}
		msg.Commit()

		m := &model.BinLog{}
		if err = json.Unmarshal(msg.Value, m); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
			continue
		}
		if m.Table == "archive" {
			s.handleArchiveBinlog(m)
		} else if m.Table == "archive_video" {
			s.handleVideoBinlog(m)
		}
	}
}

func (s *Service) handleArchiveBinlog(m *model.BinLog) {
	var (
		err error
	)
	na := &moniMdl.BinlogArchive{}
	oa := &moniMdl.BinlogArchive{}
	if err = json.Unmarshal(m.New, na); err != nil {
		log.Error("json.Unmarshal(%s,%+v) error(%v)", m.New, na, err)
		return
	}
	if err = json.Unmarshal(m.New, oa); err != nil {
		log.Error("json.Unmarshal(%s,%+v) error(%v)", m.New, oa, err)
		return
	}
	s.monitorArchive(m.Action, na, oa)

}

func (s *Service) handleVideoBinlog(m *model.BinLog) {
	var (
		err error
	)
	nv := &moniMdl.BinlogVideo{}
	ov := &moniMdl.BinlogVideo{}
	if err = json.Unmarshal(m.New, nv); err != nil {
		log.Error("json.Unmarshal(%s,%+v) error(%v)", m.New, nv, err)
		return
	}
	if err = json.Unmarshal(m.New, ov); err != nil {
		log.Error("json.Unmarshal(%s,%+v) error(%v)", m.New, ov, err)
		return
	}
	s.monitorVideo(m.Action, nv, ov)

}

func (s *Service) handleBinLog(msg *databus.Message) {
	defer msg.Commit()

	bmsg := new(model.BinLog)
	if err := json.Unmarshal(msg.Value, bmsg); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", string(msg.Value), err)
		return
	}
	if bmsg.Table == _taskTable {
		old := new(model.Task)
		new := new(model.Task)
		if err := json.Unmarshal(bmsg.New, new); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", string(bmsg.New), err)
			return
		}
		if bmsg.Action == "update" {
			if err := json.Unmarshal(bmsg.Old, old); err != nil {
				log.Error("json.Unmarshal(%s) error(%v)", string(bmsg.New), err)
				return
			}
		}
		s.handleBinLogMsg(context.Background(), bmsg.Action, old, new)
	}

	// use specify goroutine to merge messages
	log.Info("handleBinlog table:%s key:%s partition:%d offset:%d", bmsg.Table, msg.Key, msg.Partition, msg.Offset)
}

//各种状态简写
const (
	INT = model.TaskStateInit
	DSP = model.TaskStateDispatch
	DEY = model.TaskStateDelay
	SUB = model.TaskStateSubmit
	RSB = model.TaskStateRscSb
	CLO = model.TaskStateClosed
	LTD = model.LogTypeTaskDispatch
)

func (s *Service) handleRelease(c context.Context, old, new *model.Task) {
	s.dao.RemovePersonalTask(c, old.BusinessID, old.FlowID, old.UID, old.ID)
	s.dao.PushPublicTask(c, new)
	s.sendTaskLog(c, new, LTD, "release", new.UID, "")
	s.dao.IncresByField(c, old.BusinessID, old.FlowID, old.UID, model.Release, 1)
}

func (s *Service) handleDisptach(c context.Context, old, new *model.Task) {
	//这里不做缓存同步，顺序会乱
	s.sendTaskLog(c, new, LTD, "dispatch", new.UID, "")
	s.dao.IncresByField(c, new.BusinessID, new.FlowID, new.UID, model.Dispatch, 1)
}

func (s *Service) handleDelay(c context.Context, old, new *model.Task) {
	//这里不做缓存同步，顺序会乱
	s.sendTaskLog(c, new, LTD, "delay", new.UID, "")
	s.dao.IncresByField(c, new.BusinessID, new.FlowID, new.UID, model.Delay, 1)
}

/*
数据统计时，容易产生误差的几种数据
1. 任务被a领取，被b在资源列表提交
2. 任务被a延迟，被b在资源列表 或者 延迟列表提交
*/
func (s *Service) handleSubmit(c context.Context, old, new *model.Task) {
	switch old.State {
	case INT: // 未分配直接提交，资源列表里操作
		s.dao.RemovePublicTask(c, old.BusinessID, old.FlowID, old.ID)
	case DSP: // 领取后提交，也可能是资源列表操作
		s.dao.RemovePersonalTask(c, old.BusinessID, old.FlowID, old.UID, old.ID)
	case DEY: // 延迟列表提交，也可能是资源列表操作
		s.dao.RemoveDelayTask(c, old.BusinessID, old.FlowID, old.UID, old.ID)
	default: // 其他未知情况
		log.Error("handleSubmit UNEXPECTED old(%+v) new(%v)", old, new)
	}

	switch new.State {
	case SUB:
		s.sendTaskLog(c, new, LTD, "tasksubmit", new.UID, "")
	case RSB:
		s.sendTaskLog(c, new, LTD, "rscsubmit", new.UID, "")
	case CLO:
		s.sendTaskLog(c, new, LTD, "close", new.UID, "")
	}
	s.reportSubmit(c, old, new)
}

func (s *Service) handleCreate(c context.Context, new *model.Task) {
	s.dao.PushPublicTask(c, new)
	s.sendTaskLog(c, new, LTD, "create", 399, "aegis-job")
	s.reportTaskCreate(c, new)
}

func (s *Service) handleBinLogMsg(c context.Context, act string, old, new *model.Task) {
	log.Info("handleTaskBinlog act(%s) old(%+v) new(%+v)", act, old, new)
	s.dao.SetTask(c, new)
	if act == "insert" {
		s.handleCreate(c, new)
	}

	if act == "update" {
		switch {
		case old.State != new.State: //状态变更
			switch new.State {
			case INT: // 初始
				switch old.State {
				case DSP: //释放
					s.handleRelease(c, old, new)
				default: //其他情况
					s.dao.PushPublicTask(c, new)
					log.Error("handleTaskBinlog UNEXPECTED INT old(%+v) new(%+v)", old, new)
				}
			case DSP: // 领取
				switch old.State {
				case INT:
					s.handleDisptach(c, old, new)
				default:
					log.Error("handleTaskBinlog UNEXPECTED DSP old(%+v) new(%+v)", old, new)
				}
			case DEY: // 延迟
				switch old.State {
				case DSP:
					s.handleDelay(c, old, new)
				default:
					log.Error("handleTaskBinlog UNEXPECTED DEY old(%+v) new(%+v)", old, new)
				}
			case SUB, RSB, CLO: // 提交,关闭
				s.handleSubmit(c, old, new)
			}
		case old.AdminID != new.AdminID: //指派变更
		default:
			log.Info("其他变更 old(%+v)->new(%+v)", old, new)
		}
	}
}

func (s *Service) setAssign(c context.Context, task *model.Task) bool {
	log.Info("指派判断 setAssign(%+v)", task)
	auids := s.hitAssignUids(c, task)
	log.Info("指派判断 hitAssignUids(%v)", auids)
	if len(auids) == 0 {
		return false
	}

	log.Info("task(%d) 命中指派配置 (%v)", task.ID, auids)
	var huids []int64
	for auid, uids := range auids {
		task.AdminID = auid
		huids = s.hitActiveUids(c, task, uids)
		length := len(huids)
		if length != 0 {
			break
		}
	}

	log.Info("task(%d) 指派在线 (%v)", task.ID, huids)
	length := len(huids)
	if length == 0 {
		return false
	}

	if length == 1 {
		task.UID = huids[0]
	} else {
		// 随机数选一个
		task.UID = huids[rand.Intn(length)]
	}
	log.Info("task(%d) admin(%d) 指派成功 (%d)", task.ID, task.AdminID, task.UID)

	return true
}

func (s *Service) hitAssignUids(c context.Context, task *model.Task) (uids map[int64][]int64) {
	key := fmt.Sprintf("%d-%d", task.BusinessID, task.FlowID)
	uids = make(map[int64][]int64)
	if assignC, ok := s.assignConfig[key]; ok {
		for _, item := range assignC {
			log.Info("指派判断 task(%+v) item(%+v)", task, item)
			for _, mid := range item.Mids {
				if mid == task.MID {
					if aus, ok := uids[item.Admin]; ok {
						uids[item.Admin] = append(aus, item.Uids...)
					} else {
						uids[item.Admin] = item.Uids
					}
				}
			}
		}
	}
	return
}

func (s *Service) hitActiveUids(c context.Context, task *model.Task, uids []int64) (hitid []int64) {
	s.ccMux.RLock()
	defer s.ccMux.RUnlock()
	key := fmt.Sprintf("%d-%d", task.BusinessID, task.FlowID)
	if uidCache, ok := s.consumerCache[key]; ok {
		for _, uid := range uids {
			if _, ok := uidCache[uid]; ok {
				if on, _ := s.dao.IsConsumerOn(c, int(task.BusinessID), int(task.FlowID), uid); on {
					hitid = append(hitid, uid)
				}
			}
		}
	}
	return
}
