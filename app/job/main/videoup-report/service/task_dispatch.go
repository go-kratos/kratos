package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-common/app/job/main/videoup-report/model/archive"
	tmod "go-common/app/job/main/videoup-report/model/task"
	"go-common/app/job/main/videoup-report/model/utils"
	"go-common/library/database/sql"
	"go-common/library/log"
)

func (s *Service) hdlVideoTask(c context.Context, fn string) (err error) {
	var (
		v     *archive.Video
		a     *archive.Archive
		state int8
		dID   int64
	)
	if v, a, err = s.archiveVideo(c, fn); err != nil {
		log.Error("s.archiveVideo(%s) error(%v)", fn, err)
		return
	}
	if a.State == archive.StateForbidUpDelete {
		log.Info("task archive(%d) deleted", a.ID)
		return
	}

	if v.Status != archive.VideoStatusWait {
		log.Info("task archive(%d) filename(%s) already status(%d)", a.ID, v.Filename, v.Status)
		return
	}

	if dID, state, err = s.arc.DispatchState(c, v.Aid, v.Cid); err != nil {
		log.Error("task s.arc.DispatchState(%d,%d) error(%v)", v.Aid, v.Cid, err)
		return
	}
	if dID != 0 && state <= tmod.StateForTaskWork {
		log.Info("task aid(%d) cid(%d) filename(%s) already in dispatch state(%d)", v.Aid, v.Cid, v.Filename, state)
		return
	}

	log.Info("archive(%d) filename(%s) video(%d) tranVideoTask begin", a.ID, v.Filename, v.Cid)
	if err = s.addVideoTask(c, a, v); err != nil {
		log.Error("task s.addVideoTask error(%v)", err)
		return
	}
	return
}

func (s *Service) archiveVideo(c context.Context, filename string) (v *archive.Video, a *archive.Archive, err error) {
	if v, err = s.arc.NewVideo(c, filename); err != nil {
		log.Error("s.arc.NewVideo(%s) error(%v)", filename, err)
		return
	}
	if v == nil {
		log.Error("s.arc.NewVideo(%s) video is nil", filename)
		err = fmt.Errorf("video(%s) is not exists", filename)
		return
	}
	if a, err = s.arc.ArchiveByAid(c, v.Aid); err != nil {
		log.Error("s.arc.ArchiveByAid(%d) filename(%s) error(%v)", v.Aid, filename, err)
		return
	}
	return
}

func (s *Service) addVideoTask(c context.Context, a *archive.Archive, v *archive.Video) (err error) {
	var (
		task         *tmod.Task
		lastID, fans int64
		cfitems      []*tmod.ConfigItem
		descb        []byte
		accfailed    bool
	)

	task = s.setTaskAssign(c, a, v)
	fans, accfailed = s.setTaskUPSpecial(c, task, a.Mid)
	s.setTaskTimed(c, task)
	cfitems = s.getConfWeight(c, task, a)

	if lastID, err = s.arc.AddDispatch(c, task); err != nil {
		log.Error("s.arc.AddDispatch error(%v)", err)
		return
	}
	// 允许日志记录错误
	if _, err = s.arc.AddTaskHis(c, tmod.PoolForFirst, 6, lastID, v.Cid, task.UID, v.Status, "videoup-job"); err != nil {
		log.Error("s.arc.AddTaskHis error(%v)", err)
	}

	log.Info("archive(%d) filename(%s) taskUid(%d)", a.ID, v.Filename, task.UID)

	// 保存权重配置信息,错误不影响正常流程
	tp := &tmod.WeightParams{
		TaskID:    lastID,
		Mid:       a.Mid,
		Special:   task.UPSpecial,
		Ctime:     utils.NewFormatTime(time.Now()),
		Ptime:     task.Ptime,
		CfItems:   cfitems,
		Fans:      fans,
		AccFailed: accfailed,

		TypeID: a.TypeID,
	}

	s.setTaskUpFrom(c, a.ID, tp)
	s.setTaskUpGroup(c, a.Mid, tp)

	s.redis.SetWeight(c, map[int64]*tmod.WeightParams{lastID: tp})
	if len(cfitems) > 0 {
		if descb, err = json.Marshal(cfitems); err != nil {
			log.Error("json.Marshal error(%v)", err)
		} else {
			if _, err = s.arc.InDispatchExtend(c, lastID, string(descb)); err != nil {
				log.Error("s.task.InDispatchExtend(%d) error(%v)", lastID, err)
			}
		}
		err = nil
	}
	return
}

