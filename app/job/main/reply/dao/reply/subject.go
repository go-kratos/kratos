package reply

import (
	"context"
	"fmt"
	"time"

	"go-common/app/job/main/reply/model/reply"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_subSharding int64 = 50
)

const (
	_incrSubCntSQL          = "UPDATE reply_subject_%d SET count=count+1,rcount=rcount+1,acount=acount+1,mtime=? WHERE oid=? AND type=?"
	_incrSubFCntSQL         = "UPDATE reply_subject_%d SET count=count+1,mtime=? WHERE oid=? AND type=?"
	_incrSubRCntSQL         = "UPDATE reply_subject_%d SET rcount=rcount+1,mtime=? WHERE oid=? AND type=?"
	_incrSubACntSQL         = "UPDATE reply_subject_%d SET acount=acount+?,mtime=? WHERE oid=? AND type=?"
	_incrSubMCntSQL         = "UPDATE reply_subject_%d SET mcount=mcount+1,mtime=? WHERE oid=? AND type=?"
	_decrSubMCntSQL         = "UPDATE reply_subject_%d SET mcount=mcount-1,mtime=? WHERE oid=? AND type=? AND mcount>0"
	_decrSubCntSQL          = "UPDATE reply_subject_%d SET rcount=rcount-1,mtime=? WHERE oid=? AND type=?"
	_upSubAttrSQL           = "UPDATE reply_subject_%d SET attr=?,mtime=? WHERE oid=? AND type=?"
	_decrSubACntSQL         = "UPDATE reply_subject_%d SET acount=acount-?,mtime=? WHERE oid=? AND type=?"
	_upSubMetaSQL           = "UPDATE reply_subject_%d SET meta=?,mtime=? WHERE oid=? AND type=?"
	_selSubjectSQL          = "SELECT oid,type,mid,count,rcount,acount,state,attr,ctime,mtime,meta FROM reply_subject_%d WHERE oid=? AND type=?"
	_selSubjectForUpdateSQL = "SELECT oid,type,mid,count,rcount,acount,state,attr,ctime,mtime,meta FROM reply_subject_%d WHERE oid=? AND type=? FOR UPDATE"
)

// SubjectDao define subject mysql stmt
type SubjectDao struct {
	selSubjectStmt []*sql.Stmt
	mysql          *sql.DB
}

// NewSubjectDao new ReplySubjectDao and return.
func NewSubjectDao(db *sql.DB) (dao *SubjectDao) {

	dao = &SubjectDao{
		mysql:          db,
		selSubjectStmt: make([]*sql.Stmt, _subSharding),
	}
	for i := int64(0); i < _subSharding; i++ {
		dao.selSubjectStmt[i] = dao.mysql.Prepared(fmt.Sprintf(_selSubjectSQL, i))
	}
	return
}

func (dao *SubjectDao) hit(oid int64) int64 {
	return oid % _subSharding
}

