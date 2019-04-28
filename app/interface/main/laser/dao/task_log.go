package dao

import (
	"context"

	"github.com/pkg/errors"
	xsql "go-common/library/database/sql"
)

const (
	_addTaskLogSql = "INSERT INTO task_log (task_id, mid, build, platform, task_state, reason) VALUES ( ?, ?, ?, ?, ?, ? )"
)

func (d *Dao) TxAddTaskLog(c context.Context, tx *xsql.Tx, taskID int64, mid int64, build string, platform int, taskState int, reason string) (insertID int64, err error) {
	res, err := tx.Exec(_addTaskLogSql, taskID, mid, build, platform, taskState, reason)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	if insertID, err = res.LastInsertId(); err == nil {
		insertID = int64(insertID)
	}
	return
}
