package service

import (
	"context"
	"time"

	"go-common/app/admin/main/up/util"
	"go-common/app/job/main/up/conf"
	"go-common/app/job/main/up/dao/upcrm"
	"go-common/app/job/main/up/model/signmodel"
	"go-common/app/job/main/up/model/upcrmmodel"
	v1 "go-common/app/service/main/archive/api"
	"go-common/library/log"
	xtime "go-common/library/time"

	"github.com/jinzhu/gorm"
)

//CheckTaskJob check task job
func (s *Service) CheckTaskJob(tm time.Time) {
	// 今天计算昨天的数据
	var yesterday = tm.AddDate(0, 0, -1)
	log.Info("start to run CheckTaskJob, date=%s, yesterday=%s", tm, yesterday)
	s.CheckTaskFinish(yesterday)
	log.Info("finish run CheckTaskJob, date=%s", tm)

}

//Archive data
type Archive struct {
	ID int64 `gorm:"column:id"`
}

//CheckTaskFinish check task finish, calculate datas in (-,date]
func (s *Service) CheckTaskFinish(date time.Time) {
	// 1.查找所有有效的合同id, begin_date <= date && end _date >= date
	// 2.找到所有合同id对应的任务id,
	// 3.根据任务类型，日、周、月、累计，计算任务周期[a,b)
	// 4.计算完成数量
	var crmdb = s.crmdb.GetDb()
	var dateStr = date.Format(upcrmmodel.TimeFmtDate)
	var offset = 0
	var limit = 200
	var actualSize = limit
	log.Info("start to check task state")
	archiveDb, err := gorm.Open("mysql", conf.Conf.ArchiveOrm.DSN)
	s.crmdb.StartTask(upcrmmodel.TaskTypeSignTaskCalculate, date)
	archiveDb.LogMode(true)
	if err != nil {
		log.Error("connect archive db fail")
		return
	}
	defer archiveDb.Close()
	defer func() {
		if err == nil {
			s.crmdb.FinishTask(upcrmmodel.TaskTypeSignTaskCalculate, date, upcrmmodel.TaskStateFinish)
		} else {
			s.crmdb.FinishTask(upcrmmodel.TaskTypeSignTaskCalculate, date, upcrmmodel.TaskStateError)
		}
	}()
	var taskTotalCount = 0
	for actualSize == limit {
		var signUps []*signmodel.SignUp
		var signUpMap = make(map[uint32]*signmodel.SignUp)
		// 1
		err = crmdb.Table(signmodel.TableNameSignUp).
			Select("id, begin_date, end_date").
			Where("begin_date <= ? and end_date >= ?", dateStr, dateStr).
			Offset(offset).
			Limit(limit).
			Find(&signUps).Error

		offset += limit
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Error("err get signs, err=%+v", err)
			return
		}
		actualSize = len(signUps)
		var signIDs []uint32
		for _, v := range signUps {
			signUpMap[v.ID] = v
			signIDs = append(signIDs, v.ID)
		}
		// 2
		var taskList []*signmodel.SignTask
		err = crmdb.Where("sign_id in (?) and state != ? and generate_date<?", signIDs, signmodel.SignTaskStateDelete, dateStr).
			Find(&taskList).Error

		if err != nil {
			log.Error("err get tasks, err=%+v", err)
			return
		}

		// 3
		for _, task := range taskList {
			taskTotalCount++
			if task.Mid == 0 {
				log.Error("task's mid is zero, please check! task id=%d", task.ID)
				continue
			}
			var signInfo, ok = signUpMap[task.SignID]
			if !ok {
				log.Error("sign not found, err=%v", err)
				continue
			}
			// 4 计算数量
			err = s.checkSingleTask(task, signInfo, date, archiveDb)
			if err != nil {
				log.Error("check task err, task_id=%d, err=%v", task.ID, err)
				continue
			}
		}
		log.Info("finish to check task state, task total num=%d", taskTotalCount)
	}

}

// get task history, if not exist, then will create it
func (s *Service) getOrCreateTaskHistory(task *signmodel.SignTask, generateDate time.Time) (res *signmodel.SignTaskHistory, err error) {
	var crmdb = s.crmdb.GetDb()
	res = new(signmodel.SignTaskHistory)
	err = crmdb.Select("*").Where("task_template_id=? and generate_date=?", task.ID, generateDate).
		Find(&res).Error

	// 创建一条，如果没找到的话
	if err == gorm.ErrRecordNotFound {
		res = &signmodel.SignTaskHistory{
			Mid:            task.Mid,
			SignID:         task.SignID,
			TaskTemplateID: task.ID,
			TaskType:       task.TaskType,
			TaskCondition:  task.TaskCondition,
			Attribute:      task.Attribute,
			GenerateDate:   xtime.Time(generateDate.Unix()),
			State:          signmodel.SignTaskStateRunning,
		}

		err = crmdb.Save(&res).Error
		if err != nil {
			log.Error("create task history fail, err=%v, task=%v", err, task)
			return
		}
	}

	return
}

