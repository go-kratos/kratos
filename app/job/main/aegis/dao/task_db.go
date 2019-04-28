package dao

import (
	"context"
	"database/sql"
	"time"

	"go-common/app/job/main/aegis/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_queryTaskSQL = "SELECT id,business_id,flow_id,rid,admin_id,uid,state,weight,utime,gtime,mid,fans,`group`,reason,ctime,mtime FROM task WHERE state=? AND mtime<=? AND id>? ORDER BY id LIMIT ?"

	_upSetWeightSQL = "UPDATE task SET weight=? WHERE id=?"

	_assignTaskSQL = "UPDATE task SET admin_id=?,uid=? WHERE id=? AND state=?"

	_checkTaskSQL = "SELECT id FROM task WHERE flow_id=? AND rid=? AND state<?"
)

// CheckTask 某资源存在未完成任务，不重复添加
func (d *Dao) CheckTask(c context.Context, flowid, rid int64) (id int64) {
	if err := d.fastdb.QueryRow(c, _checkTaskSQL, flowid, rid, model.TaskStateSubmit).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("d.db.QueryRow error(%v)", err)
		}
	}
	return
}

// AssignTask .
func (d *Dao) AssignTask(c context.Context, task *model.Task) (rows int64, err error) {
	var res sql.Result
	if res, err = d.fastdb.Exec(c, _assignTaskSQL, task.AdminID, task.UID, task.ID, model.TaskStateInit); err != nil {
		log.Error("d.db.Exec error(%v)", errors.WithStack(err))
		return
	}
	return res.RowsAffected()
}

// QueryTask .
func (d *Dao) QueryTask(c context.Context, state int8, mtime time.Time, id, limit int64) (tasks []*model.Task, lastid int64, err error) {
	var rows *xsql.Rows
	rows, err = d.slowdb.Query(c, _queryTaskSQL, state, mtime, id, limit)
	if err != nil {
		log.Error("db.Query error(%v)", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		task := &model.Task{}
		if err = rows.Scan(&task.ID, &task.BusinessID, &task.FlowID, &task.RID, &task.AdminID, &task.UID, &task.State, &task.Weight,
			&task.Utime, &task.Gtime, &task.MID, &task.Fans, &task.Group, &task.Reason, &task.Ctime, &task.Mtime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}

		tasks = append(tasks, task)
		lastid = task.ID
	}
	return
}

// SetWeightDB .
func (d *Dao) SetWeightDB(c context.Context, taskid, weight int64) (rows int64, err error) {
	res, err := d.fastdb.Exec(c, _upSetWeightSQL, weight, taskid)
	if err != nil {
		log.Error("db.Exec error(%v)", err)
		return
	}
	return res.LastInsertId()
}
