package reply

import (
	"context"
	"fmt"

	model "go-common/app/interface/main/reply/model/reply"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_repSharding int64 = 200
)

const (
	// report
	_inRptSQL  = "INSERT INTO reply_report_%d (oid,type,rpid,mid,reason,content,count,score,state,ctime,mtime) VALUES(?,?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE count=count+1, score=score+?, mtime=?"
	_selRptSQL = "SELECT oid,type,rpid,mid,reason,content,count,state,ctime,mtime,attr FROM reply_report_%d WHERE rpid=?"
	// report user
	_inRptUserSQL = "INSERT IGNORE INTO reply_report_user_%d (oid,type,rpid,mid,reason,content,state,ctime,mtime) VALUES(?,?,?,?,?,?,?,?,?)"
)

// ReportDao report dao.
type ReportDao struct {
	inStmts     []*sql.Stmt
	inUserStmts []*sql.Stmt
	db          *sql.DB
}

// NewReportDao new ReplyReportDao and return.
func NewReportDao(db *sql.DB) (dao *ReportDao) {
	dao = &ReportDao{
		db:          db,
		inStmts:     make([]*sql.Stmt, _repSharding),
		inUserStmts: make([]*sql.Stmt, _repSharding),
	}
	for i := int64(0); i < _repSharding; i++ {
		dao.inStmts[i] = dao.db.Prepared(fmt.Sprintf(_inRptSQL, i))
		dao.inUserStmts[i] = dao.db.Prepared(fmt.Sprintf(_inRptUserSQL, i))
	}
	return
}

func (dao *ReportDao) hit(oid int64) int64 {
	return oid % _repSharding
}

// Insert insert a report to mysql.
func (dao *ReportDao) Insert(c context.Context, rpt *model.Report) (id int64, err error) {
	res, err := dao.inStmts[dao.hit(rpt.Oid)].Exec(c, rpt.Oid, rpt.Type, rpt.RpID, rpt.Mid, rpt.Reason, rpt.Content, rpt.Count, rpt.Score, rpt.State, rpt.CTime, rpt.MTime, rpt.Score, rpt.MTime)
	if err != nil {
		log.Error("db.Exec error(%v)", err)
		return
	}
	return res.LastInsertId()
}

// Get return a report from mysql.
func (dao *ReportDao) Get(c context.Context, oid, rpID int64) (rpt *model.Report, err error) {
	row := dao.db.QueryRow(c, fmt.Sprintf(_selRptSQL, dao.hit(oid)), rpID)
	rpt = &model.Report{}
	err = row.Scan(&rpt.Oid, &rpt.Type, &rpt.RpID, &rpt.Mid, &rpt.Reason, &rpt.Content, &rpt.Count, &rpt.State, &rpt.CTime, &rpt.MTime, &rpt.Attr)
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

// InsertUser inser a report user to mysql.
func (dao *ReportDao) InsertUser(c context.Context, ru *model.ReportUser) (id int64, err error) {
	res, err := dao.inUserStmts[dao.hit(ru.Oid)].Exec(c, ru.Oid, ru.Type, ru.RpID, ru.Mid, ru.Reason, ru.Content, ru.State, ru.CTime, ru.MTime)
	if err != nil {
		log.Error("db.Exec error(%v)", err)
		return
	}
	return res.LastInsertId()
}
