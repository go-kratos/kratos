package dao

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go-common/app/admin/main/dm/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_countDMTask       = "SELECT count(*) FROM dm_task"
	_selectDMTask      = "SELECT id,title,creator,reviewer,state,result,ctime,mtime FROM dm_task%s order by ctime limit ?,? "
	_insertDMTask      = "INSERT INTO dm_task(title,creator,regex,keywords,ips,mids,cids,start,end,state,sub) VALUES(?,?,?,?,?,?,?,?,?,?,?)"
	_insertDMSubTask   = "INSERT INTO dm_sub_task(task_id,operation,start,rate) VALUES(?,?,?,?)"
	_updateDMTaskState = "UPDATE dm_task SET state=? WHERE id IN (%s) AND state!=?"
	_reviewDmTask      = "UPDATE dm_task SET state=?,reviewer=?,topic=? WHERE id=? AND state=0"
	_selectTaskByID    = "SELECT id,title,creator,reviewer,regex,keywords,ips,mids,cids,start,end,qcount,state,result,ctime,mtime FROM dm_task WHERE id=?"
	_selectSubTask     = "SELECT id,operation,rate,tcount,start,end FROM dm_sub_task WHERE task_id=?"
	_editTaskPriority  = "UPDATE dm_task SET priority=? WHERE id IN (%s)"
)

