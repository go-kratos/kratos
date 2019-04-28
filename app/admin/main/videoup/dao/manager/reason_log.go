package manager

import (
	"context"
	"go-common/library/log"
	"time"
)

const (
	_inLogSQL = "INSERT INTO reason_log (oid, type, category_id, reason_id, uid, typeid, ctime, mtime) VALUES (?,?,?,?,?,?,?,?)"
)

// AddReasonLog add a reason log
func (d *Dao) AddReasonLog(c context.Context, oid int64, tp int8, cateID int64, rid int64, uid int64, tid int16, ctime, mtime time.Time) (rows int64, err error) {
	res, err := d.managerDB.Exec(c, _inLogSQL, oid, tp, cateID, rid, uid, tid, ctime, mtime)
	if err != nil {
		log.Error("d.AddReasonLog.Exec error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}
