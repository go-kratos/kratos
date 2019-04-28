package reply

import (
	"context"
	"fmt"
	"time"

	model "go-common/app/job/main/reply/model/reply"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_repSharding int64 = 200
)

const ( // report
	_upRepSQL       = "UPDATE reply_report_%d SET state=?,mtime=?,reason=?,content=?,attr=? WHERE rpid=?"
	_selRepSQL      = "SELECT oid,type,rpid,mid,reason,content,count,score,state,ctime,mtime,attr FROM reply_report_%d WHERE rpid=?"
	_selRepByOidSQL = "SELECT oid,type,rpid,mid,reason,content,count,score,state,ctime,mtime,attr FROM reply_report_%d WHERE oid=? and type=?"
	// report user
	_getRptUsersSQL     = "SELECT oid,type,rpid,mid,reason,content,state,ctime,mtime FROM reply_report_user_%d WHERE rpid=? and state=?"
	_setRptUserStateSQL = "UPDATE reply_report_user_%d SET state=?,mtime=? WHERE rpid=?"
)

//ReportDao define report mysql stmt
type ReportDao struct {
	upRepStmts        []*sql.Stmt
	selRepStmts       []*sql.Stmt
	getUsersStmts     []*sql.Stmt
	setUserStateStmts []*sql.Stmt
	mysql             *sql.DB
}

// NewReportDao new ReplyReportDao and return.
func NewReportDao(db *sql.DB) (dao *ReportDao) {
	dao = &ReportDao{
		mysql:             db,
		upRepStmts:        make([]*sql.Stmt, _repSharding),
		selRepStmts:       make([]*sql.Stmt, _repSharding),
		getUsersStmts:     make([]*sql.Stmt, _repSharding),
		setUserStateStmts: make([]*sql.Stmt, _repSharding),
	}
	for i := int64(0); i < _repSharding; i++ {
		dao.upRepStmts[i] = dao.mysql.Prepared(fmt.Sprintf(_upRepSQL, i))
		dao.selRepStmts[i] = dao.mysql.Prepared(fmt.Sprintf(_selRepSQL, i))
		dao.getUsersStmts[i] = dao.mysql.Prepared(fmt.Sprintf(_getRptUsersSQL, i))
		dao.setUserStateStmts[i] = dao.mysql.Prepared(fmt.Sprintf(_setRptUserStateSQL, i))
	}
	return
}

func (dao *ReportDao) hit(oid int64) int64 {
	return oid % _repSharding
}

// Update update reply report.
func (dao *ReportDao) Update(c context.Context, rpt *model.Report) (rows int64, err error) {
	res, err := dao.upRepStmts[dao.hit(rpt.Oid)].Exec(c, rpt.State, rpt.MTime, rpt.Reason, rpt.Content, rpt.Attr, rpt.RpID)
	if err != nil {
		log.Error("mysqlDB.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// Get get a reply report.
func (dao *ReportDao) Get(c context.Context, oid, rpID int64) (rpt *model.Report, err error) {
	row := dao.selRepStmts[dao.hit(oid)].QueryRow(c, rpID)
	rpt = &model.Report{}
	err = row.Scan(&rpt.Oid, &rpt.Type, &rpt.RpID, &rpt.Mid, &rpt.Reason, &rpt.Content, &rpt.Count, &rpt.Score, &rpt.State, &rpt.CTime, &rpt.MTime, &rpt.Attr)
	if err != nil {
		if err == sql.ErrNoRows {
			rpt = nil
			err = nil
		} else {
			log.Error("Mysql error(%v)", err)
		}
	}
	return
}

// GetMapByOid return report map by oid.
func (dao *ReportDao) GetMapByOid(c context.Context, oid int64, typ int8) (res map[int64]*model.Report, err error) {
	rows, err := dao.mysql.Query(c, fmt.Sprintf(_selRepByOidSQL, dao.hit(oid)), oid, typ)
	if err != nil {
		log.Error("db.Query(%s) error(%v)", _selRepByOidSQL, err)
		return
	}
	defer rows.Close()
	res = make(map[int64]*model.Report)
	for rows.Next() {
		rpt := &model.Report{}
		if err = rows.Scan(&rpt.Oid, &rpt.Type, &rpt.RpID, &rpt.Mid, &rpt.Reason, &rpt.Content, &rpt.Count, &rpt.Score, &rpt.State, &rpt.CTime, &rpt.MTime, &rpt.Attr); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		res[rpt.RpID] = rpt
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.err error(%v)", err)
		return
	}
	return
}

// GetUsers return a report users from mysql.
func (dao *ReportDao) GetUsers(c context.Context, oid int64, tp int8, rpID int64) (res map[int64]*model.ReportUser, err error) {
	rows, err := dao.getUsersStmts[dao.hit(oid)].Query(c, rpID, model.ReportUserStateNew)
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
func (dao *ReportDao) SetUserReported(c context.Context, oid int64, tp int8, rpID int64, now time.Time) (rows int64, err error) {
	res, err := dao.setUserStateStmts[dao.hit(oid)].Exec(c, model.ReportUserStateReported, now, rpID)
	if err != nil {
		log.Error("db.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}