// check task state
func (s *Service) checkSingleTask(task *signmodel.SignTask, signInfo *signmodel.SignUp, date time.Time, archiveDb *gorm.DB) (err error) {
	var taskBegin, taskEnd time.Time
	if task.TaskType == signmodel.TaskTypeAccumulate {
		taskBegin = signInfo.BeginDate.Time()
		taskEnd = signInfo.EndDate.Time()
	} else {
		taskBegin, taskEnd = upcrm.GetTaskDuration(date, task.TaskType)
	}

	if task.Mid == 0 {
		log.Error("task's mid is zero, please check! task id=%d", task.ID)
		return
	}

	// get task history
	// 如果是累计任务，这里的taskBegin要设置为0
	var tBegin = taskBegin
	if task.TaskType == signmodel.TaskTypeAccumulate {
		tBegin = time.Time{}
	}
	taskHistory, err := s.getOrCreateTaskHistory(task, tBegin)
	if err != nil {
		log.Error("get task history fail, task=%+v, err=%+v", task, err)
		return
	}

	var dateStr = date.Format(upcrmmodel.TimeFmtDate)
	switch {
	default:
		var crmdb = s.crmdb.GetDb()

		// 4.去稿件库中查找对应的稿件数量
		var archiveCount = 0
		var archiveList []*Archive
		err = archiveDb.Table("archive").
			Where("mid = ? and ctime>= ? and ctime <? and (state >= 0 or state = -6)",
				task.Mid, taskBegin, taskEnd).
			Select("id").
			Find(&archiveList).
			Error

		if err != nil && err != gorm.ErrRecordNotFound {
			log.Error("check archive count fail, taskid=%+v err=%+v", task, err)
			break
		}

		var finalResult []int64
		// 5.任务完成度统计时，
		// 	若录入为不包含商单，
		// 		则任务完成数=新增稿件数-减去（绿洲/商单报备）稿件数+请假任务数；
		// 	若录入为包含商单，
		// 		则任务完成数=新增稿件数+请假任务数；
		// 如果没有archive，就直接返回
		if len(archiveList) != 0 {
			// 需要判断商单
			if task.IsAttrSet(signmodel.SignTaskAttrBitBusiness) {
				var ids []int64
				for _, v := range archiveList {
					ids = append(ids, v.ID)
				}

				ids = util.Unique(ids)
				// 查询archive服务
				archiveResult, e := s.arcRPC.Arcs(context.Background(), &v1.ArcsRequest{Aids: ids})
				if e != nil {
					err = e
					log.Error("get archive result err, err=%+v", err)
					break
				}
				for _, v := range archiveList {
					a, ok := archiveResult.Arcs[v.ID]
					// 是商单的要排除
					if !ok || a.OrderID > 0 {
						continue
					}
					finalResult = append(finalResult, v.ID)
				}

			} else {
				for _, v := range archiveList {
					finalResult = append(finalResult, v.ID)
				}
			}
			archiveCount = len(finalResult)
		}
		// 请假任务数
		var absence signmodel.SignTaskAbsence
		err = crmdb.Select("sum(absence_count) as absence_count").
			Where("task_history_id=? and state!=?", taskHistory.ID, signmodel.SignTaskAbsenceStateDelete).
			Find(&absence).Error
		if err != nil {
			log.Error("get task absence fail, task history=%+v", taskHistory)
			return
		}
		archiveCount += int(absence.AbsenceCount)
		log.Info("task count=%d, archive=%d, absence=%d, task=%+v", archiveCount, len(finalResult), absence.AbsenceCount, taskHistory)

		// 更新task history的数量
		task.TaskCounter = int32(archiveCount)

		var tx = crmdb.Begin()
		defer func() {
			if r := recover(); r != nil || err != nil {
				log.Error("roll back task update, task=%+v, r=%+v | err=%+v", task, r, err)
				tx.Rollback()
			}
		}()
		err = tx.Table(signmodel.TableNameSignTask).Where("id=?", task.ID).
			Updates(map[string]interface{}{
				"generate_date": dateStr,
			}).Error

		if err != nil {
			log.Error("update sign task fail, task=%+v, err=%+v", task, err)
			return
		}

		// update history
		var state = signmodel.SignTaskStateRunning
		if archiveCount >= int(task.TaskCondition) {
			state = signmodel.SignTaskStateFinish
		}
		err = tx.Table(signmodel.TableNameSignTaskHistory).Where("id=?", taskHistory.ID).
			Updates(map[string]interface{}{
				"task_counter":   task.TaskCounter,
				"task_condition": task.TaskCondition,
				"state":          state,
				"attribute":      task.Attribute,
				"task_type":      task.TaskType,
			}).Error
		if err != nil {
			log.Error("update sign task history fail, task=%+v, err=%+v", task, err)
			return
		}

		// update sign
		err = tx.Table(signmodel.TableNameSignUp).Where("id=?", task.SignID).
			Updates(map[string]interface{}{
				"task_state": state,
			}).Error
		if err != nil {
			log.Error("update sign up fail, task=%+v, err=%+v", task, err)
			return
		}

		err = tx.Commit().Error
		if err != nil {
			log.Error("commit err, err=%+v", err)
		}
	}
	return
}
