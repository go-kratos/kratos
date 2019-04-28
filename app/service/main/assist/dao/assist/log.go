package assist

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-common/app/service/main/assist/model/assist"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	// insert
	_inLogSQL = "INSERT IGNORE INTO assist_log (mid,assist_mid,type,action,subject_id,object_id,detail) VALUES (?,?,?,?,?,?,?)"
	// update
	_cancelLogSQL = "UPDATE assist_log SET state=1 WHERE id=? AND mid=? AND assist_mid=?"
	// info
	_logInfoSQL = "SELECT id,mid,assist_mid,type,action,subject_id,object_id,detail,state,ctime,mtime FROM assist_log WHERE id=? AND mid=? AND assist_mid=?"
	// obj
	_logObjSQL = "SELECT id,mid,assist_mid,type,action,subject_id,object_id,detail,state,ctime,mtime FROM assist_log WHERE mid=? AND object_id=? AND type=? AND action=? limit 1"
	// select
	_logsSQL                = "SELECT id,mid,assist_mid,type,action,subject_id,object_id,detail,state,ctime,mtime FROM assist_log WHERE mid=? ORDER BY id DESC LIMIT ?,?"
	_logsByAssSQL           = "SELECT id,mid,assist_mid,type,action,subject_id,object_id,detail,state,ctime,mtime FROM assist_log WHERE mid=? AND assist_mid=? ORDER BY id DESC LIMIT ?,?"
	_logsByCtimeSQL         = "SELECT id,mid,assist_mid,type,action,subject_id,object_id,detail,state,ctime,mtime FROM assist_log WHERE mid=? AND ctime>=? AND ctime<=? ORDER BY id DESC LIMIT ?,?"
	_logsByAssCtimeSQL      = "SELECT id,mid,assist_mid,type,action,subject_id,object_id,detail,state,ctime,mtime FROM assist_log WHERE mid=? AND assist_mid=? AND ctime>=? AND ctime<=? ORDER BY id DESC LIMIT ?,?"
	_logCntGroupBySQL       = "SELECT assist_mid,type,action,count(*) FROM assist_log WHERE mid=? GROUP BY assist_mid,type,action"
	_logAssMidCntGroupBySQL = "SELECT assist_mid,type,action,count(*) FROM assist_log WHERE mid=%d and assist_mid IN (%s) GROUP BY assist_mid,type,action"
	// LogCntBy*SQL
	_logCntSQL           = "SELECT count(*) FROM assist_log WHERE mid=?"
	_logCntByAssSQL      = "SELECT count(*) FROM assist_log WHERE mid=? AND assist_mid=?"
	_logCntByCtimeSQL    = "SELECT count(*) FROM assist_log WHERE mid=? AND ctime>=? AND ctime<=?"
	_logCntByAssCtimeSQL = "SELECT count(*) FROM assist_log WHERE mid=? AND assist_mid=? AND ctime>=? AND ctime<=?"
)

// AddLog add one assist log.
func (d *Dao) AddLog(c context.Context, mid, assistMid, tp, act, subID int64, objIDStr string, detail string) (id int64, err error) {
	res, err := d.db.Exec(c, _inLogSQL, mid, assistMid, tp, act, subID, objIDStr, detail)
	if err != nil {
		log.Error("d.inLog error(%v)|(%d,%d,%d,%d,%d,%s,%d)", err, mid, assistMid, tp, act, subID, objIDStr, detail)
		return
	}
	id, err = res.LastInsertId()
	return
}

