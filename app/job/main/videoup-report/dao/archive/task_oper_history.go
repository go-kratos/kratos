package archive

import (
	"context"
	"time"

	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_inTaskHisSQL          = "INSERT INTO task_oper_history(pool,action,task_id,cid,uid,result,reason) VALUE (?,?,?,?,?,?,?);"
	_mvTaskOperHisSQL      = "INSERT IGNORE INTO task_oper_history_done(id,pool,action,task_id,cid,uid,result,reason,utime,ctime,mtime) SELECT id,pool,action,task_id,cid,uid,result,reason,utime,ctime,mtime FROM task_oper_history WHERE mtime < ? LIMIT ?"
	_delTaskOperHisSQL     = "DELETE FROM task_oper_history WHERE mtime < ? LIMIT ?"
	_delTaskHistoryDoneSQL = "DELETE FROM task_oper_history_done WHERE mtime < ? LIMIT ?"
)

// AddTaskHis add task oper history
func (d *Dao) AddTaskHis(c context.Context, pool int8, action int8, taskID int64, cid int64, uid int64, result int16, reason string) (rows int64, err error) {
	res, err := d.db.Exec(c, _inTaskHisSQL, pool, action, taskID, cid, uid, result, reason)
	if err != nil {
		log.Error("d.db.Exec(%s) error(%v)", _inTaskHisSQL, err)
		return
	}
	return res.RowsAffected()
}

// TxMoveTaskOperDone select into from task_oper_history to task_oper_history_done before t.
func (d *Dao) TxMoveTaskOperDone(tx *sql.Tx, t time.Time, limit int64) (rows int64, err error) {
	res, err := tx.Exec(_mvTaskOperHisSQL, t, limit)
	if err != nil {
		log.Error("tx.Exec(%s, %s, %d) error(%v)", _mvTaskOperHisSQL, t, limit, err)
		return
	}
	return res.RowsAffected()
}

// TxDelTaskOper delete from task_oper_history before t.
func (d *Dao) TxDelTaskOper(tx *sql.Tx, t time.Time, limit int64) (rows int64, err error) {
	res, err := tx.Exec(_delTaskOperHisSQL, t, limit)
	if err != nil {
		log.Error("tx.Exec(%s, %s, %d) error(%v)", _delTaskOperHisSQL, t, limit, err)
		return
	}
	return res.RowsAffected()
}

// DelTaskHistoryDone del oper done
func (d *Dao) DelTaskHistoryDone(c context.Context, before time.Time, limit int64) (rows int64, err error) {
	res, err := d.db.Exec(c, _delTaskHistoryDoneSQL, before.Format("2006-01-02 15:04:05"), limit)
	if err != nil {
		log.Error("d.db.Exec(%s, %s, %d) error(%v)", _delTaskHistoryDoneSQL, before, limit, err)
		return
	}
	return res.RowsAffected()
}
