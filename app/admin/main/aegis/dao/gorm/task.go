package gorm

import (
	"context"
	"database/sql"

	"go-common/app/admin/main/aegis/model/task"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
)

const (
	_submitSQL = "UPDATE task SET state=?,uid=?,utime=? WHERE id=? AND state=? AND uid=?"
)

// TxSubmit 提交任务
func (d *Dao) TxSubmit(tx *gorm.DB, opt *task.SubmitOptions, state int8) (rows int64, err error) {
	rows = tx.Exec(_submitSQL, state, opt.UID, opt.Utime, opt.TaskID, opt.OldState, opt.OldUID).RowsAffected
	return
}

// TxCloseTasks close
func (d *Dao) TxCloseTasks(tx *gorm.DB, rids []int64, uid int64) (err error) {
	err = tx.Table("task").Where("rid IN (?) AND state<?", rids, task.TaskStateSubmit).Update("state", task.TaskStateClosed).Update("uid", uid).Error
	return
}

// CloseTask .
func (d *Dao) CloseTask(c context.Context, id int64) (err error) {
	return d.orm.Table("task").Where("id=?", id).Update("state", task.TaskStateClosed).Update("uid", 399).Error
}

// TaskByRID task by rid
func (d *Dao) TaskByRID(c context.Context, rid, flowid int64) (t *task.Task, err error) {
	db := d.orm.Model(&task.Task{}).Where("rid = ? AND state<?", rid, task.TaskStateSubmit)
	if flowid > 0 {
		db = db.Where("flow_id=?", flowid)
	}
	t = &task.Task{}
	if err = db.Find(t).Error; err == gorm.ErrRecordNotFound {
		err = nil
		t = nil
	}

	return
}

// MaxWeight max weight
func (d *Dao) MaxWeight(c context.Context, bizID, flowID int64) (max int64, err error) {
	if err = d.orm.Table("task").Select("max(weight)").Where("business_id = ? AND flow_id = ?", bizID, flowID).
		Where("state = ? OR state = ?", task.TaskStateInit, task.TaskStateDispatch).Row().Scan(&max); err != nil {
		max = 0
		err = nil
	}
	return
}

// UndoStat 未完成
func (d *Dao) UndoStat(c context.Context, bizID, flowID, UID int64) (stat *task.UnDOStat, err error) {
	stat = &task.UnDOStat{}

	err = d.orm.Raw(`SELECT COUNT(CASE WHEN admin_id>0 AND state = 0 THEN 1 ELSE NULL END) assign, 
	COUNT(CASE WHEN admin_id = 0 AND state = 2 THEN 1 ELSE NULL END) delay, 
	COUNT(CASE WHEN admin_id = 0 AND state = 1 THEN 1 ELSE NULL END) normal
	FROM task WHERE business_id=? AND flow_id=? AND uid=?`, bizID, flowID, UID).Scan(stat).Error
	return
}

// TaskStat 任务详情统计
func (d *Dao) TaskStat(c context.Context, bizID, flowID, UID int64) (stat *task.Stat, err error) {
	stat = &task.Stat{}

	err = d.orm.Raw(`SELECT COUNT(CASE WHEN admin_id=0 AND state = 0 THEN 1 ELSE NULL END) normal,
		COUNT(CASE WHEN admin_id>0 AND state = 0 THEN 1 ELSE NULL END) assign,
		COUNT(CASE WHEN state=2 THEN 1 ELSE NULL END) delayTotal,
		COUNT(CASE WHEN uid=? AND state=2 THEN 1 ELSE NULL END) delayPersonal
		FROM task WHERE business_id=? AND flow_id=?`, UID, bizID, flowID).Scan(stat).Error
	return
}

// TaskListSeized 停滞任务
func (d *Dao) TaskListSeized(c context.Context, opt *task.ListOptions) (ids []int64, count int64, err error) {
	return d.tasklist(c, "seized", opt.BusinessID, opt.FlowID, opt.UID, opt.Pn, opt.Ps)
}

// TaskListDelayd 延迟任务
func (d *Dao) TaskListDelayd(c context.Context, opt *task.ListOptions) (ids []int64, count int64, err error) {
	return d.tasklist(c, "delayd", opt.BusinessID, opt.FlowID, opt.UID, opt.Pn, opt.Ps)
}

// TaskListAssignd 指派停滞任务
func (d *Dao) TaskListAssignd(c context.Context, opt *task.ListOptions) (ids []int64, count int64, err error) {
	return d.tasklist(c, "assignd", opt.BusinessID, opt.FlowID, opt.UID, opt.Pn, opt.Ps)
}

func (d *Dao) tasklist(c context.Context, ltp string, bizID, flowID, UID int64, pn, ps int) (ids []int64, count int64, err error) {
	db := d.orm.Table("task").Where("business_id=? AND flow_id=?", bizID, flowID)
	switch ltp {
	case "seized":
		db = db.Where("state=?", task.TaskStateDispatch)
	case "delayd":
		db = db.Where("state=?", task.TaskStateDelay)
	case "assignd":
		db = db.Where("state=? AND admin_id>0", task.TaskStateDispatch)
	}
	if UID > 0 {
		db = db.Where("uid=?", UID)
	}

	var rows *sql.Rows
	rows, err = db.Count(&count).Select("id").Order("weight DESC").Offset((pn - 1) * ps).Limit(ps).Rows()
	if err != nil {
		log.Error("tasklist error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		if err = rows.Scan(&id); err != nil {
			log.Error("tasklist error(%v)", err)
			return
		}
		ids = append(ids, id)
	}
	return
}

//TaskHitAuditing 检查资源是否正在审核
func (d *Dao) TaskHitAuditing(c context.Context, rids []int64) (map[int64]struct{}, error) {
	hitids := make(map[int64]struct{})
	rows, err := d.orm.Table("task").Select("rid").Where("rid IN (?)", rids).
		Where("state = ? AND gtime!=0", task.TaskStateDispatch).Rows()
	if err != nil {
		return hitids, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		if err = rows.Scan(&id); err != nil {
			return hitids, err
		}
		hitids[id] = struct{}{}
	}
	return hitids, err
}
