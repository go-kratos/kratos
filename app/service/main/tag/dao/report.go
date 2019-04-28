package dao

import (
	"context"
	"fmt"

	"go-common/app/service/main/tag/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

var (
	_reportSQL = "SELECT id,oid,`type`,tid,mid,`action`,rid,count,reason,content,state,ctime,mtime FROM platform_tag.report where oid = ? AND `type` = ?;"
)

// Report select list from report  mysql.
func (d *Dao) Report(c context.Context, oid int64, typ int32) (res map[string]*model.Report, err error) {
	rows, err := d.db.Query(c, _reportSQL, oid, typ)
	if err != nil {
		log.Error("Report d.db.Query(%d,%d) err:%v", oid, typ, err)
		return
	}
	defer rows.Close()
	res = make(map[string]*model.Report)
	for rows.Next() {
		r := &model.Report{}
		if err = rows.Scan(&r.ID, &r.Oid, &r.Type, &r.Tid, &r.Mid, &r.Action, &r.TypeID, &r.Count, &r.Reason, &r.Content, &r.State, &r.CTime, &r.MTime); err != nil {
			log.Error("rows.Scan() err:%v", err)
			return
		}
		k := fmt.Sprintf("%d_%d_%d_%d_%d", r.Oid, r.Type, r.Tid, r.Mid, r.Action)
		res[k] = r
	}
	return
}

var (
	_addReportSQL = "INSERT INTO report (oid,type,tid,mid,action,rid,count,reason,content,state,score) VALUES (?,?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE count=count+1, score=score+?"
)

// TxAddReport add a report into mysql. typeid = prid
func (d *Dao) TxAddReport(tx *sql.Tx, rpt *model.Report) (rptID int64, err error) {
	res, err := tx.Exec(_addReportSQL, rpt.Oid, rpt.Type, rpt.Tid, rpt.Mid, rpt.Action, rpt.TypeID, rpt.Count, rpt.Reason, rpt.Content, rpt.State, rpt.Score, rpt.Score)
	if err != nil {
		log.Error("TxAddReport d.db.Exec() error(%v)", err)
		return
	}
	return res.LastInsertId()
}

var (
	_insertReportUserSQL = "INSERT INTO `platform_tag`.`report_user` (`rpt_id`, `mid`, `attr`) VALUES ( ?, ?, ?);"
)

// TxAddUserReport add a user report into mysql.
func (d *Dao) TxAddUserReport(tx *sql.Tx, r *model.ReportUser) (rptID int64, err error) {
	res, err := tx.Exec(_insertReportUserSQL, r.RptID, r.Mid, r.Attr)
	if err != nil {
		log.Error("d.db.Exec() error(%v)", err)
		return
	}
	return res.LastInsertId()
}

// AddUserReport add a user report into mysql.
func (d *Dao) AddUserReport(c context.Context, r *model.ReportUser) (int64, error) {
	res, err := d.db.Exec(c, _insertReportUserSQL, r.RptID, r.Mid, r.Attr)
	if err != nil {
		log.Error("d.db.Exec() error(%v)", err)
		return 0, err
	}
	return res.LastInsertId()
}

var (
	_selectReporAndtUser = "SELECT r.oid,r.`type`,r.tid,r.mid,r.`action` FROM platform_tag.report as r LEFT JOIN platform_tag.report_user as u ON r.id = u.rpt_id WHERE r.oid = ? AND r.`type` = ? AND r.tid = ? AND r.mid=? AND r.`action` = ? ORDER BY r.id DESC LIMIT 500 ; "
)

// ReportAndUser .
func (d *Dao) ReportAndUser(c context.Context, oid, mid, tid int64, typ, action int32) (res []*model.Report, err error) {
	rows, err := d.db.Query(c, _selectReporAndtUser, oid, typ, tid, mid, action)
	if err != nil {
		log.Error("d.db.Query(%d,%d,%d,%d,%d) err:%v", oid, mid, tid, typ, action, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.Report{}
		if err = rows.Scan(&r.Oid, &r.Type, &r.Tid, &r.Mid, &r.Action); err != nil {
			log.Error("rows.Scan() err:%v", err)
			return
		}
		res = append(res, r)
	}
	return
}

var (
	_selectRepoerUser = "SELECT mid FROM platform_tag.report_user as u where u.rpt_id = ?;"
)

// ReportUser .
func (d *Dao) ReportUser(c context.Context, lid int64) (res map[int64]bool, err error) {
	res = make(map[int64]bool)
	rows, err := d.db.Query(c, _selectRepoerUser, lid)
	if err != nil {
		log.Error("d.db.Query(%d) err:%v", lid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var mid int64
		if err = rows.Scan(&mid); err != nil {
			log.Error("rows.Scan() err:%v", err)
			return
		}
		res[mid] = true
	}
	return
}
