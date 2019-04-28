package dao

import (
	"context"
	sql "database/sql"
	"fmt"
	"time"

	"go-common/app/admin/main/reply/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	// report
	_upRptStateSQL    = "UPDATE reply_report_%d SET state=?,mtime=? WHERE rpid IN (%s)"
	_upRptReasonSQL   = "UPDATE reply_report_%d SET state=?,reason=?,content=?,mtime=? WHERE rpid IN (%s)"
	_upRptAttrBitSQL  = "UPDATE reply_report_%d SET attr=attr&(~(1<<?))|(?<<?),mtime=? WHERE rpid IN (%s)"
	_selReportSQL     = "SELECT oid,type,rpid,mid,reason,content,count,score,state,attr,ctime,mtime FROM reply_report_%d WHERE rpid=?"
	_selReportsSQL    = "SELECT oid,type,rpid,mid,reason,content,count,score,state,attr,ctime,mtime FROM reply_report_%d WHERE rpid IN(%s)"
	_selReportOidsSQL = "SELECT oid,type,rpid,mid,reason,content,count,score,state,attr,ctime,mtime FROM reply_report_%d WHERE oid IN(?) and type=?"
	// report_user
	_setRptUserStateSQL = "UPDATE reply_report_user_%d SET state=?,mtime=? WHERE rpid=?"
	_selRptUsersSQL     = "SELECT oid,type,rpid,mid,reason,content,state,ctime,mtime FROM reply_report_user_%d WHERE rpid=? and state=?"
)

// Report get a reply report.
func (d *Dao) Report(c context.Context, oid, rpID int64) (rpt *model.Report, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_selReportSQL, hit(oid)), rpID)
	rpt = new(model.Report)
	if err = row.Scan(&rpt.Oid, &rpt.Type, &rpt.RpID, &rpt.Mid, &rpt.Reason, &rpt.Content, &rpt.Count, &rpt.Score, &rpt.State, &rpt.Attr, &rpt.CTime, &rpt.MTime); err != nil {
		if err == xsql.ErrNoRows {
			rpt = nil
			err = nil
		}
	}
	return
}

// Reports return report map by oid.
func (d *Dao) Reports(c context.Context, oids, rpIDs []int64) (res map[int64]*model.Report, err error) {
	hits := make(map[int64][]int64)
	for i, oid := range oids {
		hit := hit(oid)
		hits[hit] = append(hits[hit], rpIDs[i])
	}
	res = make(map[int64]*model.Report)
	for hit, ids := range hits {
		var rows *xsql.Rows
		if rows, err = d.db.Query(c, fmt.Sprintf(_selReportsSQL, hit, xstr.JoinInts(ids))); err != nil {
			return
		}
		for rows.Next() {
			rpt := new(model.Report)
			if err = rows.Scan(&rpt.Oid, &rpt.Type, &rpt.RpID, &rpt.Mid, &rpt.Reason, &rpt.Content, &rpt.Count, &rpt.Score, &rpt.State, &rpt.Attr, &rpt.CTime, &rpt.MTime); err != nil {
				rows.Close()
				return
			}
			res[rpt.RpID] = rpt
		}
		if err = rows.Err(); err != nil {
			rows.Close()
			return
		}
		rows.Close()
	}
	return
}

// ReportByOids return report map by oid.
func (d *Dao) ReportByOids(c context.Context, typ int32, oids ...int64) (res map[int64]*model.Report, err error) {
	hits := make(map[int64][]int64)
	for _, oid := range oids {
		hit := hit(oid)
		hits[hit] = append(hits[hit], oid)
	}
	res = make(map[int64]*model.Report)
	for hit, oids := range hits {
		var rows *xsql.Rows
		if rows, err = d.db.Query(c, fmt.Sprintf(_selReportOidsSQL, hit), xstr.JoinInts(oids), typ); err != nil {
			return
		}
		for rows.Next() {
			rpt := new(model.Report)
			if err = rows.Scan(&rpt.Oid, &rpt.Type, &rpt.RpID, &rpt.Mid, &rpt.Reason, &rpt.Content, &rpt.Count, &rpt.Score, &rpt.State, &rpt.Attr, &rpt.CTime, &rpt.MTime); err != nil {
				rows.Close()
				return
			}
			res[rpt.RpID] = rpt
		}
		if err = rows.Err(); err != nil {
			rows.Close()
			return
		}
		rows.Close()
	}
	return
}

