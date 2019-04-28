package dao

import (
	"context"
	"fmt"
	"strings"

	"go-common/app/admin/main/tag/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_reportSQL             = "SELECT id,oid,type,tid,mid,action,prid,rid,count,reason,content,moral,score,state,ctime,mtime FROM report WHERE id=?"
	_reportsSQL            = "SELECT id,oid,type,tid,mid,action,prid,rid,count,reason,content,moral,score,state,ctime,mtime FROM report WHERE id in (%s)"
	_updateReportStateSQL  = "UPDATE report SET state=? WHERE id=?"
	_updateReportStatesSQL = "UPDATE report SET state=? WHERE id in (%s)"
	_insertReportLogSQL    = "INSERT INTO report_log(rpt,username,points,oid,type,mid,tid,rid,reason,handle_type,notice) VALUES (%d,%q,%d,%d,%d,%d,%d,%d,%q,%d,%d)"
	_insertReportLogsSQL   = "INSERT INTO report_log(rpt,username,points,oid,type,mid,tid,rid,reason,handle_type,notice) VALUES %s"
	_updateReportMoralSQL  = "UPDATE report SET moral=? WHERE id in (%s)"
	_reportByOidMidSQL     = "SELECT id,oid,type,tid,mid,action,prid,rid,count,reason,content,moral,state,ctime,mtime FROM report WHERE mid=? AND oid=? AND type=? ORDER BY id DESC"
	_reportUsersSQL        = "SELECT id,rpt_id,mid,attr,ctime,mtime FROM report_user WHERE rpt_id in (%s)"
	_reportUserByRptIDSQL  = "SELECT id,rpt_id,mid,attr,ctime,mtime FROM report_user WHERE rpt_id=? ; "
	_reportLogSQL          = "SELECT id,username,points,oid,type,mid,tid,rid,reason,handle_type,notice,rpt,ctime,mtime FROM report_log WHERE rpt=?"
	_reportLogCountSQL     = "SELECT count(*) FROM report_log r %s ORDER BY r.id DESC;"
	_reportLogListSQL      = "SELECT r.id,r.username,r.points,r.oid,r.type,r.mid,r.tid,r.rid,r.reason,r.handle_type,r.notice,r.rpt,r.ctime,r.mtime FROM report_log r %s ORDER BY r.id DESC LIMIT ?,?"
	// _reportUserListSQL     = "SELECT id,rpt_id,mid,attr,ctime,mtime FROM report_user %s ORDER BY ctime DESC"
	// _reportCountSQL        = "SELECT count(*) FROM report rpt %s ORDER BY %s DESC"
	// _reportListSQL       = "SELECT id,oid,type,tid,mid,action,rid,count,reason,content,moral,score,state,ctime,mtime FROM report as rpt %s ORDER BY %s DESC LIMIT ?,?"
	// _reportFirstUserSQL  = "SELECT id,rpt_id,mid,attr,ctime,mtime FROM report_user WHERE rpt_id in (%s) AND attr in (1,3)"
	_reportInfoSQL       = "SELECT rpt.id,rpt.oid,rpt.type,rpt.tid,rpt.mid,rpt.action,rpt.rid,rpt.count,rpt.reason,rpt.moral,rpt.score,rpt.state,rpt.ctime,rpt.mtime,ru.mid FROM report rpt LEFT OUTER JOIN report_user ru ON rpt.id=ru.rpt_id WHERE rpt.state=? AND ru.attr in (1,3) %s ORDER BY %s DESC LIMIT ?,?"
	_reportInfoCountSQL  = "SELECT count(*) FROM report rpt LEFT OUTER JOIN report_user ru ON rpt.id=ru.rpt_id WHERE rpt.state=? AND ru.attr in (1,3) %s ORDER BY %s DESC"
	_reportLogByRptIDSQL = "SELECT id,rpt,username,points,oid,type,mid,tid,rid,reason,handle_type,notice,ctime,mtime FROM report_log  WHERE rpt in (%s) order by ctime DESC"
)

// Report Report.
func (d *Dao) Report(c context.Context, id int64) (res *model.Report, err error) {
	res = new(model.Report)
	row := d.db.QueryRow(c, _reportSQL, id)
	if err = row.Scan(&res.ID, &res.Oid, &res.Type, &res.Tid, &res.Mid, &res.Action, &res.Prid, &res.Rid, &res.Count, &res.Reason, &res.Content, &res.Moral, &res.Score, &res.State, &res.CTime, &res.MTime); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Error("row.Scan() error(%v)", err)
	}
	return
}

