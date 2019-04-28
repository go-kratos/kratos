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
	_subSharding int64 = 50
)

const (
	_inSubSQL       = "INSERT IGNORE INTO reply_subject_%d (oid,type,mid,state,ctime,mtime) VALUES(?,?,?,?,?,?)"
	_upSubStateSQL  = "UPDATE reply_subject_%d SET state=?,mtime=? WHERE oid=? AND type=?"
	_upSubMidSQL    = "UPDATE reply_subject_%d SET mid=?,mtime=? WHERE oid=? AND type=?"
	_selSubjectSQL  = "SELECT oid,type,mid,count,rcount,acount,state,attr,ctime,mtime,meta FROM reply_subject_%d WHERE oid=? AND type=?"
	_selSubjectsSQL = "SELECT oid,type,mid,count,rcount,acount,state,attr,ctime,mtime,meta FROM reply_subject_%d WHERE type=? AND oid IN(%s)"
	_setSubjectSQL  = "INSERT INTO reply_subject_%d (oid,type,mid,state,ctime,mtime) VALUES(?,?,?,?,?,?) ON DUPLICATE KEY UPDATE state=?,mid=?,mtime=?"
)

// SubjectDao subject dao.
type SubjectDao struct {
	db *sql.DB
}

// NewSubjectDao new ReplySubjectDao and return.
func NewSubjectDao(db *sql.DB) (dao *SubjectDao) {
	dao = &SubjectDao{
		db: db,
	}
	return
}

func (dao *SubjectDao) hit(oid int64) int64 {
	return oid % _subSharding
}

// Set insert or update subject state.
func (dao *SubjectDao) Set(c context.Context, sub *reply.Subject) (id int64, err error) {
	res, err := dao.db.Exec(c, fmt.Sprintf(_setSubjectSQL, dao.hit(sub.Oid)), sub.Oid, sub.Type, sub.Mid, sub.State, sub.CTime, sub.MTime, sub.State, sub.Mid, sub.MTime)
	if err != nil {
		log.Error("mysqlDB.Exec error(%v)", err)
		return
	}
	return res.LastInsertId()
}

// Insert insert a reply subject.
func (dao *SubjectDao) Insert(c context.Context, sub *reply.Subject) (id int64, err error) {
	res, err := dao.db.Exec(c, fmt.Sprintf(_inSubSQL, dao.hit(sub.Oid)), sub.Oid, sub.Type, sub.Mid, sub.State, sub.CTime, sub.MTime)
	if err != nil {
		log.Error("mysqlDB.Exec error(%v)", err)
		return
	}
	return res.LastInsertId()
}

// UpState update subject state.
func (dao *SubjectDao) UpState(c context.Context, oid int64, tp int8, state int8, now time.Time) (rows int64, err error) {
	res, err := dao.db.Exec(c, fmt.Sprintf(_upSubStateSQL, dao.hit(oid)), state, now, oid, tp)
	if err != nil {
		log.Error("mysqlDB.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// UpMid update subject mid.
func (dao *SubjectDao) UpMid(c context.Context, mid, oid int64, tp int8, now time.Time) (rows int64, err error) {
	res, err := dao.db.Exec(c, fmt.Sprintf(_upSubMidSQL, dao.hit(oid)), mid, now, oid, tp)
	if err != nil {
		log.Error("mysqlDB.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// Get get a subject.
func (dao *SubjectDao) Get(c context.Context, oid int64, tp int8) (sub *reply.Subject, err error) {
	sub = &reply.Subject{}
	row := dao.db.QueryRow(c, fmt.Sprintf(_selSubjectSQL, dao.hit(oid)), oid, tp)
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

// Gets get a subject.
func (dao *SubjectDao) Gets(c context.Context, oids []int64, tp int8) (res map[int64]*reply.Subject, err error) {
	hits := make(map[int64][]int64)
	for _, oid := range oids {
		hit := dao.hit(oid)
		hits[hit] = append(hits[hit], oid)
	}
	res = make(map[int64]*reply.Subject, len(oids))
	for hit, oids := range hits {
		var rows *sql.Rows
		if rows, err = dao.db.Query(c, fmt.Sprintf(_selSubjectsSQL, hit, xstr.JoinInts(oids)), tp); err != nil {
			log.Error("dao.db.Query error(%v)", err)
			return
		}
		for rows.Next() {
			sub := new(reply.Subject)
			if err = rows.Scan(&sub.Oid, &sub.Type, &sub.Mid, &sub.Count, &sub.RCount, &sub.ACount, &sub.State, &sub.Attr, &sub.CTime, &sub.MTime, &sub.Meta); err != nil {
				if err == sql.ErrNoRows {
					sub = nil
					err = nil
					continue
				} else {
					log.Error("row.Scan error(%v)", err)
					rows.Close()
					return
				}
			}
			res[sub.Oid] = sub
		}
		if err = rows.Err(); err != nil {
			log.Error("rows.err error(%v)", err)
			rows.Close()
			return
		}
		rows.Close()
	}
	return
}
