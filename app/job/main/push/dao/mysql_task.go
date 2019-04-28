package dao

import (
	"context"
	"database/sql"
	"encoding/json"
	"strconv"
	"time"

	pamdl "go-common/app/admin/main/push/model"
	pushmdl "go-common/app/service/main/push/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_addTaskSQL          = "INSERT INTO push_tasks (job,type,app_id,business_id,platform,platform_id,title,summary,link_type,link_value,build,sound,vibration,pass_through,mid_file,push_time,expire_time,status,`group`,image_url,extra) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	_delTasksSQL         = `DELETE FROM push_tasks where mtime <= ? limit ?`
	_upadteTaskStatusSQL = "UPDATE push_tasks SET status=? WHERE id=?"
	_taskByStatusSQL     = "SELECT id,job,type,app_id,business_id,platform,title,summary,link_type,link_value,build,sound,vibration,pass_through,mid_file,progress,push_time,expire_time,status,`group`,image_url,extra FROM push_tasks WHERE status=? AND dtime=0 LIMIT 1 FOR UPDATE"
	_upadteTaskSQL       = "UPDATE push_tasks SET mid_file=?,status=? WHERE id=?"
	// dataplatform
	_txDpCondByStatusSQL   = `SELECT id,job,task,conditions,sql_stmt,status,status_url,file FROM push_dataplatform_conditions WHERE status=? LIMIT 1 FOR UPDATE`
	_updateDpCondSQL       = `UPDATE push_dataplatform_conditions SET job=?,task=?,conditions=?,sql_stmt=?,status=?,status_url=?,file=? WHERE id=?`
	_UpdateDpCondStatusSQL = `UPDATE push_dataplatform_conditions SET status=? WHERE id=?`
)

// DelTasks deletes tasks.
func (d *Dao) DelTasks(c context.Context, t time.Time, limit int) (rows int64, err error) {
	res, err := d.delTasksStmt.Exec(c, t, limit)
	if err != nil {
		log.Error("d.DelTasks(%v) error(%v)", t, err)
		PromError("mysql:DelTasks")
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxTaskByStatus gets task by status by tx.
func (d *Dao) TxTaskByStatus(tx *xsql.Tx, status int8) (t *pushmdl.Task, err error) {
	var (
		id       int64
		platform string
		build    string
		progress string
		extra    string
		now      = time.Now()
	)
	t = &pushmdl.Task{Progress: &pushmdl.Progress{}, Extra: &pushmdl.TaskExtra{}}
	if err = tx.QueryRow(_taskByStatusSQL, status).Scan(&id, &t.Job, &t.Type, &t.APPID, &t.BusinessID, &platform, &t.Title, &t.Summary, &t.LinkType, &t.LinkValue, &build,
		&t.Sound, &t.Vibration, &t.PassThrough, &t.MidFile, &progress, &t.PushTime, &t.ExpireTime, &t.Status, &t.Group, &t.ImageURL, &extra); err != nil {
		t = nil
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Error("d.TxTaskByStatus() QueryRow(%d,%v) error(%v)", status, now, err)
		PromError("mysql:按状态查询任务")
		return
	}
	t.ID = strconv.FormatInt(id, 10)
	t.Platform = pushmdl.SplitInts(platform)
	t.Build = pushmdl.ParseBuild(build)
	if progress != "" {
		if err = json.Unmarshal([]byte(progress), t.Progress); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", progress, err)
			return
		}
	}
	if extra != "" {
		if err = json.Unmarshal([]byte(extra), t.Extra); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", extra, err)
		}
	}
	return
}

// TxUpdateTaskStatus updates task status by tx.
func (d *Dao) TxUpdateTaskStatus(tx *xsql.Tx, taskID string, status int8) (err error) {
	id, _ := strconv.ParseInt(taskID, 10, 64)
	if _, err = tx.Exec(_upadteTaskStatusSQL, status, id); err != nil {
		log.Error("d.TxUpdateTaskStatus() Exec(%s,%d) error(%v)", taskID, status, err)
		PromError("mysql:更新推送任务状态")
	}
	return
}

// UpdateTaskStatus update task status.
func (d *Dao) UpdateTaskStatus(c context.Context, taskID int64, status int8) (err error) {
	if _, err = d.updateTaskStatusStmt.Exec(c, status, taskID); err != nil {
		log.Error("d.updateTaskStatusStmt.Exec(%d,%d) error(%v)", taskID, status, err)
		PromError("mysql:更新推送任务状态")
	}
	return
}

// UpdateTask update task.
func (d *Dao) UpdateTask(c context.Context, taskID string, file string, status int8) (err error) {
	id, _ := strconv.ParseInt(taskID, 10, 64)
	if _, err = d.updateTaskStmt.Exec(c, file, status, id); err != nil {
		log.Error("d.updateTaskFileStmt.Exec(%d,%s,%d) error(%v)", id, file, status, err)
		PromError("mysql:更新推送任务file")
	}
	return
}

// AddTask adds task
func (d *Dao) AddTask(ctx context.Context, t *pushmdl.Task) (err error) {
	var (
		platform = pushmdl.JoinInts(t.Platform)
		build, _ = json.Marshal(t.Build)
		extra, _ = json.Marshal(t.Extra)
	)
	if _, err = d.db.Exec(ctx, _addTaskSQL, t.Job, t.Type, t.APPID, t.BusinessID, platform, t.PlatformID, t.Title, t.Summary, t.LinkType, t.LinkValue,
		build, t.Sound, t.Vibration, t.PassThrough, t.MidFile, t.PushTime, t.ExpireTime, t.Status, t.Group, t.ImageURL, extra); err != nil {
		log.Error("d.AddTask(%+v) error(%v)", t, err)
	}
	return
}

// TxCondByStatus gets condition by status.
func (d *Dao) TxCondByStatus(tx *xsql.Tx, status int) (cond *pamdl.DPCondition, err error) {
	cond = new(pamdl.DPCondition)
	if err = tx.QueryRow(_txDpCondByStatusSQL, status).Scan(&cond.ID, &cond.Job, &cond.Task, &cond.Condition, &cond.SQL, &cond.Status, &cond.StatusURL, &cond.File); err != nil {
		if err == sql.ErrNoRows {
			cond = nil
			err = nil
		}
		return
	}
	return
}

// UpdateDpCond update data platform query condition
func (d *Dao) UpdateDpCond(ctx context.Context, cond *pamdl.DPCondition) (err error) {
	if _, err = d.updateDpCondStmt.Exec(ctx, cond.Job, cond.Task, cond.Condition, cond.SQL, cond.Status, cond.StatusURL, cond.File, cond.ID); err != nil {
		log.Error("d.UpdateDpCond(%+v) error(%v)", cond, err)
	}
	return
}

// UpdateDpCondStatus update data platform query condition
func (d *Dao) UpdateDpCondStatus(ctx context.Context, id int64, status int) (err error) {
	if _, err = d.db.Exec(ctx, _UpdateDpCondStatusSQL, status, id); err != nil {
		log.Error("d.UpdateCondStatus(%d,%d) error(%v)", id, status, err)
	}
	return
}

// TxUpdateCondStatus update data platform query condition status
func (d *Dao) TxUpdateCondStatus(tx *xsql.Tx, id int64, status int) (err error) {
	if _, err = tx.Exec(_UpdateDpCondStatusSQL, status, id); err != nil {
		log.Error("d.TxUpdateCondStatus(%d,%d) error(%v)", id, status, err)
	}
	return
}
