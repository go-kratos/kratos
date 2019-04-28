package reply

import (
	"context"
	"fmt"
	"time"

	"go-common/app/interface/main/reply/model/reply"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_contSharding int64 = 200
)

const (
	_selContSQL   = "SELECT rpid,message,ats,ip,plat,device,topics FROM reply_content_%d WHERE rpid=?"
	_selContsSQL  = "SELECT rpid,message,ats,ip,plat,device,topics FROM reply_content_%d WHERE rpid IN (%s)"
	_upContMsgSQL = "UPDATE reply_content_%d SET message=?,mtime=? WHERE rpid=?"
)

// ContentDao ContentDao
type ContentDao struct {
	upContMsgStmts []*sql.Stmt
	db             *sql.DB
	dbSlave        *sql.DB
}

// NewContentDao new replyDao and return.
func NewContentDao(db *sql.DB, dbSlave *sql.DB) (dao *ContentDao) {
	dao = &ContentDao{
		db:             db,
		dbSlave:        dbSlave,
		upContMsgStmts: make([]*sql.Stmt, _contSharding),
	}
	for i := int64(0); i < _contSharding; i++ {
		dao.upContMsgStmts[i] = dao.db.Prepared(fmt.Sprintf(_upContMsgSQL, i))
	}
	return
}

func (dao *ContentDao) hit(oid int64) int64 {
	return oid % int64(_contSharding)
}

// UpMessage update content's message.
func (dao *ContentDao) UpMessage(c context.Context, oid int64, rpID int64, msg string, now time.Time) (rows int64, err error) {
	res, err := dao.upContMsgStmts[dao.hit(oid)].Exec(c, msg, now, rpID)
	if err != nil {
		log.Error("contentDao.UpMessage error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// Get get reply content.
func (dao *ContentDao) Get(c context.Context, oid int64, rpID int64) (rc *reply.Content, err error) {
	row := dao.db.QueryRow(c, fmt.Sprintf(_selContSQL, dao.hit(oid)), rpID)
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

// GetByIds get reply contents by reply ids.
func (dao *ContentDao) GetByIds(c context.Context, oid int64, rpIds []int64) (rcMap map[int64]*reply.Content, err error) {
	if len(rpIds) == 0 {
		return
	}
	rows, err := dao.dbSlave.Query(c, fmt.Sprintf(_selContsSQL, dao.hit(oid), xstr.JoinInts(rpIds)))
	if err != nil {
		log.Error("contentDao.Query error(%v)", err)
		return
	}
	defer rows.Close()
	rcMap = make(map[int64]*reply.Content, len(rpIds))
	for rows.Next() {
		rc := &reply.Content{}
		if err = rows.Scan(&rc.RpID, &rc.Message, &rc.Ats, &rc.IP, &rc.Plat, &rc.Device, &rc.Topics); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		rcMap[rc.RpID] = rc
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.err error(%v)", err)
		return
	}
	return
}
