package reply

import (
	"context"
	"fmt"

	"go-common/app/job/main/reply/model/reply"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_contSharding int64 = 200
)

const (
	_inContSQL  = "INSERT IGNORE INTO reply_content_%d (rpid,message,ats,ip,plat,device,version,ctime,mtime,topics) VALUES(?,?,?,?,?,?,?,?,?,?)"
	_selContSQL = "SELECT rpid,message,ats,ip,plat,device,topics FROM reply_content_%d WHERE rpid=?"
)

// ContentDao define content mysql stmt
type ContentDao struct {
	selContStmts []*sql.Stmt
	mysql        *sql.DB
}

// NewContentDao new contentDao and return.
func NewContentDao(db *sql.DB) (dao *ContentDao) {
	dao = &ContentDao{
		mysql:        db,
		selContStmts: make([]*sql.Stmt, _contSharding),
	}
	for i := int64(0); i < _contSharding; i++ {
		dao.selContStmts[i] = dao.mysql.Prepared(fmt.Sprintf(_selContSQL, i))
	}
	return
}

func (dao *ContentDao) hit(oid int64) int64 {
	return oid % int64(_contSharding)
}

// TxInsert insert reply content by transaction.
func (dao *ContentDao) TxInsert(tx *sql.Tx, oid int64, rc *reply.Content) (rows int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_inContSQL, dao.hit(oid)), rc.RpID, rc.Message, rc.Ats, rc.IP, rc.Plat, rc.Device, rc.Version, rc.CTime, rc.MTime, rc.Topics)
	if err != nil {
		log.Error("mysqlDB.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// Get get reply content.
func (dao *ContentDao) Get(c context.Context, oid int64, rpID int64) (rc *reply.Content, err error) {
	row := dao.selContStmts[dao.hit(oid)].QueryRow(c, rpID)
	rc = &reply.Content{}
	if err = row.Scan(&rc.RpID, &rc.Message, &rc.Ats, &rc.IP, &rc.Plat, &rc.Device, &rc.Topics); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			rc = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}
