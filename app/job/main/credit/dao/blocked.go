package dao

import (
	"context"
	"time"

	"go-common/library/log"
)

const (
	_countBlockedSQL = "SELECT COUNT(*) FROM blocked_info WHERE uid=? AND ctime > ? AND status = 0"
	_blockedInfoID   = "SELECT id FROM blocked_info WHERE uid=? ORDER BY id DESC"
)

// CountBlocked get user block count ofter ts.
func (d *Dao) CountBlocked(c context.Context, uid int64, ts time.Time) (count int64, err error) {
	row := d.db.QueryRow(c, _countBlockedSQL, uid, ts)
	if err = row.Scan(&count); err != nil {
		log.Error("d.CountBlocked err(%v)", err)
	}
	return
}

// BlockedInfoID get user blocked new info.
func (d *Dao) BlockedInfoID(c context.Context, uid int64) (id int64, err error) {
	row := d.db.QueryRow(c, _blockedInfoID, uid)
	if err = row.Scan(&id); err != nil {
		log.Error("d.BlockedInfoID err(%v)", err)
	}
	return
}
