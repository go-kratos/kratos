package dao

import (
	"context"

	"github.com/pkg/errors"

	"database/sql"
	"go-common/app/interface/main/laser/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_queryTaskInfoSql  = "SELECT mid, log_date, source_type, platform FROM task WHERE state = 0 AND is_deleted = 0 AND mid = ?"
	_updateStateSql    = "UPDATE task SET state = ? WHERE is_deleted = 0 AND id = ? "
	_queryTaskIDSql    = "SELECT id FROM task WHERE is_deleted = 0 AND state = 0 and mid = ? "
	_selectByPrimaryID = "SELECT id, admin_id, mid, log_date, contact_email, source_type, platform,state, is_deleted, mtime, ctime FROM task WHERE id = ? AND is_deleted = 0"
)

func (d *Dao) QueryUndoneTaskInfo(c context.Context, mid int64) (t *model.TaskInfo, err error) {
	t = &model.TaskInfo{}
	row := d.db.QueryRow(c, _queryTaskInfoSql, mid)
	if err = row.Scan(&t.MID, &t.LogDate, &t.SourceType, &t.Platform); err != nil {
		if err == sql.ErrNoRows {
			t = nil
			err = nil
		} else {
			err = errors.WithStack(err)
			log.Error("row.Scan() error(%v)", err)
		}
	}
	return
}

func (d *Dao) TxUpdateTaskState(c context.Context, tx *xsql.Tx, state int, taskID int64) (rows int64, err error) {
	res, err := tx.Exec(_updateStateSql, state, taskID)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

func (d *Dao) QueryTaskID(c context.Context, mid int64) (taskID int64, err error) {
	row := d.db.QueryRow(c, _queryTaskIDSql, mid)
	if err = row.Scan(&taskID); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			err = errors.WithStack(err)
			log.Error("row.Scan() error(%v)", err)
		}
	}
	return
}

func (d *Dao) DetailTask(c context.Context, taskID int64) (t *model.Task, err error) {
	t = &model.Task{}
	row := d.db.QueryRow(c, _selectByPrimaryID, taskID)
	if err = row.Scan(&t.ID, &t.AdminID, &t.MID, &t.LogDate, &t.ContactEmail, &t.SourceType, &t.Platform, &t.State, &t.IsDeleted, &t.CTime, &t.MTime); err != nil {
		err = errors.WithStack(err)
		log.Error("rows.Scan error(%v)", err)
		return
	}
	return
}