func (s *Service) moveDispatch() (err error) {
	var (
		tx               *sql.Tx
		c                = context.TODO()
		mtime            = time.Now().Add(-24 * time.Hour)
		startTime        = time.Date(mtime.Year(), mtime.Month(), mtime.Day(), mtime.Hour(), 0, 0, 0, mtime.Location())
		endTime          = time.Date(mtime.Year(), mtime.Month(), mtime.Day(), mtime.Hour(), 59, 59, 0, mtime.Location())
		dispatchRows     int64
		dispatchDoneRows int64
	)
	if tx, err = s.arc.BeginTran(c); err != nil {
		log.Error("s.arc.BeginTran error(%v)")
		return
	}
	if dispatchRows, err = s.arc.TxAddDispatchDone(c, tx, startTime, endTime); err != nil {
		tx.Rollback()
		log.Error("s.arc.TxAddDispatchDone(%s) error(%v)", mtime.Format("2006-01-02 15:04:05"), err)
		return
	}
	if dispatchDoneRows, err = s.arc.TxDelDispatchByTime(c, tx, startTime, endTime); err != nil {
		tx.Rollback()
		log.Error("s.arc.TxDelDispatchByTime(%s) error(%v)", mtime.Format("2006-01-02 15:04:05"), err)
		return
	}
	if dispatchRows != dispatchDoneRows {
		// no way here !
		tx.Rollback()
		log.Error("moveDispatch error dispatchRows(%d) dispatchDoneRows(%d)", dispatchRows, dispatchDoneRows)
		return
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit error(%v)")
		return
	}
	log.Info("moveDispatch mtime(%s) to mtime(%s) rows(%d)", startTime.Format("2006-01-02 15:04:05"), endTime.Format("2006-01-02 15:04:05"), dispatchDoneRows)
	return
}

func (s *Service) moveTaskOperHis(limit int64) (moved int64, err error) {
	var (
		tx        *sql.Tx
		c         = context.TODO()
		mtime     = time.Now().Add(-2 * 30 * 24 * time.Hour)
		before    = time.Date(mtime.Year(), mtime.Month(), mtime.Day(), mtime.Hour(), 0, 0, 0, mtime.Location())
		movedRows int64
		delRows   int64
	)
	if tx, err = s.arc.BeginTran(c); err != nil {
		log.Error("s.arc.BeginTran error(%v)")
		return
	}
	if movedRows, err = s.arc.TxMoveTaskOperDone(tx, before, limit); err != nil {
		tx.Rollback()
		log.Error("s.arc.TxMoveTaskOperDone(%s) error(%v)", mtime.Format("2006-01-02 15:04:05"), err)
		return
	}
	if delRows, err = s.arc.TxDelTaskOper(tx, before, limit); err != nil {
		tx.Rollback()
		log.Error("s.arc.TxDelTaskOper(%s) error(%v)", mtime.Format("2006-01-02 15:04:05"), err)
		return
	}
	if movedRows != delRows {
		tx.Rollback()
		log.Error("moveOperHistory error mvRows(%d) delRows(%d)", movedRows, delRows)
		return
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit error(%v)")
		return
	}
	log.Info("moveTaskOperHistory before mtime(%s) rows(%d)", before.Format("2006-01-02 15:04:05"), movedRows)
	return movedRows, nil
}

