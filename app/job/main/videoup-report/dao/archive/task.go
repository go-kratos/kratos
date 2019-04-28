package archive

import (
	"context"
	"time"

	"go-common/app/job/main/videoup-report/model/archive"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	// task
	_taskByMtimeSQL     = "SELECT id,state,ctime,mtime FROM task_dispatch WHERE mtime>? and ptime=0"
	_taskDoneByMtimeSQL = "SELECT id,state,ctime,mtime FROM task_dispatch_done WHERE mtime>? and ptime=0"
	_taskByUntreatedSQL = "SELECT id,state,ctime,mtime FROM task_dispatch WHERE (state=0 OR state=1) and ptime=0"
	// task took in and sel
	_addTaskTookSQL         = "INSERT INTO task_dispatch_took(m50,m60,m80,m90,type,ctime,mtime) VALUE(?,?,?,?,?,?,?)"
	_taskTooksSQL           = "SELECT id,m50,m60,m80,m90,type,ctime,mtime FROM task_dispatch_took WHERE type=1 AND ctime>?"
	_taskTookByHalfHourSQL  = "SELECT id,m50,m60,m80,m90,type,ctime,mtime FROM task_dispatch_took WHERE type=2 ORDER BY ctime DESC LIMIT 1"
	_taskTooksByHalfHourSQL = "SELECT id,m50,m60,m80,m90,type,ctime,mtime FROM task_dispatch_took WHERE type=2 AND ctime>=? AND ctime<=? ORDER BY ctime ASC"
)

// TaskByMtime gets to took the task by mtime
func (d *Dao) TaskByMtime(c context.Context, stime time.Time) (tasks []*archive.Task, err error) {
	rows, err := d.db.Query(c, _taskByMtimeSQL, stime)
	if err != nil {
		log.Error("d.taskStmt.Query(%v) error(%v)", stime, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		task := &archive.Task{}
		if err = rows.Scan(&task.ID, &task.State, &task.Ctime, &task.Mtime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		tasks = append(tasks, task)
	}
	return
}

// TaskDoneByMtime gets to took the task done by mtime
func (d *Dao) TaskDoneByMtime(c context.Context, stime time.Time) (tasks []*archive.Task, err error) {
	rows, err := d.db.Query(c, _taskDoneByMtimeSQL, stime)
	if err != nil {
		log.Error("d.taskStmt.Query(%v) error(%v)", stime, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		task := &archive.Task{}
		if err = rows.Scan(&task.ID, &task.State, &task.Ctime, &task.Mtime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		tasks = append(tasks, task)
	}
	return
}

// TaskByUntreated gets to took the task by untreated
func (d *Dao) TaskByUntreated(c context.Context) (tasks []*archive.Task, err error) {
	rows, err := d.db.Query(c, _taskByUntreatedSQL)
	if err != nil {
		log.Error("d.taskStmt.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		task := &archive.Task{}
		if err = rows.Scan(&task.ID, &task.State, &task.Ctime, &task.Mtime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		tasks = append(tasks, task)
	}
	return
}

// AddTaskTook add TaskTook
func (d *Dao) AddTaskTook(c context.Context, took *archive.TaskTook) (lastID int64, err error) {
	res, err := d.db.Exec(c, _addTaskTookSQL, took.M50, took.M60, took.M80, took.M90, took.TypeID, took.Ctime, took.Mtime)
	if err != nil {
		log.Error("d.TaskTookAddStmt.Exec error(%v)", err)
		return
	}
	lastID, err = res.LastInsertId()
	return
}

// TaskTooks gets TaskTook by ctime
func (d *Dao) TaskTooks(c context.Context, stime time.Time) (tooks []*archive.TaskTook, err error) {
	rows, err := d.db.Query(c, _taskTooksSQL, stime)
	if err != nil {
		log.Error("d.TaskTookStmt.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		took := &archive.TaskTook{}
		if err = rows.Scan(&took.ID, &took.M50, &took.M60, &took.M80, &took.M90, &took.TypeID, &took.Ctime, &took.Mtime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		tooks = append(tooks, took)
	}
	return
}

// TaskTookByHalfHour get TaskTook by half hour
func (d *Dao) TaskTookByHalfHour(c context.Context) (took *archive.TaskTook, err error) {
	row := d.db.QueryRow(c, _taskTookByHalfHourSQL)
	took = &archive.TaskTook{}
	if err = row.Scan(&took.ID, &took.M50, &took.M60, &took.M80, &took.M90, &took.TypeID, &took.Ctime, &took.Mtime); err != nil {
		if err == sql.ErrNoRows {
			took = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// TaskTooksByHalfHour get TaskTooks by half hour
func (d *Dao) TaskTooksByHalfHour(c context.Context, stime time.Time, etime time.Time) (tooks []*archive.TaskTook, err error) {
	rows, err := d.db.Query(c, _taskTooksByHalfHourSQL, stime, etime)
	if err != nil {
		log.Error("d.TaskTooksByHalfHour.Query(%v,%v) error(%v)", stime, etime, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		took := &archive.TaskTook{}
		if err = rows.Scan(&took.ID, &took.M50, &took.M60, &took.M80, &took.M90, &took.TypeID, &took.Ctime, &took.Mtime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		tooks = append(tooks, took)
	}
	return
}
