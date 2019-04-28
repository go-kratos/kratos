package archive

import (
	"context"
	"time"

	"go-common/app/job/main/videoup-report/model/task"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_dispatchSQL          = "SELECT id,state FROM task_dispatch WHERE aid=? AND cid=? ORDER BY id DESC"
	_inDispatchSQL        = "INSERT INTO task_dispatch(pool,subject,adminid,aid,cid,uid,state,conf_id,conf_state,conf_weight,upspecial,cftime,ptime) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?)"
	_delDispatchSQL       = "UPDATE task_dispatch SET state=? WHERE aid=? AND cid=? AND state!=?"
	_delDispatchByAidSQL  = "UPDATE task_dispatch SET state=? WHERE aid=? AND state IN (?,?,?)"
	_delDispatchByTimeSQL = "DELETE FROM task_dispatch WHERE mtime>=? AND mtime<=? AND state in (2,6)"
	_taskIDforWeightSQL   = "SELECT id FROM task_dispatch WHERE state=0 AND id>? ORDER BY id ASC limit 1000"
	_upTaskWeightSQL      = "UPDATE task_dispatch set weight=?,uptime=now() where id=? and state=0"
	_upSpecialSQL         = "UPDATE task_dispatch SET upspecial=? WHERE id=?"
)

// DispatchState get dipatch state.
func (d *Dao) DispatchState(c context.Context, aid, cid int64) (id int64, state int8, err error) {
	row := d.db.QueryRow(c, _dispatchSQL, aid, cid)
	if err = row.Scan(&id, &state); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Error("row.Scan(%d) error(%v)", err)
		return
	}
	return
}

// AddDispatch add task dispatch
func (d *Dao) AddDispatch(c context.Context, t *task.Task) (lastID int64, err error) {
	res, err := d.db.Exec(c, _inDispatchSQL, t.Pool, t.Subject, t.AdminID, t.Aid, t.Cid, t.UID, t.State,
		t.ConfigID, t.ConfigState, t.ConfigWeight, t.UPSpecial, t.CFtime, t.Ptime)
	if err != nil {
		log.Error("d.db.Exec(%s) error(%v)", _inDispatchSQL, err)
		return
	}
	return res.LastInsertId()
}

// DelDispatch del dispatch.
func (d *Dao) DelDispatch(c context.Context, aid, cid int64) (rows int64, err error) {
	res, err := d.db.Exec(c, _delDispatchSQL, task.StateForTaskUserDeleted, aid, cid, task.StateForTaskUserDeleted)
	if err != nil {
		log.Error("d.db.Exec(%s, %d, %d) error(%v)", _delDispatchSQL, aid, cid, err)
		return
	}
	return res.RowsAffected()
}

// DelDispatchByAid del dispatch by aid.
func (d *Dao) DelDispatchByAid(c context.Context, aid int64) (rows int64, err error) {
	res, err := d.db.Exec(c, _delDispatchByAidSQL, task.StateForTaskUserDeleted, aid, task.StateForTaskDefault, task.StateForTaskWork, task.StateForTaskDelay)
	if err != nil {
		log.Error("d.db.Exec(%s, %d, %d) error(%v)", _delDispatchByAidSQL, aid, err)
		return
	}
	return res.RowsAffected()
}

// TxDelDispatchByTime del dispatch by time segment
func (d *Dao) TxDelDispatchByTime(c context.Context, tx *sql.Tx, startTime, endTime time.Time) (rows int64, err error) {
	res, err := tx.Exec(_delDispatchByTimeSQL, startTime.Format("2006-01-02 15:04:05"), endTime.Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Error("tx.Exec(%s, %s, %s) error(%v)", _delDispatchByTimeSQL, startTime.Format("2006-01-02 15:04:05"), endTime.Format("2006-01-02 15:04:05"), err)
		return
	}
	return res.RowsAffected()
}

// TaskIDforWeight 获取需要更新权重的任务id(用于给redis批量读取)
func (d *Dao) TaskIDforWeight(c context.Context, lastid int64) (ids []int64, last int64, err error) {
	rows, err := d.db.Query(c, _taskIDforWeightSQL, lastid) //获取一批待审核任务
	if err != nil {
		log.Error("d.db.Query(%s, %d) error(%v)", _taskIDforWeightSQL, lastid, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		if err = rows.Scan(&id); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		ids = append(ids, id)
		last = id
	}
	return
}

// UpTaskWeight 更新单条权重
func (d *Dao) UpTaskWeight(c context.Context, taskid int64, weight int64) (rows int64, err error) {
	res, err := d.db.Exec(c, _upTaskWeightSQL, weight, taskid)
	if err != nil {
		log.Error("d.db.Exec(%s,%d,%d) error(%v)", _upTaskWeightSQL, weight, taskid, err)
		return
	}
	return res.RowsAffected()
}

// SetUpSpecial 更新单条权重
func (d *Dao) SetUpSpecial(c context.Context, taskid int64, special int8) (rows int64, err error) {
	res, err := d.db.Exec(c, _upSpecialSQL, special, taskid)
	if err != nil {
		log.Error("d.db.Exec(%s,%d,%d) error(%v)", _upSpecialSQL, special, taskid, err)
		return
	}
	return res.RowsAffected()
}