func (s *Service) delTaskDispatchDone(limit int64) (delRows int64, err error) {
	var (
		c      = context.TODO()
		mtime  = time.Now().Add(-30 * 24 * time.Hour)
		before = time.Date(mtime.Year(), mtime.Month(), mtime.Day(), mtime.Hour(), 0, 0, 0, mtime.Location())
	)
	if delRows, err = s.arc.DelTaskDoneBefore(c, before, limit); err != nil {
		log.Error("s.arc.DelTaskDoneBefore(%s) error(%v)", mtime.Format("2006-01-02 15:04:05"), err)
		return
	}
	if delRows > 0 {
		log.Info("delTaskDispatchDone before mtime(%s) rows(%d)", before.Format("2006-01-02 15:04:05"), delRows)
	}

	if delRows, err = s.arc.DelTaskBefore(c, before, limit); err != nil {
		log.Error("s.arc.DelTaskBefore(%s) error(%v)", mtime.Format("2006-01-02 15:04:05"), err)
		return
	}
	if delRows > 0 {
		log.Info("DelTaskBefore before mtime(%s) rows(%d)", before.Format("2006-01-02 15:04:05"), delRows)
	}
	return
}

func (s *Service) delTaskHistoryDone(limit int64) (delRows int64, err error) {
	var (
		c      = context.TODO()
		mtime  = time.Now().Add(-3 * 30 * 24 * time.Hour)
		before = time.Date(mtime.Year(), mtime.Month(), mtime.Day(), mtime.Hour(), 0, 0, 0, mtime.Location())
	)
	if delRows, err = s.arc.DelTaskHistoryDone(c, before, limit); err != nil {
		log.Error("s.arc.DelTaskHistoryDone(%s) error(%v)", mtime.Format("2006-01-02 15:04:05"), err)
		return
	}
	log.Info("delTaskHistoryDone before mtime(%s) rows(%d)", before.Format("2006-01-02 15:04:05"), delRows)
	return
}

func (s *Service) delTaskExtend(limit int64) (delRows int64, err error) {
	var (
		c      = context.TODO()
		mtime  = time.Now().Add(-20 * 24 * time.Hour)
		before = time.Date(mtime.Year(), mtime.Month(), mtime.Day(), mtime.Hour(), 0, 0, 0, mtime.Location())
	)
	if delRows, err = s.arc.DelTaskExtend(c, before, limit); err != nil {
		log.Error("s.arc.DelTaskExtend(%s) error(%v)", mtime.Format("2006-01-02 15:04:05"), err)
		return
	}
	log.Info("delTaskExtend before mtime(%s) rows(%d)", before.Format("2006-01-02 15:04:05"), delRows)
	return
}

/*
	1.移动task_dispatch到task_dispatch_done
	2.移动task_oper_history到task_oper_history_done
	3.删除过于久远的task_dispatch_done,task_oper_history_done,task_dispatch_extend
*/
func (s *Service) movetaskproc() {
	defer s.waiter.Done()
	for {
		if s.closed {
			return
		}
		s.moveDispatch()
		time.Sleep(1 * time.Hour)
	}
}

func (s *Service) deltaskproc() {
	for {
		for {
			rows, _ := s.moveTaskOperHis(100)
			time.Sleep(1 * time.Second)
			if rows == 0 {
				break
			}
		}
		for {
			rows, _ := s.delTaskDispatchDone(100)
			time.Sleep(1 * time.Second)
			if rows == 0 {
				break
			}
		}
		for {
			rows, _ := s.delTaskHistoryDone(100)
			time.Sleep(1 * time.Second)
			if rows == 0 {
				break
			}
		}
		for {
			rows, _ := s.delTaskExtend(100)
			time.Sleep(1 * time.Second)
			if rows == 0 {
				break
			}
		}
		time.Sleep(nextDay(10))
	}
}