// TaskList dm task list
func (d *Dao) TaskList(c context.Context, taskSQL []string, pn, ps int64) (tasks []*model.TaskInfo, total int64, err error) {
	var sql string
	tasks = make([]*model.TaskInfo, 0)
	if len(taskSQL) > 0 {
		sql = fmt.Sprintf(" WHERE %s", strings.Join(taskSQL, " AND "))
	}
	countRow := d.biliDM.QueryRow(c, _countDMTask+sql)
	if err = countRow.Scan(&total); err != nil {
		log.Error("row.ScanCount(%s) error(%v)", _countDMTask+sql, err)
		return
	}
	rows, err := d.biliDM.Query(c, fmt.Sprintf(_selectDMTask, sql), (pn-1)*ps, ps)
	if err != nil {
		log.Error("biliDM.Query(%s) error(%v)", fmt.Sprintf(_selectDMTask, sql), err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		task := &model.TaskInfo{}
		var cTime, mTime time.Time
		if err = rows.Scan(&task.ID, &task.Title, &task.Creator, &task.Reviewer, &task.State, &task.Result, &cTime, &mTime); err != nil {
			log.Error("biliDM.Scan(%s) error(%v)", fmt.Sprintf(_selectDMTask, sql), err)
			return
		}
		task.Ctime = cTime.Format("2006-01-02 15:04:05")
		task.Mtime = mTime.Format("2006-01-02 15:04:05")
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		log.Error("biliDM.rows.Err() error(%v)", err)
	}
	return
}

// AddTask add dm task
func (d *Dao) AddTask(tx *sql.Tx, v *model.AddTaskArg, sub int32) (taskID int64, err error) {
	var sTime, eTime time.Time
	if sTime, err = time.ParseInLocation("2006-01-02 15:04:05", v.Start, time.Local); err != nil {
		log.Error("d.AddTask time.Parse(%s) error(%v)", v.Start, err)
		return
	}
	if eTime, err = time.ParseInLocation("2006-01-02 15:04:05", v.End, time.Local); err != nil {
		log.Error("d.AddTask time.Parse(%s) error(%v)", v.End, err)
		return
	}
	// regex add slash
	rows, err := tx.Exec(_insertDMTask, v.Title, v.Creator, v.Regex, v.KeyWords, v.IPs, v.Mids, v.Cids, sTime, eTime, v.State, sub)
	if err != nil {
		log.Error("tx.Exec(%s params:%+v) error(%v)", _insertDMTask, v, err)
		return
	}
	return rows.LastInsertId()
}

// AddSubTask add dm sub task
func (d *Dao) AddSubTask(tx *sql.Tx, taskID int64, operation int32, start string, rate int32) (id int64, err error) {
	sTime, err := time.ParseInLocation("2006-01-02 15:04:05", start, time.Local)
	if err != nil {
		log.Error("d.AddSubTask time.Parse(%s) error(%v)", start, err)
		return
	}
	rows, err := tx.Exec(_insertDMSubTask, taskID, operation, sTime, rate)
	if err != nil {
		log.Error("tx.Exec(%s,%d,%d,%s,%d) error(%v)", _insertDMSubTask, taskID, operation, start, rate, err)
		return
	}
	return rows.LastInsertId()
}

// EditTaskState .
func (d *Dao) EditTaskState(c context.Context, v *model.EditTasksStateArg) (affected int64, err error) {
	updateSQL := fmt.Sprintf(_updateDMTaskState, v.IDs)
	rows, err := d.biliDM.Exec(c, updateSQL, v.State, v.State)
	if err != nil {
		log.Error("d.EditTaskState.Exec(id:%s, state:%d) error(%v)", v.IDs, v.State, err)
		return
	}
	return rows.RowsAffected()
}

// EditTaskPriority .
func (d *Dao) EditTaskPriority(c context.Context, ids string, priority int64) (affected int64, err error) {
	updateSQL := fmt.Sprintf(_editTaskPriority, ids)
	rows, err := d.biliDM.Exec(c, updateSQL, priority)
	if err != nil {
		log.Error("d.EditTaskPriority.Exec(ids:%s, priority:%d) error(%v)", ids, priority, err)
		return
	}
	return rows.RowsAffected()
}

// ReviewTask .
func (d *Dao) ReviewTask(c context.Context, v *model.ReviewTaskArg) (affected int64, err error) {
	row, err := d.biliDM.Exec(c, _reviewDmTask, v.State, v.Reviewer, v.Topic, v.ID)
	if err != nil {
		log.Error("d.ReviewTask.Exec(id:%d, state:%d) error(%v)", v.ID, v.State, err)
		return
	}
	return row.RowsAffected()
}

// TaskView .
func (d *Dao) TaskView(c context.Context, id int64) (task *model.TaskView, err error) {
	task = new(model.TaskView)
	row := d.biliDM.QueryRow(c, _selectTaskByID, id)
	var sTime, eTime, cTime, mTime time.Time
	if err = row.Scan(&task.ID, &task.Title, &task.Creator, &task.Reviewer, &task.Regex, &task.KeyWords, &task.IPs, &task.Mids, &task.Cids, &sTime, &eTime, &task.QCount, &task.State, &task.Result, &cTime, &mTime); err != nil {
		if err == sql.ErrNoRows {
			err = ecode.NothingFound
		}
		log.Error("biliDM.Scan(%s, id:%d) error(%v)", _selectTaskByID, id, err)
		return
	}
	task.Start = sTime.Format("2006-01-02 15:04:05")
	task.End = eTime.Format("2006-01-02 15:04:05")
	task.Ctime = cTime.Format("2006-01-02 15:04:05")
	task.Mtime = mTime.Format("2006-01-02 15:04:05")
	return
}

// SubTask .
func (d *Dao) SubTask(c context.Context, id int64) (subTask *model.SubTask, err error) {
	// TODO: operation time
	subTask = new(model.SubTask)
	row := d.biliDM.QueryRow(c, _selectSubTask, id)
	var sTime time.Time
	var eTime time.Time
	if err = row.Scan(&subTask.ID, &subTask.Operation, &subTask.Rate, &subTask.Tcount, &sTime, &eTime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			subTask = nil
		}
		log.Error("biliDM.Scan(%s, taskID:%d) error*(%v)", _selectSubTask, id, err)
		return
	}
	subTask.Start = sTime.Format("2006-01-02 15:04:05")
	subTask.End = eTime.Format("2006-01-02 15:04:05")
	return
}
