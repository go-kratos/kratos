package archive

import (
	"context"
	xsql "database/sql"
	"fmt"
	"strings"
	"time"

	"go-common/app/admin/main/videoup/model/archive"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_upTaskByIDSQL     = "UPDATE task_dispatch SET %s WHERE id=?"
	_upGtimeByIDSQL    = "UPDATE task_dispatch SET gtime=? WHERE id=?"
	_releaseByIDSQL    = "UPDATE task_dispatch SET subject=0,state=0,uid=0,gtime='0000-00-00 00:00:00' WHERE id=?"
	_releaseMtimeSQL   = "UPDATE task_dispatch SET subject=0,state=0,uid=0,gtime='0000-00-00 00:00:00' WHERE id IN (%s) AND mtime<=?"
	_timeOutTaskSQL    = "SELECT id,cid,subject,mtime FROM task_dispatch WHERE (state=1 AND mtime<?) OR (state=0 AND uid<>0 AND ctime<?)"
	_getRelTaskSQL     = "SELECT id,cid,subject,mtime,gtime FROM task_dispatch WHERE state IN (0,1) AND uid=?"
	_releaseSpecialSQL = "UPDATE task_dispatch SET subject=0,state=0,uid=0 WHERE id=? AND gtime='0000-00-00 00:00:00' AND mtime<=? AND state=? AND uid=?"
)

// UpGtimeByID update gtime
func (d *Dao) UpGtimeByID(c context.Context, id int64, gtime string) (rows int64, err error) {
	var res xsql.Result

	if res, err = d.db.Exec(c, _upGtimeByIDSQL, gtime, id); err != nil {
		log.Error("d.db.Exec(%s, %v, %d) error(%v)", _upGtimeByIDSQL, gtime, id)
		return
	}
	return res.RowsAffected()
}

// TxUpTaskByID 更新任务状态
func (d *Dao) TxUpTaskByID(tx *sql.Tx, id int64, paras map[string]interface{}) (rows int64, err error) {
	arrSet := []string{}
	arrParas := []interface{}{}
	for k, v := range paras {
		arrSet = append(arrSet, k+"=?")
		arrParas = append(arrParas, v)
	}
	arrParas = append(arrParas, id)
	sqlstring := fmt.Sprintf(_upTaskByIDSQL, strings.Join(arrSet, ","))
	res, err := tx.Exec(sqlstring, arrParas...)
	if err != nil {
		log.Error("tx.Exec(%v %v) error(%v)", sqlstring, arrParas, err)
		return
	}
	return res.RowsAffected()
}

// TxReleaseByID 释放指定任务
func (d *Dao) TxReleaseByID(tx *sql.Tx, id int64) (rows int64, err error) {
	res, err := tx.Exec(_releaseByIDSQL, id)
	if err != nil {
		log.Error("tx.Exec(%s, %d) error(%v)", _releaseByIDSQL, id, err)
		return
	}
	return res.RowsAffected()
}

// MulReleaseMtime 批量释放任务,加时间防止释放错误
func (d *Dao) MulReleaseMtime(c context.Context, ids []int64, mtime time.Time) (rows int64, err error) {
	sqlstring := fmt.Sprintf(_releaseMtimeSQL, xstr.JoinInts(ids))
	res, err := d.db.Exec(c, sqlstring, mtime)
	if err != nil {
		log.Error("tx.Exec(%s,  %v) error(%v)", sqlstring, mtime, err)
		return
	}
	return res.RowsAffected()
}

// GetTimeOutTask 释放正在处理且超时的,释放指派后但长时间未审核的
func (d *Dao) GetTimeOutTask(c context.Context) (rts []*archive.TaskForLog, err error) {
	var (
		rows *sql.Rows
	)
	if rows, err = d.rddb.Query(c, _timeOutTaskSQL, time.Now().Add(-10*time.Minute), time.Now().Add(-80*time.Minute)); err != nil {
		log.Error("d.rddb.Query(%s) error(%v)", _timeOutTaskSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		rt := &archive.TaskForLog{}
		if err = rows.Scan(&rt.ID, &rt.Cid, &rt.Subject, &rt.Mtime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		rts = append(rts, rt)
	}
	return
}

// GetRelTask 用户登出或者主动释放(分配给该用户的都释放)
func (d *Dao) GetRelTask(c context.Context, uid int64) (rts []*archive.TaskForLog, lastid int64, err error) {
	var (
		gtime time.Time
		rows  *sql.Rows
	)
	if rows, err = d.rddb.Query(c, _getRelTaskSQL, uid); err != nil {
		log.Error("d.rddb.Query(%s, %d) error(%v)", _getRelTaskSQL, uid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		rt := &archive.TaskForLog{}
		if err = rows.Scan(&rt.ID, &rt.Cid, &rt.Subject, &rt.Mtime, &gtime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		if gtime.IsZero() {
			rts = append(rts, rt)
		} else {
			lastid = rt.ID
		}
	}
	return
}

// TxReleaseSpecial 延时固定时间释放的任务，需要校验释放时的状态，时间，认领人等
func (d *Dao) TxReleaseSpecial(tx *sql.Tx, mtime time.Time, state int8, taskid, uid int64) (rows int64, err error) {
	res, err := tx.Exec(_releaseSpecialSQL, taskid, mtime, state, uid)
	if err != nil {
		log.Error("tx.Exec(%s, %d, %v, %d, %d) error(%v)", _releaseSpecialSQL, taskid, mtime, state, uid, err)
		return
	}
	return res.RowsAffected()
}
