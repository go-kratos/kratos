package service

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"go-common/app/job/main/dm2/model"
	"go-common/app/job/main/dm2/model/oplog"
	"go-common/library/log"
)

func (s *Service) taskResProc() {
	var (
		c     = context.Background()
		tasks []*model.TaskInfo
		err   error
	)
	ticker := time.NewTicker(time.Duration(s.conf.TaskConf.ResInterval))
	defer ticker.Stop()
	for range ticker.C {
		if tasks, err = s.dao.TaskInfos(c, model.TaskStateSearch); err != nil {
			log.Error("s.dao.TaskInfos error(%v)", err)
			continue
		}
		for _, task := range tasks {
			count, url, state, err := s.dao.TaskSearchRes(c, task)
			if err != nil {
				log.Error("s.dao.TaskSearchRes(%+v) error(%v)", task, err)
				continue
			}
			if state == model.TaskSearchFail {
				task.State = model.TaskStateFail
			} else if state == model.TaskSearchSuc {
				task.Result = url
				task.Count = count
				if task.Sub > 0 {
					task.State = model.TaskStateWait
				} else {
					task.State = model.TaskStateSuc
				}
			}
			s.dao.UpdateTask(c, task)
		}
	}
}

func (s *Service) taskDelProc() {
	var (
		c   = context.Background()
		err error
	)
	ticker := time.NewTicker(time.Duration(s.conf.TaskConf.DelInterval))
	defer ticker.Stop()
	for range ticker.C {
		if err = s.taskSchedule(c); err != nil {
			log.Error("taskDelProc error(%v)", err)
			continue
		}
	}
}

func (s *Service) taskSchedule(c context.Context) (err error) {
	var (
		ok                               bool
		now                              = time.Now()
		expire                           = now.Add(time.Duration(s.conf.TaskConf.DelInterval))
		expireStr                        = expire.Format(time.RFC3339)
		oldExpireStr, oldExpireGetSetStr string
		oldExpire                        time.Time
	)
	if ok, err = s.dao.SetnxTaskJob(c, expireStr); err != nil {
		return
	}
	// redis中不存在
	if ok {
		if err = s.taskDelJob(c); err != nil {
			s.dao.DelTaskJob(c)
			log.Error("taskDelJob,error(%v)", err)
			return
		}
		return
	}
	// redis中已经存在
	// 判断是否过期了
	if oldExpireStr, err = s.dao.GetTaskJob(c); err != nil {
		return
	}
	if oldExpire, err = time.Parse(time.RFC3339, oldExpireStr); err != nil {
		return
	}
	if oldExpire.Sub(now) > 0 {
		return
	}
	if oldExpireGetSetStr, err = s.dao.GetSetTaskJob(c, expireStr); err != nil {
		return
	}
	if oldExpireGetSetStr != oldExpireStr {
		return
	}
	if err = s.taskDelJob(c); err != nil {
		s.dao.DelTaskJob(c)
		log.Error("taskDelJob,error(%v)", err)
		return
	}
	return
}

// TODO: operation_time && operation_rate
func (s *Service) taskDelJob(c context.Context) (err error) {
	var (
		task *model.TaskInfo
	)
	if task, err = s.dao.OneTask(c); err != nil || task == nil {
		return
	}
	task.State = model.TaskStateDelDM
	s.dao.UpdateTask(c, task)
	var delCount int64
	if delCount, task.LastIndex, task.State, err = s.taskDelDM(c, task); err != nil {
		return
	}
	if task.State == model.TaskStateDelDM {
		task.State = model.TaskStateSuc
	}
	if _, err = s.dao.UptSubTask(c, task.ID, delCount, time.Now()); err != nil {
		return
	}
	_, err = s.dao.UpdateTask(c, task)
	return
}

