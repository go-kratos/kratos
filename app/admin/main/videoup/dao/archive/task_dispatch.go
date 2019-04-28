package archive

import (
	"context"
	xsql "database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go-common/app/admin/main/videoup/model/archive"
	"go-common/app/admin/main/videoup/model/utils"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

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
	_taskDispatchByIDSQL = `SELECT id,subject,aid,cid,uid,state,ctime,utime,mtime,dtime,gtime FROM task_dispatch WHERE id=?`
)

// UserUndoneSpecTask get undone dispatch which belongs to someone.
func (d *Dao) UserUndoneSpecTask(c context.Context, uid int64) (tasks []*archive.Task, err error) {
	rows, err := d.db.Query(c, _userUndoneSpecifiedSQL, uid)
	if err != nil {
		log.Error("d.db.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		t := &archive.Task{}
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
func (d *Dao) GetDispatchTask(c context.Context, uid int64) (tls []*archive.TaskForLog, err error) {
	rows, err := d.rddb.Query(c, _dispatchTaskSQL, uid)
	if err != nil {
		log.Error("d.rddb.Query(%s, %d) error(%v)", _dispatchTaskSQL, uid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		taskLog := &archive.TaskForLog{}
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

	res, err = d.db.Exec(c, sqlstring, uid)
	if err != nil {
		log.Error("d.db.Exec(%s %d %v) error(%v)", sqlstring, uid, err)
		return
	}
	return res.RowsAffected()
}

// GetNextTask 获取一条任务
func (d *Dao) GetNextTask(c context.Context, uid int64) (task *archive.Task, err error) {
	task = new(archive.Task)
	err = d.rddb.QueryRow(c, _getNextTaskSQL, uid).Scan(&task.ID, &task.Pool, &task.Subject, &task.AdminID,
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
		_, err = d.db.Exec(c, _upTaskGtimeSQL, timeNow, task.ID)
		if err != nil {
			log.Error("d.db.Exec(%v,%d) error(%v)", timeNow, task.ID, err)
			return nil, err
		}
		task.GTime = utils.NewFormatTime(timeNow)
	}

	return
}

// TaskByID get task
func (d *Dao) TaskByID(c context.Context, id int64) (task *archive.Task, err error) {
	task = new(archive.Task)
	err = d.rddb.QueryRow(c, _taskByIDSQL, id, id).Scan(&task.ID, &task.Pool, &task.Subject, &task.AdminID,
		&task.Aid, &task.Cid, &task.UID, &task.State, &task.UTime, &task.CTime, &task.MTime, &task.DTime, &task.GTime, &task.PTime, &task.Weight)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Error("db.QueryRow(%d) error(%v)", err)
		return nil, err
	}

	return
}

// ListByCondition 从数据库获取读取任务列表
func (d *Dao) ListByCondition(c context.Context, uid int64, pn, ps int, ltype, leader int8) (tasks []*archive.Task, err error) {
	var task *archive.Task
	tasks = []*archive.Task{}
	if !archive.IsDispatch(ltype) {
		log.Error("ListByCondition listtype(%d) error", ltype)
		return
	}

	listSQL := d.sqlHelper(uid, pn, ps, ltype, leader)
	rows, err := d.rddb.Query(c, listSQL)
	if err != nil {
		log.Error("rddb.Query(%s) error(%v)", listSQL, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		task = &archive.Task{}
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

	if uid != 0 && (ltype != archive.TypeRealTime && leader != 1) { //实时任务或者组长不区分uid
		wherecase = append(wherecase, fmt.Sprintf("uid=%d", uid))
	}
	ordercase = append(ordercase, "weight desc,ctime asc")

	switch ltype {
	case archive.TypeRealTime:
		wherecase = append(wherecase, "state=0")
	case archive.TypeDispatched:
		wherecase = append(wherecase, "state=1 AND subject=0")
		ordercase = append(ordercase, "utime desc")
	case archive.TypeDelay:
		wherecase = append(wherecase, "state=3")
		ordercase = append(ordercase, "dtime asc")
	case archive.TypeSpecial:
		wherecase = append(wherecase, "state=5 AND subject=1")
		ordercase = append(ordercase, "mtime asc")
	case archive.TypeSpecialWait:
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
func (d *Dao) GetWeightDB(c context.Context, ids []int64) (mcases map[int64]*archive.TaskPriority, err error) {
	var (
		rows *sql.Rows
		desc xsql.NullString
	)
	sqlstring := fmt.Sprintf(_getWeightDBSQL, xstr.JoinInts(ids))
	if rows, err = d.db.Query(c, sqlstring); err != nil {
		log.Error("d.db.Query(%s) error(%v)", sqlstring, err)
		return
	}
	defer rows.Close()
	mcases = make(map[int64]*archive.TaskPriority)
	for rows.Next() {
		tp := new(archive.TaskPriority)
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

//TaskDispatchByID task by id
func (d *Dao) TaskDispatchByID(c context.Context, id int64) (tk *archive.Task, err error) {
	tk = &archive.Task{}
	if err = d.rddb.QueryRow(c, _taskDispatchByIDSQL, id).Scan(&tk.ID, &tk.Subject, &tk.Aid, &tk.Cid, &tk.UID, &tk.State, &tk.CTime, &tk.UTime, &tk.MTime, &tk.DTime, &tk.GTime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("TaskDispatchByID rows.Scan error(%v) id(%d)", err, id)
		}
	}

	return
}
