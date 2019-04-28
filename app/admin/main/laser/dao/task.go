package dao

import (
	"go-common/library/log"

	"context"
	"fmt"
	"github.com/pkg/errors"
	"go-common/app/admin/main/laser/model"
	"go-common/library/database/sql"
)

const (
	_findMIDTaskSQL       = " SELECT id, admin_id, username, mid, log_date, contact_email, source_type, platform, state, is_deleted, ctime, mtime  FROM task WHERE mid = ? AND state = ? and is_deleted = 0"
	_queryTaskInfoByIDSQL = " SELECT mid, log_date, source_type, platform FROM task WHERE state = 0 AND is_deleted = 0 AND id = ? "
	_insertTaskSQL        = " INSERT INTO task (mid, admin_id, username, log_date, contact_email, platform, source_type) VALUES (?, ?, ?, ?, ?, ?, ?) "
	_deleteTaskSQL        = " UPDATE task SET is_deleted = 1 , username = ? , admin_id = ? WHERE id = ? AND is_deleted = 0 "
	_countTaskSQL         = " SELECT count(*) FROM task WHERE %s "
	_queryTaskSQL         = " SELECT id, admin_id, username, mid, log_date, contact_email, source_type, platform, state, is_deleted, ctime, mtime FROM task WHERE %s ORDER BY %s LIMIT %d,%d "
	_updateTaskSQL        = " UPDATE task SET %s WHERE id = ? AND is_deleted = 0 AND state = ? "
)

// AddTask is add a unique task by mid and state(0).
func (d *Dao) AddTask(ctx context.Context, mid int64, username string, adminID int64, logDate string, contactEmail string, platform int, sourceType int) (lastInsertID int64, err error) {
	res, err := d.laserDB.Exec(ctx, _insertTaskSQL, mid, adminID, username, logDate, contactEmail, platform, sourceType)
	if err != nil {
		log.Error("d.AddTask() error(%v)", err)
	}
	return res.LastInsertId()
}

// FindTask is find task by mid and state.
func (d *Dao) FindTask(context context.Context, mid int64, state int) (t *model.Task, err error) {
	t = &model.Task{}
	row := d.laserDB.QueryRow(context, _findMIDTaskSQL, mid, state)
	if err = row.Scan(&t.ID, &t.AdminID, &t.Username, &t.MID, &t.LogDate, &t.ContactEmail, &t.SourceType, &t.Platform, &t.State, &t.IsDeleted, &t.CTime, &t.MTime); err != nil {
		if err == sql.ErrNoRows {
			t = nil
			err = nil
			return
		}
		log.Error("row.Scan() error(%v)", err)
	}
	return
}

// QueryTaskInfoByIDSQL is query task by task id.
func (d *Dao) QueryTaskInfoByIDSQL(c context.Context, id int64) (t *model.TaskInfo, err error) {
	t = &model.TaskInfo{}
	row := d.laserDB.QueryRow(c, _queryTaskInfoByIDSQL, id)
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

// DeleteTask is delete task by TaskID.
func (d *Dao) DeleteTask(ctx context.Context, taskID int64, username string, adminID int64) (err error) {
	_, err = d.laserDB.Exec(ctx, _deleteTaskSQL, username, adminID, taskID)
	if err != nil {
		log.Error("d.DeleteTask() error(%v)", err)
	}
	return
}

// UpdateTask is update undone task where state = 0.
func (d *Dao) UpdateTask(ctx context.Context, taskID int64, state int, updateStmt string) (err error) {
	_, err = d.laserDB.Exec(ctx, fmt.Sprintf(_updateTaskSQL, updateStmt), taskID, state)
	if err != nil {
		log.Error("d.UpdateTask() error(%v)", err)

	}
	return
}

// QueryTask is query task by condition.
func (d *Dao) QueryTask(ctx context.Context, queryStmt string, sort string, offset int, limit int) (tasks []*model.Task, count int64, err error) {
	row := d.laserDB.QueryRow(ctx, fmt.Sprintf(_countTaskSQL, queryStmt))
	if err = row.Scan(&count); err != nil {
		return
	}
	rows, err := d.laserDB.Query(ctx, fmt.Sprintf(_queryTaskSQL, queryStmt, sort, offset, limit))
	if err != nil {
		log.Error("d.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		item := &model.Task{}
		if err = rows.Scan(&item.ID, &item.AdminID, &item.Username, &item.MID, &item.LogDate, &item.ContactEmail, &item.SourceType, &item.Platform, &item.State, &item.IsDeleted, &item.CTime, &item.MTime); err != nil {
			return
		}
		tasks = append(tasks, item)
	}
	return
}
