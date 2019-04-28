package dao

import (
	"context"
	"fmt"
	"go-common/app/admin/main/laser/model"
)

const (
	_countTaskLogSQL = "SELECT count(*) FROM task_log %s"
	_queryTaskLogSQL = "SELECT * FROM task_log %s ORDER BY %s LIMIT %d, %d"
)

// QueryTaskLog is query finished task.
func (d *Dao) QueryTaskLog(ctx context.Context, queryStmt string, sort string, offset int, limit int) (taskLogs []*model.TaskLog, count int64, err error) {
	row := d.laserDB.QueryRow(ctx, fmt.Sprintf(_countTaskLogSQL, queryStmt))
	if err = row.Scan(&count); err != nil {
		return
	}
	rows, err := d.laserDB.Query(ctx, fmt.Sprintf(_queryTaskLogSQL, queryStmt, sort, offset, limit))
	if err != nil {
		return
	}
	for rows.Next() {
		t := &model.TaskLog{}
		if err = rows.Scan(&t.ID, &t.TaskID, &t.MID, &t.Build, &t.Platform, &t.TaskState, &t.Reason, &t.CTime, &t.MTime); err != nil {
			return
		}
		taskLogs = append(taskLogs, t)
	}
	return
}
