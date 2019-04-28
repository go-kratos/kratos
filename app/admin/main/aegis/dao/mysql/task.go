package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"go-common/app/admin/main/aegis/model/common"
	modtask "go-common/app/admin/main/aegis/model/task"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_taskSQL      = "SELECT id,business_id,flow_id,rid,admin_id,uid,state,weight,utime,gtime,mid,fans,`group`,reason,ctime,mtime from task WHERE id=?"
	_listCheckSQL = "SELECT id FROM task WHERE id IN (%s)"

	_dispatchByIDSQL = "UPDATE task SET gtime=? WHERE id=? AND state=? AND uid=? AND gtime=0"
	_queryGtimeSQL   = "SELECT gtime FROM task WHERE id=? AND state=? AND uid=?"
	_dispatchSQL     = "UPDATE task SET gtime=? WHERE state=? AND uid=? AND gtime='0000-00-00 00:00:00' ORDER BY weight LIMIT ?"
	_releaseSQL      = "UPDATE task SET admin_id=0,uid=0,state=0,gtime='0000-00-00 00:00:00' WHERE business_id=? AND flow_id=? AND uid=? AND (state=? OR (state=0 AND admin_id>0))"
	_resetGtimeSQL   = "UPDATE task SET gtime='0000-00-00 00:00:00' WHERE state=? AND business_id=? AND flow_id=? AND uid=?"
	_seizeSQL        = "UPDATE task SET state=?,uid=? WHERE id=? AND state=?"
	_submitSQL       = "UPDATE task SET state=?,uid=?,utime=? WHERE id=? AND state=? AND uid=?"
	_delaySQL        = "UPDATE task SET state=?,uid=?,reason=?,gtime='0000-00-00 00:00:00' WHERE id=? AND state=? AND uid=?"

	_consumerSQL     = "INSERT INTO task_consumer (business_id,flow_id,uid,state) VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE state=?"
	_onlinesSQL      = "SELECT uid,mtime FROM task_consumer WHERE business_id=? AND flow_id=? AND state=?"
	_isconsumerOnSQL = "SELECT state FROM task_consumer WHERE business_id=? AND flow_id=? AND uid=?"

	_queryTaskSQL     = "SELECT id,business_id,flow_id,uid,weight FROM task WHERE state=? AND mtime<=? AND id>? ORDER BY id LIMIT ?"
	_countPersonalSQL = "SELECT count(*) FROM task WHERE state=? AND business_id=? AND flow_id=? AND uid=?"
	_queryForSeizeSQL = "SELECT id FROM task WHERE state=? AND business_id=? AND flow_id=? AND uid IN (0,?) ORDER BY weight DESC LIMIT ?"
	_listTasksSQL     = "SELECT `id`,`business_id`,`flow_id`,`rid`,`admin_id`,`uid`,`state`,`weight`,`utime`,`gtime`,`mid`,`fans`,`group`,`reason`,`ctime`,`mtime` FROM task %s ORDER BY weight DESC LIMIT ?,?"
)

// TaskFromDB .
func (d *Dao) TaskFromDB(c context.Context, id int64) (task *modtask.Task, err error) {
	task = &modtask.Task{}
	err = d.db.QueryRow(c, _taskSQL, id).
		Scan(&task.ID, &task.BusinessID, &task.FlowID, &task.RID, &task.AdminID, &task.UID, &task.State,
			&task.Weight, &task.Utime, &task.Gtime, &task.MID, &task.Fans, &task.Group, &task.Reason, &task.Ctime, &task.Mtime)
	if err != nil {
		task = nil
		if err == sql.ErrNoRows {
			log.Error("TaskFromDB(%d) norows", id)
			err = nil
			return
		}
		log.Error("TaskFromDB(%d) error(%v)", id, errors.WithStack(err))
	}
	return
}

// DispatchByID 派遣任务,更新gtime
func (d *Dao) DispatchByID(c context.Context, mtasks map[int64]*modtask.Task, ids []int64, args ...interface{}) (missids map[int64]struct{}, err error) {
	var (
		gtime = time.Now()
		uid   = args[0].(int64)
	)

	missids = make(map[int64]struct{})

	for _, id := range ids {
		var (
			rows int64
			gt   time.Time
			res  sql.Result
		)
		if err = d.db.QueryRow(c, _queryGtimeSQL, id, modtask.TaskStateDispatch, uid).Scan(&gt); err != nil {
			if err == sql.ErrNoRows {
				missids[id] = struct{}{}
				err = nil
				continue
			}
			log.Error("d.db.QueryRow error(%v)", errors.WithStack(err))
			return
		}

		if gt.IsZero() {
			res, err = d.db.Exec(c, _dispatchByIDSQL, gtime, id, modtask.TaskStateDispatch, uid)
			if err != nil {
				log.Error("Exec error(%v)", errors.WithStack(err))
				return
			}
			if rows, err = res.RowsAffected(); err != nil {
				log.Error("RowsAffected error(%v)", errors.WithStack(err))
				return
			}
			if rows == 0 {
				missids[id] = struct{}{}
			} else {
				mtasks[id].Gtime = common.IntTime(gtime.Unix())
			}
		} else {
			mtasks[id].Gtime = common.IntTime(gt.Unix())
		}
	}
	return
}

