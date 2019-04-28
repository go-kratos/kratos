package dao

import (
	"context"
	"time"

	"go-common/app/admin/main/videoup-task/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_avgUtimeSQL     = "SELECT IFNULL(avg(utime),0.0) FROM task_oper_history WHERE action=2 AND uid=? AND ctime>=? AND ctime<?"
	_sumDurationSQL  = "SELECT IFNULL(sum(v.duration),0) FROM task_oper_history as t LEFT JOIN video as v ON t.cid=v.id WHERE action IN (2,5) AND t.uid=? AND t.ctime>=? AND t.ctime<?"
	_actionCountSQL  = "SELECT action,count(*) FROM task_oper_history WHERE uid=? AND ctime>=? AND ctime<? AND action IN (2,5,6,7) GROUP BY action"
	_passCountSQL    = "SELECT count(*) FROM task_oper_history WHERE action IN (2,5) AND uid=? AND ctime>=? AND ctime<? AND result=0"
	_subjectCountSQL = `SELECT count(*) FROM (
							SELECT id FROM task_dispatch WHERE uid=? AND state=? AND mtime>=? AND ctime<? AND subject=1 
							UNION ALL SELECT task_id as id FROM task_dispatch_done WHERE uid=? AND state=? AND mtime>=? AND ctime<? AND subject=1) as t
						`
	_activeUidsSQL = "SELECT DISTINCT uid from task_oper_history WHERE ctime>=? AND ctime<? AND uid!=0 AND action IN (2,5)"
)

// AvgUtimeByUID 平均处理耗时, 只统计action=2的
func (d *Dao) AvgUtimeByUID(c context.Context, uid int64, stime, etime time.Time) (utime float64, err error) {
	if err = d.arcReadDB.QueryRow(c, _avgUtimeSQL, uid, stime, etime).Scan(&utime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		log.Error("d.arcReadDB.QueryRow error(%v)", err)
	}
	return
}

// SumDurationByUID 视频总时长，统计action=2,5的
func (d *Dao) SumDurationByUID(c context.Context, uid int64, stime, etime time.Time) (duration int64, err error) {
	if err = d.arcReadDB.QueryRow(c, _sumDurationSQL, uid, stime, etime).Scan(&duration); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		log.Error("d.arcReadDB.QueryRow error(%v)", err)
	}
	return
}

// ActionCountByUID 操作个数统计
func (d *Dao) ActionCountByUID(c context.Context, uid int64, stime, etime time.Time) (mapAction map[int8]int64, err error) {
	mapAction = make(map[int8]int64)
	rows, err := d.arcReadDB.Query(c, _actionCountSQL, uid, stime, etime)
	if err != nil {
		log.Error("d.arcReadDB.Query error(%v)", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			action int8
			count  int64
		)
		if err = rows.Scan(&action, &count); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		mapAction[action] = count
	}
	return
}

// PassCountByUID 总过审个数
func (d *Dao) PassCountByUID(c context.Context, uid int64, stime, etime time.Time) (count int64, err error) {
	if err = d.arcReadDB.QueryRow(c, _passCountSQL, uid, stime, etime).Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		log.Error("d.arcReadDB.QueryRow error(%v)", err)
	}
	return
}

// SubjectCountByUID 总指派个数
func (d *Dao) SubjectCountByUID(c context.Context, uid int64, stime, etime time.Time) (count int64, err error) {
	if err = d.arcReadDB.QueryRow(c, _subjectCountSQL, uid, model.TaskStateCompleted, stime, etime, uid, model.TaskStateCompleted, stime, etime).Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		log.Error("d.arcReadDB.QueryRow error(%v)", err)
	}
	return
}

// ActiveUids 统计24小时内有提交的
func (d *Dao) ActiveUids(c context.Context, stime, etime time.Time) (uids []int64, err error) {
	st := time.Now()
	defer func() {
		log.Info("ActiveUids du(%.2fm) wait(%.2fs)", etime.Sub(stime).Minutes(), time.Since(st).Seconds())
	}()
	rows, err := d.arcReadDB.Query(c, _activeUidsSQL, stime, etime)
	if err != nil {
		log.Error("d.arcReadDB.Query error(%v)", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var uid int64
		if err = rows.Scan(&uid); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		uids = append(uids, uid)
	}
	return
}
