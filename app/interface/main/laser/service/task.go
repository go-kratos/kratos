package service

import (
	"context"

	"go-common/app/interface/main/laser/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

// QueryTaskState is query satisfied logDate from memcache if key exist,Otherwise query from db.
func (s *Service) QueryUndoneTaskLogdate(c context.Context, mid int64, platform int, sourceType int) (logDate int64, err error) {
	// cache flag: cached to memcache only when query db.
	cache := true
	var v *model.TaskInfo
	if v, err = s.dao.TaskInfoCache(c, mid); err != nil {
		err = nil
		cache = false
	} else if v != nil {
		// if memcache key exist
		if v.Empty {
			return
		}
		if logDateExist(v, sourceType, platform) {
			logDate = v.LogDate.Time().Unix()
			return
		}
		return
	}

	// if memcache key not exist , then query db.
	v, err = s.dao.QueryUndoneTaskInfo(c, mid)
	if err != nil {
		return
	}
	if v != nil {
		if logDateExist(v, sourceType, platform) {
			logDate = v.LogDate.Time().Unix()
		}
	} else {
		// if not found from db, new empty Instance and set Flag:Empty true.
		v = &model.TaskInfo{
			Empty: true,
		}
	}

	if cache {
		s.addCache(func() {
			s.dao.AddTaskInfoCache(context.Background(), mid, v)
		})
	}
	return
}

func logDateExist(t *model.TaskInfo, sourceType int, platform int) bool {
	if sourceType == t.SourceType && (platform == t.Platform || platform == model.ALL_PLATFORM) {
		return true
	}
	return false
}

// UpdateTaskState is update task set state = 1 ,insert taskLog to table task_log and remove memcache.
func (s *Service) UpdateTaskState(c context.Context, mid int64, build string, platform int, taskState int, reason string) (err error) {
	// avoid err when remove the specified-mid task  on laser-admin application.
	taskID, err := s.dao.QueryTaskID(c, mid)
	if err != nil {
		log.Error("s.UpdateTaskState() error(%v)", err)
		return
	}
	if taskID == 0 {
		return
	}
	tx, err := s.dao.BeginTran(c)
	if err != nil {
		log.Error("s.dao.BeginTran() error(%v)", err)
		return
	}
	rows, err := s.dao.TxUpdateTaskState(c, tx, 1, taskID)
	if err != nil || rows <= 0 {
		tx.Rollback()
		log.Error("s.UpdateTaskState() error(%v)", err)
		return
	}
	err = s.addTaskLog(c, tx, taskID, mid, build, platform, taskState, reason)
	if err != nil {
		tx.Rollback()
		log.Error("s.UpdateTaskState() error(%v)", err)
		return
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit() error(%v)", err)
		return
	}
	s.dao.RemoveTaskInfoCache(c, mid)
	go s.SendEmail(context.Background(), taskID)
	return
}

// addTaskLog is insert TaskLog to table task_log.
func (s *Service) addTaskLog(c context.Context, tx *xsql.Tx, taskID int64, mid int64, build string, platform int, taskState int, reason string) (err error) {
	insertID, err := s.dao.TxAddTaskLog(c, tx, taskID, mid, build, platform, taskState, reason)
	if err != nil || insertID <= 0 {
		log.Error("s.UpdateTaskState() error(%v)", err)
		return
	}
	return
}
