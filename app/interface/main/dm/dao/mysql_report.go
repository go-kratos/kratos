package dao

import (
	"context"
	"fmt"

	"go-common/app/interface/main/dm/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_rptSharding     = 100
	_rptUserSharding = 100
	_rptLogSharding  = 100
	_insertRpt       = "INSERT INTO dm_report_%d (cid,dmid,uid,reason,content,count,state,up_op,rp_time,score,ctime,mtime) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)" +
		" ON DUPLICATE KEY UPDATE uid=?,reason=?,state=?,content=?,count=count+1,rp_time=?,score=score+?"
	_selectRpt         = "SELECT id,cid,dmid,uid,reason,content,count,state,up_op,score,rp_time,ctime,mtime FROM dm_report_%d WHERE dmid=?"
	_updateRptUPOp     = "UPDATE dm_report_%d SET up_op=? WHERE dmid=?"
	_updateRpt         = "UPDATE dm_report_%d SET state=? WHERE dmid=?"
	_insertUser        = "INSERT IGNORE INTO dm_report_user_%d (dmid,uid,reason,state,content,ctime,mtime) VALUES(?,?,?,?,?,?,?)"
	_selectRptLog      = "SELECT id,dmid,adminid,reason,result,remark,elapsed,ctime,mtime FROM dm_report_admin_log_%d WHERE dmid=? ORDER BY mtime"
	_insertRptLog      = "INSERT INTO dm_report_admin_log_%d (dmid,adminid,reason,result,remark,elapsed) VALUES (?,?,?,?,?,?)"
	_selectRptUser     = "SELECT id,dmid,uid,reason,state,ctime,mtime FROM dm_report_user_%d WHERE dmid=? AND state=?"
	_updateRptUserStat = "UPDATE dm_report_user_%d SET state=? WHERE dmid=? AND state!=?"
)

func rptTable(cid int64) int64 {
	return cid % _rptSharding
}

func rptUserTable(dmid int64) int64 {
	return dmid % _rptUserSharding
}

// RptLogTable return log table id by dmid
func RptLogTable(dmid int64) int64 {
	return dmid % _rptLogSharding
}

// AddReport insert or update dm report.
func (d *Dao) AddReport(c context.Context, rpt *model.Report) (id int64, err error) {
	res, err := d.biliDM.Exec(c, fmt.Sprintf(_insertRpt, rptTable(rpt.Cid)), rpt.Cid, rpt.Did, rpt.UID, rpt.Reason, rpt.Content, rpt.Count, rpt.State, rpt.UpOP, rpt.RpTime, rpt.Score, rpt.Ctime, rpt.Mtime, rpt.UID, rpt.Reason, rpt.State, rpt.Content, rpt.RpTime, rpt.Score)
	if err != nil {
		log.Error("d.AddReport(%v) error(%v)", rpt, err)
		return
	}
	return res.LastInsertId()
}

// AddReportUser add dm report user.
func (d *Dao) AddReportUser(c context.Context, u *model.User) (id int64, err error) {
	res, err := d.biliDM.Exec(c, fmt.Sprintf(_insertUser, rptUserTable(u.Did)), u.Did, u.UID, u.Reason, u.State, u.Content, u.Ctime, u.Mtime)
	if err != nil {
		log.Error("d.AddReportUser(%v) error(%v)", u, err)
		return
	}
	return res.LastInsertId()
}

// ReportLog get report log of dmid.
func (d *Dao) ReportLog(c context.Context, dmid int64) (res []*model.RptLog, err error) {
	rows, err := d.biliDM.Query(c, fmt.Sprintf(_selectRptLog, RptLogTable(dmid)), dmid)
	if err != nil {
		log.Error("d.ReportLog(dmid:%d) error(%v)", dmid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.RptLog{}
		if err = rows.Scan(&r.ID, &r.Did, &r.AdminID, &r.Reason, &r.Result, &r.Remark, &r.Elapsed, &r.Ctime, &r.Mtime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		res = append(res, r)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// Report dm report info by cid and dmid.
func (d *Dao) Report(c context.Context, cid, dmid int64) (rpt *model.Report, err error) {
	rpt = &model.Report{}
	row := d.biliDM.QueryRow(c, fmt.Sprintf(_selectRpt, rptTable(cid)), dmid)
	err = row.Scan(&rpt.ID, &rpt.Cid, &rpt.Did, &rpt.UID, &rpt.Reason, &rpt.Content, &rpt.Count, &rpt.State, &rpt.UpOP, &rpt.Score, &rpt.RpTime, &rpt.Ctime, &rpt.Mtime)
	if err != nil {
		if err == sql.ErrNoRows {
			rpt = nil
			err = nil
		} else {
			log.Error("row.Scan() error(%v)", err)
		}
	}
	return
}

// UpdateReportStat update report state
func (d *Dao) UpdateReportStat(c context.Context, cid, dmid int64, state int8) (affect int64, err error) {
	sqlStr := fmt.Sprintf(_updateRpt, rptTable(cid))
	res, err := d.biliDM.Exec(c, sqlStr, state, dmid)
	if err != nil {
		log.Error("d.UpdateReportStat(cid:%d, dmid:%d) error(%v)", cid, dmid, err)
		return
	}
	return res.RowsAffected()
}

// UpdateReportUPOp update dm report state.
func (d *Dao) UpdateReportUPOp(c context.Context, cid, dmid int64, op int8) (affect int64, err error) {
	res, err := d.biliDM.Exec(c, fmt.Sprintf(_updateRptUPOp, rptTable(cid)), op, dmid)
	if err != nil {
		log.Error("d.UpdateReportUPOp(cid:%d, dmid:%d) error(%v)", cid, dmid, err)
		return
	}
	return res.RowsAffected()
}

// AddReportLog add report log.
func (d *Dao) AddReportLog(c context.Context, lg *model.RptLog) (err error) {
	sqlStr := fmt.Sprintf(_insertRptLog, RptLogTable(lg.Did))
	_, err = d.biliDM.Exec(c, sqlStr, lg.Did, lg.AdminID, lg.Reason, lg.Result, lg.Remark, lg.Elapsed)
	if err != nil {
		log.Error("d.AddReportLog(%v) error(%v)", lg, err)
	}
	return
}

// ReportUser return report use list of dmid.
func (d *Dao) ReportUser(c context.Context, dmid int64) (users []*model.User, err error) {
	sqlStr := fmt.Sprintf(_selectRptUser, rptUserTable(dmid))
	rows, err := d.biliDM.Query(c, sqlStr, dmid, model.NoticeUnsend)
	if err != nil {
		log.Error("d.ReportUser(query:%s, dmid:%d) error(%v)", sqlStr, dmid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		u := &model.User{}
		if err = rows.Scan(&u.ID, &u.Did, &u.UID, &u.Reason, &u.State, &u.Ctime, &u.Mtime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		users = append(users, u)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// SetReportUserFinished set dmid state to be noticesend.
func (d *Dao) SetReportUserFinished(c context.Context, dmid int64) (err error) {
	sqlStr := fmt.Sprintf(_updateRptUserStat, rptUserTable(dmid))
	if _, err = d.biliDM.Exec(c, sqlStr, model.NoticeSend, dmid, model.NoticeSend); err != nil {
		log.Error("d.SetReportUserFinished(sql:%s, dmid:%d) error(%v)", sqlStr, dmid, err)
	}
	return
}
