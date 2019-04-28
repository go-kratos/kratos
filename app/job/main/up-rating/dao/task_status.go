package dao

import (
	"context"
	"fmt"
)

const (
	// insert
	_inTaskStatusSQL = "INSERT INTO task_status(type, status, date, message) VALUES %s ON DUPLICATE KEY UPDATE status=VALUES(status), message=VALUES(message)"
)

// InsertTaskStatus insert task status
func (d *Dao) InsertTaskStatus(c context.Context, val string) (rows int64, err error) {
	if val == "" {
		return
	}
	res, err := d.db.Exec(c, fmt.Sprintf(_inTaskStatusSQL, val))
	if err != nil {
		return
	}
	return res.RowsAffected()
}