// DBDispatch 直接数据库派遣
func (d *Dao) DBDispatch(c context.Context, opt *modtask.NextOptions) (tasks []*modtask.Task, count int64, err error) {
	var (
		res   sql.Result
		gtime = time.Now()
	)

	// 1.直接更新派遣时间
	res, err = d.db.Exec(c, _dispatchSQL, gtime, modtask.TaskStateDispatch, opt.UID, opt.DispatchCount)
	if err != nil {
		log.Error("Exec error(%v)", errors.WithStack(err))
		return
	}
	if count, err = res.RowsAffected(); err != nil {
		log.Error("RowsAffected error(%v)", errors.WithStack(err))
		return
	}

	// 2.读取任务
	wherecache := fmt.Sprintf("WHERE state=%d AND uid=%d AND gtime!='0000-00-00 00:00:00'", modtask.TaskStateDispatch, opt.UID)

	return d.listTasks(c, &modtask.ListOptions{BaseOptions: opt.BaseOptions, Pager: common.Pager{Pn: 1, Ps: int(opt.DispatchCount)}}, wherecache)
}

// Release 释放任务
func (d *Dao) Release(c context.Context, opt *common.BaseOptions, delay bool) (rows int64, err error) {
	sql := _releaseSQL
	if delay {
		sql = _releaseSQL + " AND gtime='0000-00-00 00:00:00'"
	}

	log.Info("Mysql Release(%+v) delay(%v)", opt, delay)
	res, err := d.db.Exec(c, sql, opt.BusinessID, opt.FlowID, opt.UID, modtask.TaskStateDispatch)
	if err != nil {
		log.Error("db.Exec(%s)[%d,%d,%d,%d] error(%v)", sql, opt.BusinessID, opt.FlowID, opt.UID, modtask.TaskStateDispatch, err)
		return
	}
	// 已经下发的延迟5分钟释放
	if delay {
		_, err = d.db.Exec(c, _resetGtimeSQL, modtask.TaskStateDispatch, opt.BusinessID, opt.FlowID, opt.UID)
		if err != nil {
			log.Error("db.Exec(%s)[%d,%d,%d,%d] error(%v)", sql, modtask.TaskStateDispatch, opt.BusinessID, opt.FlowID, opt.UID, err)
		}

		time.AfterFunc(5*time.Minute, func() {
			d.Release(context.Background(), opt, false)
		})
	}

	return res.RowsAffected()
}

// Seize 抢占任务
func (d *Dao) Seize(c context.Context, mapids map[int64]int64) (count int64, err error) {
	tx, err := d.db.Begin(c)
	if err != nil {
		log.Error("d.Seize.Begin error(%v)", errors.WithStack(err))
		return
	}
	defer tx.Commit()

	for tid, uid := range mapids {
		var (
			rows int64
			res  sql.Result
		)
		res, err = tx.Exec(_seizeSQL, modtask.TaskStateDispatch, uid, tid, modtask.TaskStateInit)
		if err != nil {
			log.Error("Exec error(%v)", errors.WithStack(err))
			tx.Rollback()
			return
		}
		if rows, err = res.RowsAffected(); err != nil {
			log.Error("RowsAffected error(%v)", errors.WithStack(err))
			tx.Rollback()
			return
		}
		if rows == 1 {
			count++
		}
	}
	return
}

// Delay 延迟任务
func (d *Dao) Delay(c context.Context, opt *modtask.DelayOptions) (rows int64, err error) {
	var (
		res sql.Result
	)
	res, err = d.db.Exec(c, _delaySQL, modtask.TaskStateDelay, opt.UID, opt.Reason, opt.TaskID, modtask.TaskStateDispatch, opt.UID)
	if err != nil {
		log.Error("Exec error(%v)", errors.WithStack(err))
		return
	}
	if rows, err = res.RowsAffected(); err != nil {
		log.Error("RowsAffected error(%v)", errors.WithStack(err))
		return
	}
	return
}

// ListCheckUnSeized .
func (d *Dao) ListCheckUnSeized(c context.Context, mtasks map[int64]*modtask.Task, ids []int64, args ...interface{}) (missids map[int64]struct{}, err error) {
	wherecase := fmt.Sprintf("state = %d", modtask.TaskStateInit)
	return d.listCheck(c, wherecase, ids)
}

