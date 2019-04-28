package archive

import (
	"context"
	"fmt"
	"strings"

	"go-common/app/admin/main/videoup/model/archive"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_inTaskHisSQL    = "INSERT INTO task_oper_history(pool,action,task_id,cid,uid,result,reason,utime) VALUE (?,?,?,?,?,?,?,?)"
	_mulinTaskHisSQL = "INSERT INTO task_oper_history(action,task_id,cid,uid) VALUES "
)

// TxAddTaskHis add task oper history
func (d *Dao) TxAddTaskHis(tx *sql.Tx, pool, action int8, taskID, cid, uid, utime int64, result int16, reason string) (rows int64, err error) {
	res, err := tx.Exec(_inTaskHisSQL, pool, action, taskID, cid, uid, result, reason, utime)
	if err != nil {
		log.Error("tx.Exec(%s) error(%v)", _inTaskHisSQL, err)
		return
	}
	return res.RowsAffected()
}

// AddTaskHis 非事务
func (d *Dao) AddTaskHis(c context.Context, pool, action int8, taskID, cid, uid, utime int64, result int16, reason string) (rows int64, err error) {
	res, err := d.db.Exec(c, _inTaskHisSQL, pool, action, taskID, cid, uid, result, reason, utime)
	if err != nil {
		log.Error("d.db.Exec(%s) error(%v)", _inTaskHisSQL, err)
		return
	}
	return res.RowsAffected()
}

// MulAddTaskHis 批量插入日志
func (d *Dao) MulAddTaskHis(c context.Context, tls []*archive.TaskForLog, action int8, uid int64) (rows int64, err error) {
	params := []string{}
	for _, item := range tls {
		var itemstr string
		itemstr += fmt.Sprintf("(%d,", action)
		itemstr += fmt.Sprintf("%d,", item.ID)
		itemstr += fmt.Sprintf("%d,", item.Cid)
		itemstr += fmt.Sprintf("%d)", uid)
		params = append(params, itemstr)
	}
	if len(params) == 0 {
		log.Warn("MulAddTaskHis empty params")
		return
	}
	sqlsring := strings.Join(params, ",")
	res, err := d.db.Exec(c, _mulinTaskHisSQL+sqlsring)
	if err != nil {
		log.Error("d.db.Exec(%s, %s) error(%v)", _mulinTaskHisSQL, sqlsring, err)
		return
	}

	return res.RowsAffected()
}
