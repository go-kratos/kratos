package reply

import (
	"context"
	"time"

	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_inAdminSQL = "INSERT INTO reply_admin_log (oid,type,rpID,adminid,result,remark,isnew,isreport,state,ctime,mtime) VALUES(?,?,?,?,?,?,?,?,?,?,?)"
	_upAdminSQL = "UPDATE reply_admin_log SET isnew=0,mtime=? WHERE rpID=? AND isnew=1"
)

// AdminDao define admin mysql info
type AdminDao struct {
	mysql *sql.DB
}

// NewAdminDao new ReplyReportDao and return.
func NewAdminDao(db *sql.DB) (dao *AdminDao) {
	dao = &AdminDao{
		mysql: db,
	}
	return
}

// Insert insert reply report.
func (dao *AdminDao) Insert(c context.Context, adminid, oid, rpID int64, tp int8, result, remark string, isnew, isreport, state int8, now time.Time) (id int64, err error) {
	res, err := dao.mysql.Exec(c, _inAdminSQL, oid, tp, rpID, adminid, result, remark, isnew, isreport, state, now, now)
	if err != nil {
		log.Error("mysqlDB.Exec error(%v)", err)
		return
	}
	return res.LastInsertId()
}

// UpIsNotNew update reply report.
func (dao *AdminDao) UpIsNotNew(c context.Context, rpID int64, now time.Time) (rows int64, err error) {
	res, err := dao.mysql.Exec(c, _upAdminSQL, now, rpID)
	if err != nil {
		log.Error("mysqlDB.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}