// ListCheckSeized .
func (d *Dao) ListCheckSeized(c context.Context, mtasks map[int64]*modtask.Task, ids []int64, args ...interface{}) (missids map[int64]struct{}, err error) {
	if len(args) < 1 {
		return
	}
	uid := args[0].(int64)
	wherecase := fmt.Sprintf("state = %d", modtask.TaskStateDispatch)
	if uid != 0 {
		wherecase += fmt.Sprintf(" AND uid=%d", uid)
	}
	return d.listCheck(c, wherecase, ids)
}

// ListCheckDelay .
func (d *Dao) ListCheckDelay(c context.Context, mtasks map[int64]*modtask.Task, ids []int64, args ...interface{}) (missids map[int64]struct{}, err error) {
	if len(args) < 1 {
		return
	}
	uid := args[0].(int64)
	wherecase := fmt.Sprintf("state=%d", modtask.TaskStateDelay)
	if uid != 0 {
		wherecase += fmt.Sprintf(" AND uid=%d", uid)
	}
	return d.listCheck(c, wherecase, ids)
}

// ListTasks .
func (d *Dao) ListTasks(c context.Context, opt *modtask.ListOptions) (tasks []*modtask.Task, count int64, err error) {
	var (
		wherecase string
		cases     []string
		state     int8
		isDefault bool
	)

	switch opt.State {
	case 1:
		state = modtask.TaskStateInit
	case 2:
		state = modtask.TaskStateDispatch
	case 3:
		state = modtask.TaskStateDelay
	case 4:
		state = modtask.TaskStateDispatch
		cases = append(cases, "admin_id>0")
	default:
		isDefault = true
		cases = append(cases, fmt.Sprintf("state<%d", modtask.TaskStateSubmit))
	}
	if !isDefault {
		cases = append(cases, fmt.Sprintf("state=%d", state))
		if !opt.BisLeader && (opt.State == 2 || opt.State == 3 || opt.State == 4) {
			cases = append(cases, fmt.Sprintf("uid=%d", opt.UID))
		}
	}

	wherecase = fmt.Sprintf("WHERE business_id=%d AND flow_id=%d AND ", opt.BusinessID, opt.FlowID) + strings.Join(cases, " AND ")
	return d.listTasks(c, opt, wherecase)
}

func (d *Dao) listTasks(c context.Context, opt *modtask.ListOptions, wherecase string) (tasks []*modtask.Task, count int64, err error) {
	countSQL := fmt.Sprintf("SELECT count(*) FROM task %s", wherecase)

	if err = d.db.QueryRow(c, countSQL).Scan(&count); err != nil {
		log.Error("QueryRow error(%v)", err)
		return
	}

	if count > 0 {
		var (
			rows    *xsql.Rows
			listSQL = fmt.Sprintf(_listTasksSQL, wherecase)
		)

		if rows, err = d.db.Query(c, listSQL, (opt.Pn-1)*opt.Ps, opt.Pn*opt.Ps); err != nil {
			log.Error("Query error(%v)", err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			task := &modtask.Task{}
			if err = rows.Scan(&task.ID, &task.BusinessID, &task.FlowID, &task.RID, &task.AdminID, &task.UID, &task.State,
				&task.Weight, &task.Utime, &task.Gtime, &task.MID, &task.Fans, &task.Group, &task.Reason, &task.Ctime, &task.Mtime); err != nil {
				log.Error("Scan error(%v)", err)
				return
			}
			tasks = append(tasks, task)
		}
	}
	return
}

func (d *Dao) listCheck(c context.Context, wherecase string, ids []int64) (missids map[int64]struct{}, err error) {
	if len(ids) == 0 {
		return
	}
	missids = make(map[int64]struct{})
	mapids := make(map[int64]struct{})

	log.Info("listCheck ids(%v)", ids)
	defer func() {
		log.Info("listCheck missids(%v)", missids)
	}()

	for _, id := range ids {
		mapids[id] = struct{}{}
	}

	var (
		rows      *xsql.Rows
		sqlstring = fmt.Sprintf(_listCheckSQL, xstr.JoinInts(ids)) + " AND " + wherecase
	)

	if rows, err = d.db.Query(c, sqlstring); err != nil {
		log.Error("db.Query(%s) error(%v)", sqlstring, errors.WithStack(err))
		return
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		if err = rows.Scan(&id); err != nil {
			log.Error("rows.Scan error(%v)", errors.WithStack(err))
			return
		}
		delete(mapids, id)
	}

	for id := range mapids {
		missids[id] = struct{}{}
	}

	return
}

// ConsumerOn .
func (d *Dao) ConsumerOn(c context.Context, opt *common.BaseOptions) (err error) {
	return d.consumer(c, opt, modtask.ActionConsumerOn)
}

// ConsumerOff .
func (d *Dao) ConsumerOff(c context.Context, opt *common.BaseOptions) (err error) {
	return d.consumer(c, opt, modtask.ActionConsumerOff)
}

// IsConsumerOn .
func (d *Dao) IsConsumerOn(c context.Context, opt *common.BaseOptions) (on bool, err error) {
	var state int8
	if err = d.db.QueryRow(c, _isconsumerOnSQL, opt.BusinessID, opt.FlowID, opt.UID).Scan(&state); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Error("d.db.QueryRow error(%v)", err)
		return
	}
	if state == modtask.ActionConsumerOn {
		on = true
	}
	return
}

