package dao

import (
	"bytes"
	"context"
	"fmt"

	"go-common/app/admin/main/dm/model"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_userSharding   = 100
	_logSharding    = 100
	_reportSharding = 100
	// dm report
	_updateStatSQL = "UPDATE dm_report_%d SET state=? WHERE dmid IN (%s)"
	_ignoreStatSQL = "UPDATE dm_report_%d SET state=?,count=0 WHERE dmid IN (%s)"
	_selectRpt     = "SELECT id,dmid,cid,uid,reason,count,up_op,state,score,rp_time,ctime,mtime FROM dm_report_%d WHERE dmid IN (%s)"
	_selectUsers   = "SELECT id,dmid,uid,reason,state,ctime,mtime FROM dm_report_user_%d WHERE dmid IN (%s) AND state=? ORDER BY id"
	_updateUser    = "UPDATE dm_report_user_%d SET state=? WHERE dmid IN (%s) AND state !=?"
	_insertLog     = "INSERT dm_report_admin_log_%d (dmid,adminid,reason,result,remark,elapsed) VALUES"
	_selectLog     = "SELECT id,dmid,adminid,reason,result,remark,elapsed,ctime,mtime FROM dm_report_admin_log_%d WHERE dmid=? ORDER BY mtime"
)

// RptTable return report table id by cid
func RptTable(cid int64) int64 {
	return cid % _reportSharding
}

// UserTable return user table id by dmid
func UserTable(dmid int64) int64 {
	return dmid % _userSharding
}

// LogTable return log table id by dmid
func LogTable(dmid int64) int64 {
	return dmid % _logSharding
}

// ChangeReportStat change report state
func (d *Dao) ChangeReportStat(c context.Context, cid int64, dmids []int64, state int8) (err error) {
	updateSQL := fmt.Sprintf(_updateStatSQL, RptTable(cid), xstr.JoinInts(dmids))
	if _, err = d.biliDM.Exec(c, updateSQL, state); err != nil {
		log.Error("d.biliDM.Exec(%d) error(%v)", state, err)
	}
	return
}

// IgnoreReport change report state to SecondIgnore or FirstIngnore
func (d *Dao) IgnoreReport(c context.Context, cid int64, dmids []int64, state int8) (err error) {
	updateSQL := fmt.Sprintf(_ignoreStatSQL, RptTable(cid), xstr.JoinInts(dmids))
	if _, err = d.biliDM.Exec(c, updateSQL, state); err != nil {
		log.Error("d.biliDM.Exec(%d) error(%v)", state, err)
	}
	return
}

// Reports get multi dm report info.
func (d *Dao) Reports(c context.Context, cid int64, dmids []int64) (res []*model.Report, err error) {
	res = []*model.Report{}
	selectSQL := fmt.Sprintf(_selectRpt, RptTable(cid), xstr.JoinInts(dmids))
	rows, err := d.biliDM.Query(c, selectSQL)
	if err != nil {
		log.Error("d.biliDM.Exec(cid:%d,dmids:%v) error(%v)", cid, dmids, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.Report{}
		err = rows.Scan(&r.ID, &r.Did, &r.Cid, &r.UID, &r.RpType, &r.Count, &r.UpOP, &r.State, &r.Score, &r.RpTime, &r.Ctime, &r.Mtime)
		if err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		res = append(res, r)
	}
	return
}

// ReportUsers get user list.
func (d *Dao) ReportUsers(c context.Context, tableID int64, dmids []int64, state int8) (res map[int64][]*model.ReportUser, err error) {
	res = make(map[int64][]*model.ReportUser, 100)
	selectSQL := fmt.Sprintf(_selectUsers, tableID, xstr.JoinInts(dmids))
	rows, err := d.biliDM.Query(c, selectSQL, state)
	if err != nil {
		log.Error("d.biliDM.Query(sql:%s) error(%v)", selectSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		u := &model.ReportUser{}
		if err = rows.Scan(&u.ID, &u.Did, &u.UID, &u.Reason, &u.State, &u.Ctime, &u.Mtime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		res[u.Did] = append(res[u.Did], u)
	}
	return
}

// UpReportUserState update report user state.
func (d *Dao) UpReportUserState(c context.Context, tableID int64, dmids []int64, state int8) (affect int64, err error) {
	selectSQL := fmt.Sprintf(_updateUser, tableID, xstr.JoinInts(dmids))
	res, err := d.biliDM.Exec(c, selectSQL, state, state)
	if err != nil {
		log.Error("d.updateUserStmt.Exec(dmid:%v) error(%v)", dmids, err)
		return
	}
	return res.RowsAffected()
}

// AddReportLog add dm report admin log.
func (d *Dao) AddReportLog(c context.Context, tableID int64, lg []*model.ReportLog) (err error) {
	var (
		buffer   bytes.Buffer
		insertTp string
	)
	insertTp = "(%d,%d,%d,%d,'%s',%d),"
	buffer.WriteString(fmt.Sprintf(_insertLog, tableID))
	for _, v := range lg {
		buffer.WriteString(fmt.Sprintf(insertTp, v.Did, v.AdminID, v.Reason, v.Result, v.Remark, v.Elapsed))
	}
	//truncate the last ','
	buffer.Truncate(buffer.Len() - 1)
	_, err = d.biliDM.Exec(c, buffer.String())
	if err != nil {
		log.Error("d.insertLogStmt.Exec(%v) error(%v)", lg, err)
		return
	}
	return
}

// ReportLog get dm report log.
func (d *Dao) ReportLog(c context.Context, dmid int64) (res []*model.ReportLog, err error) {
	selectSQL := fmt.Sprintf(_selectLog, LogTable(dmid))
	rows, err := d.biliDM.Query(c, selectSQL, dmid)
	if err != nil {
		log.Error("dmreport:d.biliDM.Query(sql:%s) error(%v)", selectSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.ReportLog{}
		if err = rows.Scan(&r.ID, &r.Did, &r.AdminID, &r.Reason, &r.Result, &r.Remark, &r.Elapsed, &r.Ctime, &r.Mtime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		res = append(res, r)
	}
	return
}