// UpMeta update subject meta.
func (dao *SubjectDao) UpMeta(c context.Context, oid int64, tp int8, meta string, now time.Time) (rows int64, err error) {
	res, err := dao.mysql.Exec(c, fmt.Sprintf(_upSubMetaSQL, dao.hit(oid)), meta, now, oid, tp)
	if err != nil {
		log.Error("mysqlDB.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// TxUpMeta update subject meta.
func (dao *SubjectDao) TxUpMeta(tx *sql.Tx, oid int64, tp int8, meta string, now time.Time) (rows int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_upSubMetaSQL, dao.hit(oid)), meta, now, oid, tp)
	if err != nil {
		log.Error("mysqlDB.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// TxUpAttr update subject attr.
func (dao *SubjectDao) TxUpAttr(tx *sql.Tx, oid int64, tp int8, attr uint32, now time.Time) (rows int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_upSubAttrSQL, dao.hit(oid)), attr, now, oid, tp)
	if err != nil {
		log.Error("mysqlDB.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// TxIncrCount incr subject count and rcount by transaction.
func (dao *SubjectDao) TxIncrCount(tx *sql.Tx, oid int64, tp int8, now time.Time) (rows int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_incrSubCntSQL, dao.hit(oid)), now, oid, tp)
	if err != nil {
		log.Error("mysqlDB.Exec() error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// TxIncrFCount incr subject count and rcount by transaction.
func (dao *SubjectDao) TxIncrFCount(tx *sql.Tx, oid int64, tp int8, now time.Time) (rows int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_incrSubFCntSQL, dao.hit(oid)), now, oid, tp)
	if err != nil {
		log.Error("mysqlDB.Exec() error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// TxIncrMCount incr subject mcount by transaction.
func (dao *SubjectDao) TxIncrMCount(tx *sql.Tx, oid int64, tp int8, now time.Time) (rows int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_incrSubMCntSQL, dao.hit(oid)), now, oid, tp)
	if err != nil {
		log.Error("mysqlDB.Exec() error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// TxDecrMCount decr subject mcount by transaction.
func (dao *SubjectDao) TxDecrMCount(tx *sql.Tx, oid int64, tp int8, now time.Time) (rows int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_decrSubMCntSQL, dao.hit(oid)), now, oid, tp)
	if err != nil {
		log.Error("mysqlDB.Exec() error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// TxIncrRCount incr subject rcount by transaction
func (dao *SubjectDao) TxIncrRCount(tx *sql.Tx, oid int64, tp int8, now time.Time) (rows int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_incrSubRCntSQL, dao.hit(oid)), now, oid, tp)
	if err != nil {
		log.Error("mysqlDB.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// TxDecrCount decr subject count by transaction.
func (dao *SubjectDao) TxDecrCount(tx *sql.Tx, oid int64, tp int8, now time.Time) (rows int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_decrSubCntSQL, dao.hit(oid)), now, oid, tp)
	if err != nil {
		log.Error("mysqlDB.Exec() error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// TxIncrACount incr subject acount by transaction.
func (dao *SubjectDao) TxIncrACount(tx *sql.Tx, oid int64, tp int8, count int, now time.Time) (rows int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_incrSubACntSQL, dao.hit(oid)), count, now, oid, tp)
	if err != nil {
		log.Error("mysqlDB.Exec() error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// TxDecrACount decr subject rcount by transaction.
func (dao *SubjectDao) TxDecrACount(tx *sql.Tx, oid int64, tp int8, count int, now time.Time) (rows int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_decrSubACntSQL, dao.hit(oid)), count, now, oid, tp)
	if err != nil {
		log.Error("mysqlDB.Exec() error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// Get get a subject.
func (dao *SubjectDao) Get(c context.Context, oid int64, tp int8) (sub *reply.Subject, err error) {
	sub = &reply.Subject{}
	row := dao.selSubjectStmt[dao.hit(oid)].QueryRow(c, oid, tp)
	if err = row.Scan(&sub.Oid, &sub.Type, &sub.Mid, &sub.Count, &sub.RCount, &sub.ACount, &sub.State, &sub.Attr, &sub.CTime, &sub.MTime, &sub.Meta); err != nil {
		if err == sql.ErrNoRows {
			sub = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// GetForUpdate get a subject for update.
func (dao *SubjectDao) GetForUpdate(tx *sql.Tx, oid int64, tp int8) (sub *reply.Subject, err error) {
	sub = &reply.Subject{}
	if err = tx.QueryRow(fmt.Sprintf(_selSubjectForUpdateSQL, dao.hit(oid)), oid, tp).Scan(&sub.Oid, &sub.Type, &sub.Mid, &sub.Count, &sub.RCount, &sub.ACount, &sub.State, &sub.Attr, &sub.CTime, &sub.MTime, &sub.Meta); err != nil {
		return
	}
	return
}