// UpReportsState update the report state.
func (d *Dao) UpReportsState(c context.Context, oids, rpIDs []int64, state int32, now time.Time) (rows int64, err error) {
	hitMap := make(map[int64][]int64)
	for i, oid := range oids {
		hitMap[hit(oid)] = append(hitMap[hit(oid)], rpIDs[i])
	}

	for hit, ids := range hitMap {
		var res sql.Result
		res, err = d.db.Exec(c, fmt.Sprintf(_upRptStateSQL, hit, xstr.JoinInts(ids)), state, now)
		if err != nil {
			log.Error("mysqlDB.Exec error(%v)", err)
			return
		}
		var row int64
		row, err = res.RowsAffected()
		if err != nil {
			log.Error("res.RowsAffected error(%v)", err)
			return
		}
		rows += row
	}
	return
}

// UpReportsStateWithReason update the report state.
func (d *Dao) UpReportsStateWithReason(c context.Context, oids, rpIDs []int64, state, reason int32, content string, now time.Time) (rows int64, err error) {
	hitMap := make(map[int64][]int64)
	for i, oid := range oids {
		hitMap[hit(oid)] = append(hitMap[hit(oid)], rpIDs[i])
	}

	for hit, ids := range hitMap {
		var res sql.Result
		res, err = d.db.Exec(c, fmt.Sprintf(_upRptReasonSQL, hit, xstr.JoinInts(ids)), state, reason, content, now)
		if err != nil {
			log.Error("mysqlDB.Exec error(%v)", err)
			return
		}
		var row int64
		row, err = res.RowsAffected()
		if err != nil {
			log.Error("res.RowsAffected error(%v)", err)
			return
		}
		rows += row
	}
	return
}

// UpReportsAttrBit update the report attr.
func (d *Dao) UpReportsAttrBit(c context.Context, oids, rpIDs []int64, bit uint32, val uint32, now time.Time) (rows int64, err error) {
	hitMap := make(map[int64][]int64)
	for i, oid := range oids {
		hitMap[hit(oid)] = append(hitMap[hit(oid)], rpIDs[i])
	}
	for hit, ids := range hitMap {
		var res sql.Result
		res, err = d.db.Exec(c, fmt.Sprintf(_upRptAttrBitSQL, hit, xstr.JoinInts(ids)), bit, val, bit, now)
		if err != nil {
			log.Error("mysqlDB.Exec error(%v)", err)
			return
		}
		var row int64
		row, err = res.RowsAffected()
		if err != nil {
			log.Error("res.RowsAffected error(%v)", err)
			return
		}
		rows += row
	}
	return
}

// ReportUsers return a report users from mysql.
func (d *Dao) ReportUsers(c context.Context, oid int64, tp int32, rpID int64) (res map[int64]*model.ReportUser, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_selRptUsersSQL, hit(oid)), rpID, model.ReportUserStateNew)
	if err != nil {
		log.Error("db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	res = make(map[int64]*model.ReportUser)
	for rows.Next() {
		rpt := &model.ReportUser{}
		if err = rows.Scan(&rpt.Oid, &rpt.Type, &rpt.RpID, &rpt.Mid, &rpt.Reason, &rpt.Content, &rpt.State, &rpt.CTime, &rpt.MTime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		res[rpt.Mid] = rpt
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.err error(%v)", err)
		return
	}
	return
}

// SetUserReported set a user report state by rpID.
func (d *Dao) SetUserReported(c context.Context, oid int64, tp int32, rpID int64, now time.Time) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_setRptUserStateSQL, hit(oid)), model.ReportUserStateReported, now, rpID)
	if err != nil {
		log.Error("db.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}
