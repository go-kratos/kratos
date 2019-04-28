package dao

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-common/library/log"
)

const (
	_topDyKeyFmt = "spc_td_%d"
	_topDySQL    = `SELECT dy_id FROM top_dynamic_%d WHERE deleted_time = 0 AND mid = ?`
	_topDyAddSQL = `INSERT INTO top_dynamic_%d(mid,dy_id) VALUES (?,?) ON DUPLICATE KEY UPDATE dy_id = ?`
	_topDyDelSQL = `UPDATE top_dynamic_%d set deleted_time = ? WHERE mid = ?`
)

func topDyHit(mid int64) int64 {
	return mid % 10
}

func topDyKey(mid int64) string {
	return fmt.Sprintf(_topDyKeyFmt, mid)
}

// RawTopDynamic get top dynamic data from mysql.
func (d *Dao) RawTopDynamic(c context.Context, mid int64) (dyID int64, err error) {
	var row = d.db.QueryRow(c, fmt.Sprintf(_topDySQL, topDyHit(mid)), mid)
	if err = row.Scan(&dyID); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("RawTopDynamic row.Scan() error(%v)", err)
		}
	}
	return
}

// AddTopDynamic add top archive.
func (d *Dao) AddTopDynamic(c context.Context, mid, dyID int64) (err error) {
	if _, err = d.db.Exec(c, fmt.Sprintf(_topDyAddSQL, topDyHit(mid)), mid, dyID, dyID); err != nil {
		log.Error("AddTopDynamic error d.db.Exec(%d,%d) error(%v)", mid, dyID, err)
	}
	return
}

// DelTopDynamic delete top archive.
func (d *Dao) DelTopDynamic(c context.Context, mid int64, now time.Time) (err error) {
	if _, err = d.db.Exec(c, fmt.Sprintf(_topDyDelSQL, topDyHit(mid)), now, mid); err != nil {
		log.Error("DelTopDynamic error d.db.Exec(%d) error(%v)", mid, err)
	}
	return
}