// CancelLog cancel assist oper log.
func (d *Dao) CancelLog(c context.Context, logID, mid, assistMid int64) (rows int64, err error) {
	res, err := d.db.Exec(c, _cancelLogSQL, logID, mid, assistMid)
	if err != nil {
		log.Error("d.db.Exec error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// LogInfo log Info
func (d *Dao) LogInfo(c context.Context, id, mid, assistMid int64) (a *assist.Log, err error) {
	row := d.db.QueryRow(c, _logInfoSQL, id, mid, assistMid)
	a = &assist.Log{}
	if err = row.Scan(&a.ID, &a.Mid, &a.AssistMid, &a.Type, &a.Action, &a.SubjectID, &a.ObjectID, &a.Detail, &a.State, &a.CTime, &a.MTime); err != nil {
		if err == sql.ErrNoRows {
			a = nil
			err = ecode.AssistLogNotExist
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// Logs get logs from assist db.
func (d *Dao) Logs(c context.Context, mid int64, start, offset int) (logs []*assist.Log, err error) {
	logs = make([]*assist.Log, 0)
	rows, err := d.db.Query(c, _logsSQL, mid, start, offset)
	if err != nil {
		log.Error("db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		lg := &assist.Log{}
		if err = rows.Scan(&lg.ID, &lg.Mid, &lg.AssistMid, &lg.Type, &lg.Action, &lg.SubjectID, &lg.ObjectID, &lg.Detail, &lg.State, &lg.CTime, &lg.MTime); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		logs = append(logs, lg)
	}
	return
}

// LogsByAssist get logs from assist db by assist mid.
func (d *Dao) LogsByAssist(c context.Context, mid, assistMid int64, start, offset int) (logs []*assist.Log, err error) {
	logs = make([]*assist.Log, 0)
	rows, err := d.db.Query(c, _logsByAssSQL, mid, assistMid, start, offset)
	if err != nil {
		log.Error("db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		lg := &assist.Log{}
		if err = rows.Scan(&lg.ID, &lg.Mid, &lg.AssistMid, &lg.Type, &lg.Action, &lg.SubjectID, &lg.ObjectID, &lg.Detail, &lg.State, &lg.CTime, &lg.MTime); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		logs = append(logs, lg)
	}
	return
}

// LogsByCtime get logs from assist db by ctime.
func (d *Dao) LogsByCtime(c context.Context, mid int64, stime, etime time.Time, start, offset int) (logs []*assist.Log, err error) {
	logs = make([]*assist.Log, 0)
	rows, err := d.db.Query(c, _logsByCtimeSQL, mid, stime, etime, start, offset)
	if err != nil {
		log.Error("db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		lg := &assist.Log{}
		if err = rows.Scan(&lg.ID, &lg.Mid, &lg.AssistMid, &lg.Type, &lg.Action, &lg.SubjectID, &lg.ObjectID, &lg.Detail, &lg.State, &lg.CTime, &lg.MTime); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		logs = append(logs, lg)
	}
	return
}

// LogsByAssistCtime get logs from assist db by assist oper ctime.
func (d *Dao) LogsByAssistCtime(c context.Context, mid, assistMid int64, stime, etime time.Time, start, offset int) (logs []*assist.Log, err error) {
	logs = make([]*assist.Log, 0)
	rows, err := d.db.Query(c, _logsByAssCtimeSQL, mid, assistMid, stime, etime, start, offset)
	if err != nil {
		log.Error("db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		lg := &assist.Log{}
		if err = rows.Scan(&lg.ID, &lg.Mid, &lg.AssistMid, &lg.Type, &lg.Action, &lg.SubjectID, &lg.ObjectID, &lg.Detail, &lg.State, &lg.CTime, &lg.MTime); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		logs = append(logs, lg)
	}
	return
}

// LogCount count by group from db.
func (d *Dao) LogCount(c context.Context, mid int64) (totalm map[int64]map[int8]map[int8]int, err error) {
	rows, err := d.db.Query(c, _logCntGroupBySQL, mid)
	if err != nil {
		log.Error("db.Query err(%v)", err)
		return
	}
	defer rows.Close()
	totalm = map[int64]map[int8]map[int8]int{}
	for rows.Next() {
		var (
			assMid  int64
			tp, act int8
			cnt     int
		)
		if err = rows.Scan(&assMid, &tp, &act, &cnt); err != nil {
			log.Error("row.Scan err(%v)", err)
			return
		}
		if tassMap, ok := totalm[assMid]; !ok {
			totalm[assMid] = map[int8]map[int8]int{
				tp: {act: cnt},
			}
		} else {
			if tpMap, ok := tassMap[tp]; !ok {
				tassMap[tp] = map[int8]int{act: cnt}
			} else {
				tpMap[act] = cnt
			}
		}
	}
	return
}

// LogCnt fn
func (d *Dao) LogCnt(c context.Context, mid int64) (count int64, err error) {
	row := d.db.QueryRow(c, _logCntSQL, mid)
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			count = 0
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// LogCntAssist fn
func (d *Dao) LogCntAssist(c context.Context, mid, assistMid int64) (count int64, err error) {
	row := d.db.QueryRow(c, _logCntByAssSQL, mid, assistMid)
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			count = 0
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// LogCntCtime fn
func (d *Dao) LogCntCtime(c context.Context, mid int64, stime, etime time.Time) (count int64, err error) {
	row := d.db.QueryRow(c, _logCntByCtimeSQL, mid, stime, etime)
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			count = 0
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// LogCntAssistCtime fn
func (d *Dao) LogCntAssistCtime(c context.Context, mid, assistMid int64, stime, etime time.Time) (count int64, err error) {
	row := d.db.QueryRow(c, _logCntByAssCtimeSQL, mid, assistMid, stime, etime)
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			count = 0
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// AssistsMidsTotal get all Assists from assist database.
func (d *Dao) AssistsMidsTotal(c context.Context, mid int64, assmids []int64) (totalm map[int64]map[int8]map[int8]int, err error) {
	query := fmt.Sprintf(_logAssMidCntGroupBySQL, mid, xstr.JoinInts(assmids))
	rows, err := d.db.Query(c, query)
	if err != nil {
		log.Error("db.Query err(%v)", err)
		return
	}
	defer rows.Close()
	totalm = map[int64]map[int8]map[int8]int{}
	for rows.Next() {
		var (
			assMid  int64
			tp, act int8
			cnt     int
		)
		if err = rows.Scan(&assMid, &tp, &act, &cnt); err != nil {
			log.Error("row.Scan err(%v)", err)
			return
		}
		if tassMap, ok := totalm[assMid]; !ok {
			totalm[assMid] = map[int8]map[int8]int{
				tp: {act: cnt},
			}
		} else {
			if tpMap, ok := tassMap[tp]; !ok {
				tassMap[tp] = map[int8]int{act: cnt}
			} else {
				tpMap[act] = cnt
			}
		}
	}
	return
}

// LogObj log Obj
func (d *Dao) LogObj(c context.Context, mid, objID, tp, act int64) (a *assist.Log, err error) {
	row := d.db.QueryRow(c, _logObjSQL, mid, objID, tp, act)
	a = &assist.Log{}
	if err = row.Scan(&a.ID, &a.Mid, &a.AssistMid, &a.Type, &a.Action, &a.SubjectID, &a.ObjectID, &a.Detail, &a.State, &a.CTime, &a.MTime); err != nil {
		if err == sql.ErrNoRows {
			a = nil
			err = ecode.AssistLogNotExist
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}