// Reports Reports.
func (d *Dao) Reports(c context.Context, ids []int64) (rpt []*model.Report, rptMap map[int64]*model.Report, rptIDs []int64, err error) {
	rptMap = make(map[int64]*model.Report)
	rows, err := d.db.Query(c, fmt.Sprintf(_reportsSQL, xstr.JoinInts(ids)))
	if err != nil {
		log.Error("query reports(%v), err error(%v)", ids, err)
		return
	}
	for rows.Next() {
		r := new(model.Report)
		if err = rows.Scan(&r.ID, &r.Oid, &r.Type, &r.Tid, &r.Mid, &r.Action, &r.Prid, &r.Rid, &r.Count, &r.Reason, &r.Content, &r.Moral, &r.Score, &r.State, &r.CTime, &r.MTime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		rpt = append(rpt, r)
		rptIDs = append(rptIDs, r.ID)
		rptMap[r.ID] = r
	}
	return
}

// ReportUser ReportUser.
func (d *Dao) ReportUser(c context.Context, id int64) (res *model.ReportUser, err error) {
	res = new(model.ReportUser)
	row := d.db.QueryRow(c, _reportUserByRptIDSQL, id)
	if err = row.Scan(&res.ID, &res.RptID, &res.Mid, &res.Attr, &res.CTime, &res.MTime); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Error("row.Scan() error(%v)", err)
	}
	return
}

// ReportUsers ReportUsers.
func (d *Dao) ReportUsers(c context.Context, ids []int64) (users []*model.ReportUser, userMap map[int64]*model.ReportUser, err error) {
	userMap = make(map[int64]*model.ReportUser)
	rows, err := d.db.Query(c, fmt.Sprintf(_reportUsersSQL, xstr.JoinInts(ids)))
	if err != nil {
		log.Error("query report user by report ids(%v) error(%v)", ids, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		u := &model.ReportUser{}
		if err = rows.Scan(&u.ID, &u.RptID, &u.Mid, &u.Attr, &u.CTime, &u.MTime); err != nil {
			log.Error("row.Scan() error(%v)", err)
			return
		}
		users = append(users, u)
		userMap[u.RptID] = u
	}
	return
}

// UpReportState update state by ID.
func (d *Dao) UpReportState(c context.Context, id int64, state int32) (affect int64, err error) {
	res, err := d.db.Exec(c, _updateReportStateSQL, state, id)
	if err != nil {
		log.Error("update report state(%d,%d) error(%v)", state, id, err)
		return
	}
	return res.RowsAffected()
}

// UpReportsState update state by ID.
func (d *Dao) UpReportsState(c context.Context, ids []int64, state int32) (affect int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_updateReportStatesSQL, xstr.JoinInts(ids)), state)
	if err != nil {
		log.Error("update reports state(%d,%+v) error(%v)", state, ids, err)
		return
	}
	return res.RowsAffected()
}

// AddReportLog add reportLog.
func (d *Dao) AddReportLog(c context.Context, l *model.ReportLog) (id int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_insertReportLogSQL, l.RptID, l.UserName, l.Points, l.Oid, l.Type, l.Mid, l.Tid, l.Rid, l.Reason, l.HandleType, l.Notice))
	if err != nil {
		log.Error("add report log(%+v) error(%v)", l, err)
		return
	}
	return res.LastInsertId()
}

// AddReportLogs add report logs
func (d *Dao) AddReportLogs(c context.Context, sqls []string) (id int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_insertReportLogsSQL, strings.Join(sqls, " , ")))
	if err != nil {
		log.Error("add report logs(%+v) error(%v)", sqls, err)
		return
	}
	return res.LastInsertId()
}

// UpdateMorals update moral.
func (d *Dao) UpdateMorals(c context.Context, moral int32, ids []int64) (affect int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_updateReportMoralSQL, xstr.JoinInts(ids)), moral)
	if err != nil {
		log.Error("update moral(%d,%d) error(%v)", moral, ids, err)
		return
	}
	return res.RowsAffected()
}