func (d *Dao) consumer(c context.Context, opt *common.BaseOptions, action int8) (err error) {
	var (
		res sql.Result
	)
	res, err = d.db.Exec(c, _consumerSQL, opt.BusinessID, opt.FlowID, opt.UID, action, action)
	if err != nil {
		log.Error("Exec error(%v)", errors.WithStack(err))
		return
	}
	if _, err = res.RowsAffected(); err != nil {
		log.Error("RowsAffected error(%v)", errors.WithStack(err))
		return
	}
	return
}

// ConsumerStat 24小时内有活动或者在线的用户
func (d *Dao) ConsumerStat(c context.Context, bizid, flowid int64) (items []*modtask.WatchItem, err error) {
	var rows *xsql.Rows
	sql := "SELECT uid,mtime,state from task_consumer where business_id=? AND flow_id=? AND (mtime > ? or state=1) order by mtime desc"
	if rows, err = d.db.Query(c, sql, bizid, flowid, time.Now().Add(-24*time.Hour)); err != nil {
		log.Error("ConsumerStat error(%v)", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		item := &modtask.WatchItem{}
		if err = rows.Scan(&item.UID, &item.Mtime, &item.State); err != nil {
			log.Error("ConsumerStat error(%v)", err)
			return
		}
		items = append(items, item)
	}
	return
}

// Onlines 在线列表
func (d *Dao) Onlines(c context.Context, opt *common.BaseOptions) (uids map[int64]time.Time, err error) {
	var (
		rows *xsql.Rows
	)

	rows, err = d.db.Query(c, _onlinesSQL, opt.BusinessID, opt.FlowID, modtask.ActionConsumerOn)
	if err != nil {
		log.Error("db.Query error(%v)", err)
		return
	}
	defer rows.Close()

	uids = make(map[int64]time.Time)
	for rows.Next() {
		var (
			uid   int64
			mtime time.Time
		)
		if err = rows.Scan(&uid, &mtime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		uids[uid] = mtime
	}
	return
}

// QueryTask .
func (d *Dao) QueryTask(c context.Context, state int8, mtime time.Time, id, limit int64) (tasks []*modtask.Task, lastid int64, err error) {
	var rows *xsql.Rows
	rows, err = d.db.Query(c, _queryTaskSQL, state, mtime, id, limit)
	if err != nil {
		log.Error("db.Query error(%v)", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		task := &modtask.Task{}
		if err = rows.Scan(&task.ID, &task.BusinessID, &task.FlowID, &task.UID, &task.Weight); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}

		tasks = append(tasks, task)
		lastid = task.ID
	}
	return
}

// CountPersonal count personal task
func (d *Dao) CountPersonal(c context.Context, opt *common.BaseOptions) (count int64, err error) {
	if err = d.db.QueryRow(c, _countPersonalSQL, modtask.TaskStateDispatch, opt.BusinessID, opt.FlowID, opt.UID).Scan(&count); err != nil {
		log.Error("QueryRow error(%v)", errors.WithStack(err))
		return
	}
	return
}

// QueryForSeize 查询当前可抢占的任务
func (d *Dao) QueryForSeize(c context.Context, businessID, flowID, uid, seizecount int64) (hitids []int64, err error) {
	log.Info("task-QueryForSeize businessID(%d), flowID(%d), uid(%d), seizecount(%d)", businessID, flowID, uid, seizecount)
	defer func() { log.Info("task-QueryForSeize hitids(%v), err(%v)", hitids, err) }()
	var rows *xsql.Rows
	rows, err = d.db.Query(c, _queryForSeizeSQL, modtask.TaskStateInit, businessID, flowID, uid, seizecount)
	if err != nil {
		log.Error("db.Query error(%v)", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		if err = rows.Scan(&id); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		hitids = append(hitids, id)
	}
	return
}
