package service

import (
	"context"
	"fmt"
	"time"

	"go-common/app/admin/main/videoup-task/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

// List  查看任务列表
func (s *Service) List(c context.Context, uid int64, pn, ps int, ltype, leader int8) (tasks []*model.Task, err error) {
	return s.dao.ListByCondition(c, uid, pn, ps, ltype, leader)
}

// Delay 申请延迟
func (s *Service) Delay(c context.Context, id, uid int64, reason string) (err error) {
	tx, err := s.dao.BeginTran(c)
	if err != nil {
		log.Error("s.dao.BeginTran() error(%v)", err)
		return
	}
	rows, err := s.dao.TxUpTaskByID(tx, id, map[string]interface{}{"state": model.TypeDelay, "dtime": time.Now()})
	if err != nil {
		log.Error("s.dao.TxUpTaskByID(%d) error(%v)", id, err)
		tx.Rollback()
		return
	}
	if rows > 0 {
		if _, err = s.dao.TxAddTaskHis(tx, 0, model.ActionDelay /*action*/, id /*task_id*/, 0, uid /*uid*/, 0, 0, reason /*reason*/); err != nil {
			log.Error("s.dao.AddTaskLog(%d) error(%v)", id, err)
			tx.Rollback()
			return
		}
	}
	return tx.Commit()
}

// TaskSubmit 提交审核结果
func (s *Service) TaskSubmit(c context.Context, t *model.Task, uid int64, status int16) (err error) {
	var utime int64
	switch {
	case t.State == model.TypeDelay || t.State == model.TypeReview: //延迟任务，复审任务不记录utime
		utime = 0
	case t.GTime.TimeValue().IsZero():
		utime = int64(time.Since(t.MTime.TimeValue()))
	default:
		utime = int64(time.Since(t.GTime.TimeValue()))
	}

	tx, err := s.dao.BeginTran(c)
	if err != nil {
		log.Error("s.dao.BeginTran error(%v)", err)
		return
	}
	rows, err := s.dao.TxUpTaskByID(tx, t.ID, map[string]interface{}{"state": model.TypeFinished, "utime": utime})
	if err != nil {
		log.Error("s.dao.TxUpTaskByID(%d) error(%v)", t.ID, err)
		tx.Rollback()
		return
	}
	if rows > 0 {
		if _, err = s.dao.TxAddTaskHis(tx, 0, model.ActionSubmit /*action*/, t.ID /*task_id*/, t.Cid /*cid*/, uid /*uid*/, utime /*utime*/, status /*result*/, "TaskSubmit" /*reason*/); err != nil {
			log.Error("s.dao.AddTaskLog(%d) error(%v)", t.ID, err)
			tx.Rollback()
			return
		}
	}
	return tx.Commit()
}

// Next 领取任务
func (s *Service) Next(c context.Context, uid int64) (task *model.Task, err error) {
	var rows int64
	task, err = s.dao.GetNextTask(c, uid)
	if err != nil {
		log.Error("d.getTask(%d) error(%v)", uid, err)
		return
	}
	if task != nil {
		return
	}
	// 释放超时任务
	s.Free(c, 0)
	// 从实时任务池抢占
	if rows, err = s.dispatchTask(c, uid); err != nil {
		return
	} else if rows > 0 {
		return s.dao.GetNextTask(c, uid)
	}
	return
}

// Info 查询任务信息
func (s *Service) Info(c context.Context, tid int64) (task *model.Task, err error) {
	return s.dao.TaskByID(c, tid)
}

// 抢占任务(先抢占再查,避免重复下发)
func (s *Service) dispatchTask(c context.Context, uid int64) (rows int64, err error) {
	var (
		tls   []*model.TaskForLog
		arrid []int64
	)

	if tls, err = s.dao.GetDispatchTask(c, uid); err != nil {
		log.Error("s.dao.GetDispatchTask(%d) error(%v)", uid, err)
		return
	}

	for _, item := range tls {
		arrid = append(arrid, item.ID)
	}

	if len(arrid) > 0 {
		if rows, err = s.dao.UpDispatchTask(c, uid, arrid); err != nil {
			log.Error("s.dao.UpDispatchTask(%d,%v,%v) error(%v)", uid, arrid, err)
			return
		}

		// 日志允许错误
		if int(rows) == len(arrid) {
			log.Info("UpDispatchTask 更新数量(%d)", rows)
		} else {
			log.Warn("UpDispatchTask 更新数量(%d) 日志数量(%d)", rows, len(arrid))
		}
		s.dao.MulAddTaskHis(c, tls, model.ActionDispatch, uid)
	}

	return
}

// Free 任务释放(有uid为主动释放，没有uid为被动释放)(先查再释放，有可能记录冗余释放信息)
func (s *Service) Free(c context.Context, uid int64) (rows int64) {
	var (
		rts        []*model.TaskForLog
		ids, rtids []int64
		lastid     int64
		err        error
		mtime      = time.Now()
	)

	if uid == 0 {
		if rts, err = s.dao.GetTimeOutTask(c); err != nil {
			log.Error("s.Free s.dao.GetTimeOutTask error(%v)", err)
			return
		}
	} else {
		if rts, lastid, err = s.dao.GetRelTask(c, uid); err != nil {
			log.Error("s.Free s.dao.GetRelTask(%d) error(%v)", uid, err)
			return
		}
	}

	mcases := make(map[int64]*model.WCItem)
	for _, rt := range rts {
		ids = append(ids, rt.ID)
		if rt.Subject == 1 { //指派任务回流
			rtids = append(rtids, rt.ID)
			mcases[rt.ID] = &model.WCItem{Radio: 4, Weight: model.WLVConf.SubRelease,
				Mtime: model.NewFormatTime(time.Now()), Desc: "指派回流权重"}
		}
	}
	if len(ids) > 0 {
		if rows, err = s.dao.MulReleaseMtime(c, ids, mtime); err != nil {
			log.Error("s.dao.MulReleaseMtime(%v, %v) error(%v)", ids, mtime, err)
			return
		}
		if rows > 0 {
			s.dao.MulAddTaskHis(c, rts, model.ActionRelease, uid)
		}
	}

	if lastid > 0 {
		s.dao.UpGtimeByID(c, lastid, "0000-00-00 00:00:00")

		timelogout := time.Now()
		log.Info("添加延时释放任务(%d %v)", lastid, timelogout)
		time.AfterFunc(5*time.Minute, func() {
			s.releaseSpecial(timelogout, lastid, uid)
		})
	}

	if len(rtids) > 0 {
		s.setWeightConf(c, xstr.JoinInts(rtids), mcases)
	}

	return
}

func (s *Service) releaseSpecial(tout time.Time, taskid, uid int64) {
	tx, err := s.dao.BeginTran(context.TODO())
	if err != nil {
		log.Error(" s.dao.BeginTran error(%v)", err)
		return
	}
	rows, err := s.dao.TxReleaseSpecial(tx, tout, 1, taskid, uid)
	if err != nil {
		log.Error("s.dao.TxReleaseSpecial error(%v)", err)
		tx.Rollback()
		return
	}
	if rows > 0 {
		log.Info("s.dao.TxReleaseSpecial 释放任务(%d)", taskid)
		if _, err = s.dao.TxAddTaskHis(tx, 0, model.ActionRelease, taskid, 0, uid, 0, 0, "登出延时释放"); err != nil {
			log.Error("s.dao.TxAddTaskHis error(%v)", err)
			tx.Rollback()
			return
		}
	}
	tx.Commit()
}

func (s *Service) judge(c context.Context, tid, aid, cid, uid int64) (err error) {
	var (
		rows int64
		tx   *sql.Tx
		v    *model.ArcVideo
		a    *model.Archive
	)
	// 1.校验视频

	if v, err = s.dao.ArcVideoByCID(c, cid); err != nil {
		log.Error("s.dao.ArcVideoByCID(%d) error(%v)", cid, err)
		return
	}
	if v == nil || v.Status == model.VideoStatusDelete {
		err = fmt.Errorf("视频(cid=%d)被删除", cid)
		goto DELETE
	}
	// 2.校验稿件
	if a, err = s.dao.Archive(c, aid); err != nil {
		log.Error("s.dao.Archive(%d) error(%v)", aid, err)
		return
	}
	if a == nil || a.State == model.StateForbidUpDelete {
		err = fmt.Errorf("稿件(aid=%d)被删除", aid)
	}
DELETE:
	if err != nil {
		if tx, err = s.dao.BeginTran(c); err != nil {
			log.Error("s.dao.BeginTran() error(%v)", err)
			return
		}

		if rows, err = s.dao.TxUpTaskByID(tx, tid, map[string]interface{}{"state": model.TypeFinished, "utime": 0}); err != nil {
			log.Error("s.dao.TxUpTaskByID(%d) error(%v)", tid, err)
			tx.Rollback()
			return
		}
		if rows > 0 {
			if _, err = s.dao.TxAddTaskHis(tx, 0 /*pool*/, model.ActionTaskDelete /*action*/, tid /*task_id*/, cid /*cid*/, uid /*uid*/, 0 /*utime*/, model.VideoStatusDelete /*result*/, "judge delete" /*reason*/); err != nil {
				log.Error("s.dao.AddTaskLog(%d) error(%v)", tid, err)
				tx.Rollback()
				return
			}
		}
		return tx.Commit()
	}
	return
}

// CheckOwner 检查任务状态修改权限
func (s *Service) CheckOwner(c context.Context, tid, uid int64) (err error) {
	var role int8
	var rows int64
	task, err := s.dao.TaskByID(c, tid)
	if task == nil || err != nil {
		log.Error("s.dao.TaskByID(%d) error(%v)", tid, err)
		return
	}
	if err = s.judge(c, task.ID, task.Aid, task.Cid, uid); err != nil {
		log.Error("s.judge(%+v) error(%v)", task, err)
		return
	}

	if role, err = s.dao.GetUserRole(c, uid); err != nil || role == 0 {
		err = fmt.Errorf("非法用户(%d)", uid)
		return
	}

	if task.State == model.TypeDelay || task.State == model.TypeReview {
		return
	}

	if !s.CheckOnline(c, uid) {
		err = fmt.Errorf("请先签到(%d)", uid)
		return
	}

	if role == model.TaskLeader {
		return
	}
	if task.UID != uid {
		err = fmt.Errorf("没有权限处理该任务")
		return
	}
	// 普通用户处理超时了，将任务释放掉
	if task.State == model.TypeDispatched && time.Since(task.GTime.TimeValue()).Minutes() > 10.0 {
		var tx *sql.Tx
		if tx, err = s.dao.BeginTran(c); err != nil {
			log.Error("s.dao.BeginTran() error(%v)", err)
			return
		}

		if rows, err = s.dao.TxUpTaskByID(tx, tid, map[string]interface{}{"state": model.TypeRealTime, "uid": 0, "gtime": "0000-00-00 00:00:00"}); err != nil {
			log.Error("s.dao.TxUpTaskByID(%d) error(%v)", tid, err)
			tx.Rollback()
			return
		}
		if rows > 0 {
			if _, err = s.dao.TxAddTaskHis(tx, 0 /*pool*/, model.ActionRelease /*action*/, tid /*task_id*/, 0 /*cid*/, uid /*uid*/, 0 /*utime*/, 0 /*result*/, "timeout release" /*reason*/); err != nil {
				log.Error("s.dao.AddTaskLog(%d) error(%v)", tid, err)
				tx.Rollback()
				return
			}
		}
		return tx.Commit()
	}

	return
}
