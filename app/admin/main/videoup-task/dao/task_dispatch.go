package dao

import (
	"context"
	xsql "database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go-common/app/admin/main/videoup-task/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_upTaskByIDSQL     = "UPDATE task_dispatch SET %s WHERE id=?"
	_upGtimeByIDSQL    = "UPDATE task_dispatch SET gtime=? WHERE id=?"
	_releaseByIDSQL    = "UPDATE task_dispatch SET subject=0,state=0,uid=0,gtime='0000-00-00 00:00:00' WHERE id=?"
	_releaseMtimeSQL   = "UPDATE task_dispatch SET subject=0,state=0,uid=0,gtime='0000-00-00 00:00:00' WHERE id IN (%s) AND mtime<=?"
	_timeOutTaskSQL    = "SELECT id,cid,subject,mtime FROM task_dispatch WHERE (state=1 AND mtime<?) OR (state=0 AND uid<>0 AND ctime<?)"
	_getRelTaskSQL     = "SELECT id,cid,subject,mtime,gtime FROM task_dispatch WHERE state IN (0,1) AND uid=?"
	_releaseSpecialSQL = "UPDATE task_dispatch SET subject=0,state=0,uid=0 WHERE id=? AND gtime='0000-00-00 00:00:00' AND mtime<=? AND state=? AND uid=?"
)

// UpGtimeByID update gtime
func (d *Dao) UpGtimeByID(c context.Context, id int64, gtime string) (rows int64, err error) {
	var res xsql.Result

	if res, err = d.arcDB.Exec(c, _upGtimeByIDSQL, gtime, id); err != nil {
		log.Error("d.arcDB.Exec(%s, %v, %d) error(%v)", _upGtimeByIDSQL, gtime, id)
		return
	}
	return res.RowsAffected()
}

// TxUpTaskByID 更新任务状态
func (d *Dao) TxUpTaskByID(tx *sql.Tx, id int64, paras map[string]interface{}) (rows int64, err error) {
	arrSet := []string{}
	arrParas := []interface{}{}
	for k, v := range paras {
		arrSet = append(arrSet, k+"=?")
		arrParas = append(arrParas, v)
	}
	arrParas = append(arrParas, id)
	sqlstring := fmt.Sprintf(_upTaskByIDSQL, strings.Join(arrSet, ","))
	res, err := tx.Exec(sqlstring, arrParas...)
	if err != nil {
		log.Error("tx.Exec(%v %v) error(%v)", sqlstring, arrParas, err)
		return
	}
	return res.RowsAffected()
}

// TxReleaseByID 释放指定任务
func (d *Dao) TxReleaseByID(tx *sql.Tx, id int64) (rows int64, err error) {
	res, err := tx.Exec(_releaseByIDSQL, id)
	if err != nil {
		log.Error("tx.Exec(%s, %d) error(%v)", _releaseByIDSQL, id, err)
		return
	}
	return res.RowsAffected()
}

// MulReleaseMtime 批量释放任务,加时间防止释放错误
func (d *Dao) MulReleaseMtime(c context.Context, ids []int64, mtime time.Time) (rows int64, err error) {
	sqlstring := fmt.Sprintf(_releaseMtimeSQL, xstr.JoinInts(ids))
	res, err := d.arcDB.Exec(c, sqlstring, mtime)
	if err != nil {
		log.Error("tx.Exec(%s,  %v) error(%v)", sqlstring, mtime, err)
		return
	}
	return res.RowsAffected()
}