// ReportByOidMid get more report info by oid,aid,oid_type.
func (d *Dao) ReportByOidMid(c context.Context, mid, oid int64, tp int32) (rpt []*model.Report, rptIDs, tids []int64, err error) {
	rows, err := d.db.Query(c, _reportByOidMidSQL, mid, oid, tp)
	if err != nil {
		log.Error("query report by oid,mid(%d,%d,%d) error(%v)", mid, oid, tp, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.Report{}
		if err = rows.Scan(&r.ID, &r.Oid, &r.Type, &r.Tid, &r.Mid, &r.Action, &r.Prid, &r.Rid, &r.Count, &r.Reason, &r.Content, &r.Moral, &r.State, &r.CTime, &r.MTime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		rpt = append(rpt, r)
		rptIDs = append(rptIDs, r.ID)
		tids = append(tids, r.Tid)
	}
	return
}

// ReportLog ReportLog.
func (d *Dao) ReportLog(c context.Context, id int64) (res []*model.ReportLog, tids []int64, err error) {
	rows, err := d.db.Query(c, _reportLogSQL, id)
	if err != nil {
		log.Error("query reportlog by id(%d), error(%v)", id, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.ReportLog{}
		if err = rows.Scan(&r.ID, &r.UserName, &r.Points, &r.Oid, &r.Type, &r.Mid, &r.Tid, &r.Rid, &r.Reason, &r.HandleType, &r.Notice, &r.RptID, &r.CTime, &r.MTime); err != nil {
			log.Error("row.Scan report log scan() error(%v)", err)
			return
		}
		res = append(res, r)
		tids = append(tids, r.Tid)
	}
	return
}

// ReportLogCount count of report log.
func (d *Dao) ReportLogCount(c context.Context, sql string) (count int64, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_reportLogCountSQL, sql))
	if err = row.Scan(&count); err != nil {
		log.Error("Count Notice row.Scan(%s) error(%v)", sql, err)
	}
	return
}

// ReportLogList report log list.
// TODO when two table left out join, once Columns = nil, will break
func (d *Dao) ReportLogList(c context.Context, sql string, start, end int32) (res []*model.ReportLog, tids []int64, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_reportLogListSQL, sql), start, end)
	if err != nil {
		log.Error("query report log list(%v,%d,%d) error(%v)", sql, start, end, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		l := &model.ReportLog{}
		if err = rows.Scan(&l.ID, &l.UserName, &l.Points, &l.Oid, &l.Type, &l.Mid, &l.Tid, &l.Rid, &l.Reason, &l.HandleType, &l.Notice, &l.RptID, &l.CTime, &l.MTime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			err = nil
			return
		}
		res = append(res, l)
		tids = append(tids, l.Tid)
	}
	return
}