func (s *Service) taskDelDM(c context.Context, eTask *model.TaskInfo) (delCount int64, lastIndex, state int32, err error) {
	taskDelNum := s.conf.TaskConf.DelNum
	taskResFieldLen := s.conf.TaskConf.ResFieldLen
	res, err := http.Get(eTask.Result)
	if err != nil {
		log.Error("s.taskDelDM.HttpGet(%s) error(%v)", eTask.Result, err)
		return
	}
	resp, err := ioutil.ReadAll(res.Body)
	if err != nil {
		res.Body.Close()
		log.Error("s.taskDelDM.ioutilRead error(%v)", err)
		return
	}
	res.Body.Close()
	lines := bytes.Split(resp, []byte("\n"))
	total := len(lines)
	n := (total-1)/taskDelNum + 1
	for i := int(eTask.LastIndex); i < n; i++ {
		var (
			task    *model.TaskInfo
			subTask *model.SubTask
		)
		start := i * taskDelNum
		end := (i + 1) * taskDelNum
		if end > total {
			end = total
		}
		OidDMid := make(map[int64][]int64)
		for _, line := range lines[start:end] {
			var dmid, oid int64
			fields := bytes.Split(line, []byte("\001"))
			if len(fields) < taskResFieldLen {
				log.Error("fields lenth too small:%d", len(fields))
				continue
			}
			if dmid, err = strconv.ParseInt(string(fields[0]), 10, 64); err != nil {
				log.Error("ParseInt(%s) error(%v)", string(fields[0]), err)
				continue
			}
			if oid, err = strconv.ParseInt(string(fields[1]), 10, 64); err != nil {
				log.Error("ParseInt(%s) error(%v)", string(fields[1]), err)
				continue
			}
			OidDMid[oid] = append(OidDMid[oid], dmid)
		}
		for oid, dmids := range OidDMid {
			var affected int64
			if affected, err = s.dao.DelDMs(c, oid, dmids, model.StateTaskDel); err != nil {
				log.Error("dm task(id:%d) del dm(oid:%d,dmids:%v) error(%v)", eTask.ID, oid, dmids, err)
				continue
			}
			if affected > 0 {
				s.OpLog(c, oid, 0, time.Now().Unix(), int(model.SubTypeVideo), dmids, "status", "", strconv.FormatInt(int64(model.StateTaskDel), 10), "弹幕任务删除", oplog.SourceManager, oplog.OperatorSystem)
				delCount += affected
				if _, err = s.dao.UptSubjectCount(c, model.SubTypeVideo, oid, affected); err != nil {
					log.Error("dm task update count(oid:%d,affected:%d) error(%v)", oid, affected, err)
				}
			}
			time.Sleep(50 * time.Millisecond)
		}
		if len(OidDMid) > 0 {
			log.Warn("dm task(id:%d) del dm(oid,dmids:%+v)", eTask.ID, OidDMid)
		}
		lastIndex = int32(i + 1)
		task, err = s.dao.OneTask(c)
		if err == nil && task != nil && task.ID != eTask.ID && task.Priority > eTask.Priority {
			state = model.TaskStateWait
			return
		}
		if eTask, err = s.dao.TaskInfoByID(c, eTask.ID); err != nil || task == nil {
			continue
		}
		state = eTask.State
		if state != model.TaskStateDelDM {
			return
		}
		if subTask, err = s.dao.SubTask(c, eTask.ID); err != nil || subTask == nil {
			continue
		}
		tCount := subTask.Tcount + delCount
		if tCount >= s.conf.TaskConf.DelLimit && subTask.Tcount < s.conf.TaskConf.DelLimit {
			log.Warn("task(id:%d) del dm reach limit(count:%d)", eTask.ID, tCount)
			s.sendWechatWorkMsg(c, eTask, tCount)
			state = model.TaskStatePause
			return
		}
	}
	return
}

func (s *Service) sendWechatWorkMsg(c context.Context, task *model.TaskInfo, count int64) (err error) {
	content := fmt.Sprintf(model.TaskNoticeContent, task.ID, task.Title, count)
	users := s.conf.TaskConf.MsgCC
	users = append(users, task.Creator, task.Reviewer)
	return s.dao.SendWechatWorkMsg(c, content, model.TaskNoticeTitle, users)
}

// OpLog put a new infoc format operation log into the channel
func (s *Service) OpLog(c context.Context, cid, operator, OperationTime int64, typ int, dmids []int64, subject, originVal, currentVal, remark string, source oplog.Source, operatorType oplog.OperatorType) (err error) {
	infoLog := new(oplog.Infoc)
	infoLog.Oid = cid
	infoLog.Type = typ
	infoLog.DMIds = dmids
	infoLog.Subject = subject
	infoLog.OriginVal = originVal
	infoLog.CurrentVal = currentVal
	infoLog.OperationTime = strconv.FormatInt(OperationTime, 10)
	infoLog.Source = source
	infoLog.OperatorType = operatorType
	infoLog.Operator = operator
	infoLog.Remark = remark
	select {
	case s.opsLogCh <- infoLog:
	default:
		err = fmt.Errorf("opsLogCh full")
		log.Error("opsLogCh full (%v)", infoLog)
	}
	return
}

func (s *Service) oplogproc() {
	for opLog := range s.opsLogCh {
		if len(opLog.Subject) == 0 || len(opLog.CurrentVal) == 0 || opLog.Source <= 0 ||
			opLog.Operator < 0 || opLog.OperatorType <= 0 {
			log.Warn("oplogproc() it is an illegal log, warn(%v, %v, %v)", opLog.Subject, opLog.Subject, opLog.CurrentVal)
			continue
		} else {
			for _, dmid := range opLog.DMIds {
				if dmid > 0 {
					s.dmOperationLogSvc.Info(opLog.Subject, strconv.FormatInt(opLog.Oid, 10), strconv.Itoa(opLog.Type),
						strconv.FormatInt(dmid, 10), opLog.Source.String(), opLog.OriginVal,
						opLog.CurrentVal, strconv.FormatInt(opLog.Operator, 10), opLog.OperatorType.String(),
						opLog.OperationTime, opLog.Remark)
				} else {
					log.Warn("oplogproc() it is an illegal log, for dmid value, warn(%d, %+v)", dmid, opLog)
				}
			}
		}
	}
}
