package archive

import (
	"context"
	"time"

	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_inDispatchDoneSQL = `INSERT IGNORE INTO task_dispatch_done(task_id,pool,subject,adminid,aid,cid,uid,state,utime,ctime,mtime,dtime,gtime,weight,conf_id,conf_state,conf_weight,upspecial,ptime,uptime,cftime)
		SELECT id,pool,subject,adminid,aid,cid,uid,state,utime,ctime,mtime,dtime,gtime,weight,conf_id,conf_state,conf_weight,upspecial,ptime,uptime,cftime FROM task_dispatch WHERE mtime>=? AND mtime<=? AND state IN (2,6);`
	_delTaskDoneBeforeSQL = "DELETE FROM task_dispatch_done WHERE mtime<=? LIMIT ?"
	_delTaskBeforeSQL     = "DELETE FROM task_dispatch WHERE mtime<=? LIMIT ?"
)

// TxAddDispatchDone add task dispatch done
func (d *Dao) TxAddDispatchDone(c context.Context, tx *sql.Tx, startTime, endTime time.Time) (rows int64, err error) {
	res, err := tx.Exec(_inDispatchDoneSQL, startTime.Format("2006-01-02 15:04:05"), endTime.Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Error("tx.Exec(%s, %s, %s) error(%v)", _inDispatchDoneSQL, startTime.Format("2006-01-02 15:04:05"), endTime.Format("2006-01-02 15:04:05"), err)
		return
	}
	return res.RowsAffected()
}

// DelTaskDoneBefore del task_dispatch_done
func (d *Dao) DelTaskDoneBefore(c context.Context, before time.Time, limit int64) (rows int64, err error) {
	res, err := d.db.Exec(c, _delTaskDoneBeforeSQL, before.Format("2006-01-02 15:04:05"), limit)
	if err != nil {
		log.Error("d.db.Exec(%s, %s, %d) error(%v)", _delTaskDoneBeforeSQL, before, limit, err)
		return
	}
	return res.RowsAffected()
}

// DelTaskBefore del task_dispatch
func (d *Dao) DelTaskBefore(c context.Context, before time.Time, limit int64) (rows int64, err error) {
	res, err := d.db.Exec(c, _delTaskBeforeSQL, before.Format("2006-01-02 15:04:05"), limit)
	if err != nil {
		log.Error("d.db.Exec(%s, %s, %d) error(%v)", _delTaskBeforeSQL, before, limit, err)
		return
	}
	return res.RowsAffected()
}
