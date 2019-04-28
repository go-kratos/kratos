package dao

import (
	"context"
	"database/sql"
	"fmt"

	"go-common/library/log"
)

const (
	_taskStatus   = "SELECT status FROM task_status WHERE date=? AND type=?"
	_upTaskStatus = "UPDATE task_status SET status=? WHERE date=? AND type=?"

	_inTaskStatus = "INSERT INTO task_status(type, date, status, message) VALUES (%d, '%s', %d, '%s') ON DUPLICATE KEY UPDATE status=VALUES(status), message=VALUES(message)"
)

// TaskStatus get task status
func (d *Dao) TaskStatus(c context.Context, date string, typ int) (status int, err error) {
	row := d.db.QueryRow(c, _taskStatus, date, typ)
	if err = row.Scan(&status); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("dao.GetTaskStatus error(%v)", err)
		}
	}
	return
}

// UpdateTaskStatus update task status
func (d *Dao) UpdateTaskStatus(c context.Context, date string, typ int, status int) (rows int64, err error) {
	res, err := d.db.Exec(c, _upTaskStatus, status, date, typ)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// InsertTaskStatus insert task status
func (d *Dao) InsertTaskStatus(c context.Context, typ, status int, date, message string) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_inTaskStatus, typ, date, status, message))
	if err != nil {
		return
	}
	return res.RowsAffected()
}