// ReportLogByRptID report log by rpt ids.
func (d *Dao) ReportLogByRptID(c context.Context, rptIDs []int64) (logs []*model.ReportLog, logMap map[int64][]*model.ReportLog, err error) {
	logMap = make(map[int64][]*model.ReportLog)
	rows, err := d.db.Query(c, fmt.Sprintf(_reportLogByRptIDSQL, xstr.JoinInts(rptIDs)))
	if err != nil {
		log.Error("query report log by rpt ids(%v) error(%v)", rptIDs, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		l := &model.ReportLog{}
		if err = rows.Scan(&l.ID, &l.RptID, &l.UserName, &l.Points, &l.Oid, &l.Type, &l.Mid, &l.Tid, &l.Rid, &l.Reason, &l.HandleType, &l.Notice, &l.CTime, &l.MTime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		logs = append(logs, l)
		logMap[l.RptID] = append(logMap[l.RptID], l)
	}
	return
}

// // ReportUserList report user info.
// func (d *Dao) ReportUserList(c context.Context, sqlStr []string) (userMap map[int64]*model.ReportUser, rptIDs []int64, err error) {
// 	var (
// 		sql string
// 	)
// 	userMap = make(map[int64]*model.ReportUser)
// 	if len(sqlStr) > 0 {
// 		sql = fmt.Sprintf(" WHERE %s", strings.Join(sqlStr, " AND "))
// 	}
// 	rows, err := d.db.Query(c, fmt.Sprintf(_reportUserListSQL, sql))
// 	if err != nil {
// 		log.Error("d.db.Query(%v) error(%v)", sql, err)
// 		return
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		u := &model.ReportUser{}
// 		if err = rows.Scan(&u.ID, &u.RptID, &u.Mid, &u.Attr, &u.CTime, &u.MTime); err != nil {
// 			log.Error("rows.Scan() error(%v)", err)
// 			return
// 		}
// 		rptIDs = append(rptIDs, u.RptID)
// 		userMap[u.RptID] = u
// 	}
// 	return
// }

// // ReportCount count report info.
// func (d *Dao) ReportCount(c context.Context, sqlStr []string, order string) (count int64, err error) {
// 	var sql string
// 	if len(sqlStr) > 0 {
// 		sql = fmt.Sprintf(" WHERE %s", strings.Join(sqlStr, " AND "))
// 	}
// 	row := d.db.QueryRow(c, fmt.Sprintf(_reportCountSQL, sql, order))
// 	if err = row.Scan(&count); err != nil {
// 		log.Error("CountNotice row.Scan err (%v)", err)
// 	}
// 	return
// }

// // ReportList ReportList.
// func (d *Dao) ReportList(c context.Context, sqlStr []string, order string, start, end int32) (res []*model.Report, tagIDs []int64, rptIDs []int64, oids []int64, err error) {
// 	var sql string
// 	if len(sqlStr) > 0 {
// 		sql = fmt.Sprintf(" WHERE %s", strings.Join(sqlStr, " AND "))
// 	}
// 	rows, err := d.db.Query(c, fmt.Sprintf(_reportListSQL, sql, order), start, end)
// 	if err != nil {
// 		log.Error("d.db.Query(%s) error(%v)", sql, err)
// 		return
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		r := &model.Report{}
// 		if err = rows.Scan(&r.ID, &r.Oid, &r.Type, &r.Tid, &r.Mid, &r.Action, &r.Rid, &r.Count, &r.Reason, &r.Content, &r.Moral, &r.Moral, &r.State, &r.CTime, &r.MTime); err != nil {
// 			log.Error("rows.Scan() error(%v)", err)
// 			return
// 		}
// 		res = append(res, r)
// 		rptIDs = append(rptIDs, r.ID)
// 		tagIDs = append(tagIDs, r.Tid)
// 		oids = append(oids, r.Oid)
// 	}
// 	return
// }

// ReportCount count report info.
func (d *Dao) ReportCount(c context.Context, sqlStr []string, order string, state int32) (count int64, err error) {
	var sql string
	if len(sqlStr) > 0 {
		sql = fmt.Sprintf(" AND %s", strings.Join(sqlStr, " AND "))
	}
	row := d.db.QueryRow(c, fmt.Sprintf(_reportInfoCountSQL, sql, order), state)
	if err = row.Scan(&count); err != nil {
		log.Error("CountNotice row.Scan err (%v)", err)
	}
	return
}

// ReportInfoList ReportInfoList.
func (d *Dao) ReportInfoList(c context.Context, sqlStr []string, order string, state, start, end int32) (res []*model.ReportInfo, tagIDs []int64, rptIDs []int64, oids []int64, err error) {
	var sql string
	if len(sqlStr) > 0 {
		sql = fmt.Sprintf(" AND %s", strings.Join(sqlStr, " AND "))
	}
	rows, err := d.db.Query(c, fmt.Sprintf(_reportInfoSQL, sql, order), state, start, end)
	if err != nil {
		log.Error("d.db.Query(%s) error(%v)", sql, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.ReportInfo{}
		if err = rows.Scan(&r.ID, &r.Oid, &r.Type, &r.Tid, &r.Mid, &r.Action, &r.Rid, &r.Count, &r.Reason, &r.IsDelMoral, &r.Score, &r.State, &r.CTime, &r.MTime, &r.RptMid); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		res = append(res, r)
		rptIDs = append(rptIDs, r.ID)
		tagIDs = append(tagIDs, r.Tid)
		oids = append(oids, r.Oid)
	}
	return
}

// // ReportFirstUsers  get first report user info.
// func (d *Dao) ReportFirstUsers(c context.Context, rptIDs []int64) (userMap map[int64]*model.ReportUser, err error) {
// 	userMap = make(map[int64]*model.ReportUser)
// 	sql := xstr.JoinInts(rptIDs)
// 	rows, err := d.db.Query(c, fmt.Sprintf(_reportFirstUserSQL, sql))
// 	if err != nil {
// 		log.Error("d.db.Query(%v) error(%v)", sql, err)
// 		return
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		u := &model.ReportUser{}
// 		if err = rows.Scan(&u.ID, &u.RptID, &u.Mid, &u.Attr, &u.CTime, &u.MTime); err != nil {
// 			log.Error("rows.Scan() error(%v)", err)
// 			return
// 		}
// 		userMap[u.RptID] = u
// 	}
// 	return
// }