// GetTimeOutTask 释放正在处理且超时的,释放指派后但长时间未审核的
func (d *Dao) GetTimeOutTask(c context.Context) (rts []*model.TaskForLog, err error) {
	var (
		rows *sql.Rows
	)
	if rows, err = d.arcDB.Query(c, _timeOutTaskSQL, time.Now().Add(-10*time.Minute), time.Now().Add(-80*time.Minute)); err != nil {
		log.Error("d.arcDB.Query(%s) error(%v)", _timeOutTaskSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		rt := &model.TaskForLog{}
		if err = rows.Scan(&rt.ID, &rt.Cid, &rt.Subject, &rt.Mtime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		rts = append(rts, rt)
	}
	return
}

// GetRelTask 用户登出或者主动释放(分配给该用户的都释放)
func (d *Dao) GetRelTask(c context.Context, uid int64) (rts []*model.TaskForLog, lastid int64, err error) {
	var (
		gtime time.Time
		rows  *sql.Rows
	)
	if rows, err = d.arcDB.Query(c, _getRelTaskSQL, uid); err != nil {
		log.Error("d.arcDB.Query(%s, %d) error(%v)", _getRelTaskSQL, uid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		rt := &model.TaskForLog{}
		if err = rows.Scan(&rt.ID, &rt.Cid, &rt.Subject, &rt.Mtime, &gtime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		if gtime.IsZero() {
			rts = append(rts, rt)
		} else {
			lastid = rt.ID
		}
	}
	return
}

// TxReleaseSpecial 延时固定时间释放的任务，需要校验释放时的状态，时间，认领人等
func (d *Dao) TxReleaseSpecial(tx *sql.Tx, mtime time.Time, state int8, taskid, uid int64) (rows int64, err error) {
	res, err := tx.Exec(_releaseSpecialSQL, taskid, mtime, state, uid)
	if err != nil {
		log.Error("tx.Exec(%s, %d, %v, %d, %d) error(%v)", _releaseSpecialSQL, taskid, mtime, state, uid, err)
		return
	}
	return res.RowsAffected()
}

const (
	_userUndoneSpecifiedSQL = "SELECT id,pool,subject,adminid,aid,cid,uid,state,ctime,mtime FROM task_dispatch WHERE uid = ? AND state !=2 AND subject = 1"
	_dispatchTaskSQL        = "SELECT id,cid,mtime FROM task_dispatch WHERE uid in (0,?) AND state = 0 ORDER BY `weight` DESC,`subject` DESC,`id` ASC limit 8"
	_upDispatchTaskSQL      = "UPDATE task_dispatch SET state=1,uid=?,gtime='0000-00-00 00:00:00' WHERE id IN (%s) AND state=0"
	_getNextTaskSQL         = "SELECT id,pool,subject,adminid,aid,cid,uid,state,utime,ctime,mtime,dtime,gtime,weight FROM task_dispatch WHERE uid=? AND state = 1 ORDER BY `weight` DESC,`subject` DESC,`id` ASC limit 1"
	_upTaskGtimeSQL         = "UPDATE task_dispatch SET gtime=? WHERE id=?"
	_listByConditionSQL     = "SELECT id,pool,subject,adminid,aid,cid,uid,state,utime,ctime,mtime,dtime,gtime,weight FROM task_dispatch where %s order by %s %s"
	_taskByIDSQL            = "SELECT id,pool,subject,adminid,aid,cid,uid,state,utime,ctime,mtime,dtime,gtime,ptime,weight FROM task_dispatch WHERE id =? union " +
		"SELECT task_id as id,pool,subject,adminid,aid,cid,uid,state,utime,ctime,mtime,dtime,gtime,ptime,weight FROM task_dispatch_done WHERE task_id=?"
	_getWeightDBSQL = "SELECT t.id,t.state,a.mid,t.ctime,t.upspecial,t.ptime,e.description FROM `task_dispatch` AS t " +
		"LEFT JOIN `task_dispatch_extend` AS e ON t.id=e.task_id INNER JOIN archive as a ON a.id=t.aid WHERE t.id IN (%s)"
)

// UserUndoneSpecTask get undone dispatch which belongs to someone.
func (d *Dao) UserUndoneSpecTask(c context.Context, uid int64) (tasks []*model.Task, err error) {
	rows, err := d.arcDB.Query(c, _userUndoneSpecifiedSQL, uid)
	if err != nil {
		log.Error("d.arcDB.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		t := &model.Task{}
		if err = rows.Scan(&t.ID, &t.Pool, &t.Subject, &t.AdminID, &t.Aid, &t.Cid, &t.UID, &t.State, &t.CTime, &t.MTime); err != nil {
			if err == sql.ErrNoRows {
				err = nil
				return
			}
			log.Error("row.Scan(%d) error(%v)", err)
			return
		}
		tasks = append(tasks, t)
	}
	return
}

// GetDispatchTask 获取抢占到的任务(用于记录日志)
func (d *Dao) GetDispatchTask(c context.Context, uid int64) (tls []*model.TaskForLog, err error) {
	rows, err := d.arcDB.Query(c, _dispatchTaskSQL, uid)
	if err != nil {
		log.Error("d.arcDB.Query(%s, %d) error(%v)", _dispatchTaskSQL, uid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		taskLog := &model.TaskForLog{}
		if err = rows.Scan(&taskLog.ID, &taskLog.Cid, &taskLog.Mtime); err != nil {
			log.Error("rows.Scan(%s, %d) error(%v)", _dispatchTaskSQL, uid, err)
			return
		}
		tls = append(tls, taskLog)
	}
	return
}

// UpDispatchTask 抢占任务
func (d *Dao) UpDispatchTask(c context.Context, uid int64, ids []int64) (rows int64, err error) {
	var (
		res       xsql.Result
		sqlstring = fmt.Sprintf(_upDispatchTaskSQL, xstr.JoinInts(ids))
	)

	res, err = d.arcDB.Exec(c, sqlstring, uid)
	if err != nil {
		log.Error("d.arcDB.Exec(%s %d %v) error(%v)", sqlstring, uid, err)
		return
	}
	return res.RowsAffected()
}

// GetNextTask 获取一条任务
func (d *Dao) GetNextTask(c context.Context, uid int64) (task *model.Task, err error) {
	task = new(model.Task)
	err = d.arcDB.QueryRow(c, _getNextTaskSQL, uid).Scan(&task.ID, &task.Pool, &task.Subject, &task.AdminID,
		&task.Aid, &task.Cid, &task.UID, &task.State, &task.UTime, &task.CTime, &task.MTime, &task.DTime, &task.GTime, &task.Weight)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Error("db.QueryRow(%d) error(%v)", err)
		return nil, err
	}

	if task.GTime.TimeValue().IsZero() {
		timeNow := time.Now()
		_, err = d.arcDB.Exec(c, _upTaskGtimeSQL, timeNow, task.ID)
		if err != nil {
			log.Error("d.arcDB.Exec(%v,%d) error(%v)", timeNow, task.ID, err)
			return nil, err
		}
		task.GTime = model.NewFormatTime(timeNow)
	}

	return
}

// TaskByID get task
func (d *Dao) TaskByID(c context.Context, id int64) (task *model.Task, err error) {
	task = new(model.Task)
	err = d.arcDB.QueryRow(c, _taskByIDSQL, id, id).Scan(&task.ID, &task.Pool, &task.Subject, &task.AdminID,
		&task.Aid, &task.Cid, &task.UID, &task.State, &task.UTime, &task.CTime, &task.MTime, &task.DTime, &task.GTime, &task.PTime, &task.Weight)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
			task = nil
			return
		}
		log.Error("db.QueryRow(%d) error(%v)", id, err)
		return nil, err
	}

	return
}

// ListByCondition 从数据库获取读取任务列表
func (d *Dao) ListByCondition(c context.Context, uid int64, pn, ps int, ltype, leader int8) (tasks []*model.Task, err error) {
	var task *model.Task
	tasks = []*model.Task{}
	if !model.IsDispatch(ltype) {
		log.Error("ListByCondition listtype(%d) error", ltype)
		return
	}

	listSQL := d.sqlHelper(uid, pn, ps, ltype, leader)
	rows, err := d.arcDB.Query(c, listSQL)
	if err != nil {
		log.Error("rddb.Query(%s) error(%v)", listSQL, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		task = &model.Task{}
		err = rows.Scan(&task.ID, &task.Pool, &task.Subject, &task.AdminID,
			&task.Aid, &task.Cid, &task.UID, &task.State, &task.UTime, &task.CTime, &task.MTime, &task.DTime, &task.GTime, &task.Weight)
		if err != nil {
			log.Error("rows.Scan(%s) error(%v)", listSQL, err)
			return nil, nil
		}
		tasks = append(tasks, task)
	}
	return
}

func (d *Dao) sqlHelper(uid int64, pn, ps int, ltype int8, leader int8) string {
	var (
		wherecase []string
		ordercase []string
		limitStr  string
		whereStr  string
		orderStr  string
	)
	limitStr = fmt.Sprintf("LIMIT %d,%d", (pn-1)*ps, ps)

	if uid != 0 && (ltype != model.TypeRealTime && leader != 1) { //实时任务或者组长不区分uid
		wherecase = append(wherecase, fmt.Sprintf("uid=%d", uid))
	}
	ordercase = append(ordercase, "weight desc,ctime asc")

	switch ltype {
	case model.TypeRealTime:
		wherecase = append(wherecase, "state=0")
	case model.TypeDispatched:
		wherecase = append(wherecase, "state=1 AND subject=0")
		ordercase = append(ordercase, "utime desc")
	case model.TypeDelay:
		wherecase = append(wherecase, "state=3")
		ordercase = append(ordercase, "dtime asc")
	case model.TypeReview:
		wherecase = append(wherecase, "state=5")
		ordercase = append(ordercase, "mtime asc")
	case model.TypeSpecialWait:
		wherecase = append(wherecase, "state=1 AND subject=1")
		ordercase = append(ordercase, "utime desc")
	default:
		wherecase = append(wherecase, "state=0")
	}

	whereStr = strings.Join(wherecase, " AND ")
	orderStr = strings.Join(ordercase, ",")

	return fmt.Sprintf(_listByConditionSQL, whereStr, orderStr, limitStr)
}

// GetWeightDB 从数据库读取权重配置
func (d *Dao) GetWeightDB(c context.Context, ids []int64) (mcases map[int64]*model.TaskPriority, err error) {
	var (
		rows *sql.Rows
		desc xsql.NullString
	)
	sqlstring := fmt.Sprintf(_getWeightDBSQL, xstr.JoinInts(ids))
	if rows, err = d.arcDB.Query(c, sqlstring); err != nil {
		log.Error("d.arcDB.Query(%s) error(%v)", sqlstring, err)
		return
	}
	defer rows.Close()
	mcases = make(map[int64]*model.TaskPriority)
	for rows.Next() {
		tp := new(model.TaskPriority)
		if err = rows.Scan(&tp.TaskID, &tp.State, &tp.Mid, &tp.Ctime, &tp.Special, &tp.Ptime, &desc); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		if desc.Valid && len(desc.String) > 0 {
			if err = json.Unmarshal([]byte(desc.String), &(tp.CfItems)); err != nil {
				log.Error("json.Unmarshal error(%v)", err)
				return
			}
		}
		mcases[tp.TaskID] = tp
	}
	return
}
