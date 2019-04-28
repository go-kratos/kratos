package dao

import (
	"context"
	"fmt"
	"time"

	"go-common/app/admin/main/reply/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	"go-common/library/xstr"
	"sync"
)

const _subShard int64 = 50

const (
	_selSubjectSQL = "SELECT id,oid,type,mid,count,rcount,acount,state,attr,ctime,mtime FROM reply_subject_%d WHERE oid=? AND type=?"

	_selSubjectForUpdateSQL = "SELECT id,oid,type,mid,count,rcount,acount,state,attr,ctime,mtime FROM reply_subject_%d WHERE oid=? AND type=? FOR UPDATE"
	_selSubjectsSQL         = "SELECT id,oid,type,mid,count,rcount,acount,state,attr,ctime,mtime FROM reply_subject_%d WHERE oid IN (%s) AND type=?"
	_inSubjectSQL           = "INSERT INTO reply_subject_%d (oid,type,mid,state,ctime,mtime) VALUES(?,?,?,?,?,?) ON DUPLICATE KEY UPDATE state=?,mid=?,mtime=?"
	_selSubjectMCountSQL    = "SELECT oid, mcount FROM reply_subject_%d WHERE oid IN(%s) AND type=?"
	_upSubStateSQL          = "UPDATE reply_subject_%d SET state=?,mtime=? WHERE oid=? AND type=?"
	_upSubAttrSQL           = "UPDATE reply_subject_%d SET attr=?,mtime=? WHERE oid=? AND type=?"
	_upSubStateAndAttrSQL   = "UPDATE reply_subject_%d SET state=?,attr=?,mtime=? WHERE oid=? AND type=?"
	_upSubMetaSQL           = "UPDATE reply_subject_%d SET meta=?,mtime=? WHERE oid=? AND type=?"

	_incrSubCountSQL  = "UPDATE reply_subject_%d SET count=count+1,rcount=rcount+1,acount=acount+1,mtime=? WHERE oid=? AND type=?"
	_incrSubFCountSQL = "UPDATE reply_subject_%d SET count=count+1,mtime=? WHERE oid=? AND type=?"
	_incrSubRCountSQL = "UPDATE reply_subject_%d SET rcount=rcount+1,mtime=? WHERE oid=? AND type=?"
	_incrSubACountSQL = "UPDATE reply_subject_%d SET acount=acount+?,mtime=? WHERE oid=? AND type=?"
	_decrSubRCountSQL = "UPDATE reply_subject_%d SET rcount=rcount-1,mtime=? WHERE oid=? AND type=?"
	_decrSubACountSQL = "UPDATE reply_subject_%d SET acount=acount-?,mtime=? WHERE oid=? AND type=?"
	_decrSubMCountSQL = "UPDATE reply_subject_%d SET mcount=mcount-1,mtime=? WHERE oid=? AND type=? AND mcount>0"
)

func subHit(id int64) int64 {
	return id % _subShard
}

// SubMCount get subject mcount from mysql
func (d *Dao) SubMCount(c context.Context, oids []int64, typ int32) (res map[int64]int32, err error) {
	hits := make(map[int64][]int64)
	for _, oid := range oids {
		hit := subHit(oid)
		hits[hit] = append(hits[hit], oid)
	}
	res = make(map[int64]int32, len(oids))
	wg, ctx := errgroup.WithContext(c)
	var lock = sync.RWMutex{}
	for idx, oids := range hits {
		o := oids
		i := idx
		wg.Go(func() (err error) {
			var rows *xsql.Rows
			if rows, err = d.db.Query(ctx, fmt.Sprintf(_selSubjectMCountSQL, i, xstr.JoinInts(o)), typ); err != nil {
				log.Error("dao.db.Query error(%v)", err)
				return
			}
			var mcount int32
			var oid int64
			for rows.Next() {
				if err = rows.Scan(&oid, &mcount); err != nil {
					if err == xsql.ErrNoRows {
						mcount = 0
						oid = 0
						err = nil
						continue
					} else {
						log.Error("row.Scan error(%v)", err)
						rows.Close()
						return
					}
				}
				lock.Lock()
				res[oid] = mcount
				lock.Unlock()
			}
			if err = rows.Err(); err != nil {
				log.Error("rows.err error(%v)", err)
				rows.Close()
				return
			}
			rows.Close()
			return
		})
	}
	if err = wg.Wait(); err != nil {
		return
	}
	return
}

// Subjects get  subjects from mysql.
func (d *Dao) Subjects(c context.Context, oids []int64, typ int32) (subMap map[int64]*model.Subject, err error) {
	hitMap := make(map[int64][]int64)
	for _, oid := range oids {
		hitMap[subHit(oid)] = append(hitMap[subHit(oid)], oid)
	}
	subMap = make(map[int64]*model.Subject)
	for hit, ids := range hitMap {
		var rows *xsql.Rows
		rows, err = d.db.Query(c, fmt.Sprintf(_selSubjectsSQL, hit, xstr.JoinInts(ids)), typ)
		if err != nil {
			return
		}
		for rows.Next() {
			m := new(model.Subject)
			if err = rows.Scan(&m.ID, &m.Oid, &m.Type, &m.Mid, &m.Count, &m.RCount, &m.ACount, &m.State, &m.Attr, &m.CTime, &m.MTime); err != nil {
				rows.Close()
				return
			}
			subMap[m.Oid] = m
		}
		if err = rows.Err(); err != nil {
			rows.Close()
			return
		}
		rows.Close()
	}
	return
}

// Subject get a subject from mysql.
func (d *Dao) Subject(c context.Context, oid int64, typ int32) (m *model.Subject, err error) {
	m = new(model.Subject)
	row := d.db.QueryRow(c, fmt.Sprintf(_selSubjectSQL, subHit(oid)), oid, typ)
	if err = row.Scan(&m.ID, &m.Oid, &m.Type, &m.Mid, &m.Count, &m.RCount, &m.ACount, &m.State, &m.Attr, &m.CTime, &m.MTime); err != nil {
		if err == xsql.ErrNoRows {
			m = nil
			err = nil
		}
	}
	return
}

// TxSubject get a subject from mysql.
func (d *Dao) TxSubject(tx *xsql.Tx, oid int64, typ int32) (m *model.Subject, err error) {
	m = new(model.Subject)
	row := tx.QueryRow(fmt.Sprintf(_selSubjectSQL, subHit(oid)), oid, typ)
	if err = row.Scan(&m.ID, &m.Oid, &m.Type, &m.Mid, &m.Count, &m.RCount, &m.ACount, &m.State, &m.Attr, &m.CTime, &m.MTime); err != nil {
		if err == xsql.ErrNoRows {
			m = nil
			err = nil
		}
	}
	return
}

// TxSubjectForUpdate get a subject from mysql for update.
func (d *Dao) TxSubjectForUpdate(tx *xsql.Tx, oid int64, typ int32) (m *model.Subject, err error) {
	m = new(model.Subject)
	row := tx.QueryRow(fmt.Sprintf(_selSubjectForUpdateSQL, subHit(oid)), oid, typ)
	if err = row.Scan(&m.ID, &m.Oid, &m.Type, &m.Mid, &m.Count, &m.RCount, &m.ACount, &m.State, &m.Attr, &m.CTime, &m.MTime); err != nil {
		if err == xsql.ErrNoRows {
			m = nil
			err = nil
		}
	}
	return
}

// AddSubject insert or update subject state.
func (d *Dao) AddSubject(c context.Context, mid, oid int64, typ, state int32, now time.Time) (id int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_inSubjectSQL, subHit(oid)), oid, typ, mid, state, now, now, state, mid, now)
	if err != nil {
		return
	}
	return res.LastInsertId()
}

// UpSubjectState update subject state.
func (d *Dao) UpSubjectState(c context.Context, oid int64, typ, state int32, now time.Time) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_upSubStateSQL, subHit(oid)), state, now, oid, typ)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// UpSubjectAttr update subject attr.
func (d *Dao) UpSubjectAttr(c context.Context, oid int64, typ int32, attr uint32, now time.Time) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_upSubAttrSQL, subHit(oid)), attr, now, oid, typ)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// UpStateAndAttr update subject state and attr
func (d *Dao) UpStateAndAttr(c context.Context, oid int64, typ, state int32, attr uint32, now time.Time) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_upSubStateAndAttrSQL, subHit(oid)), state, attr, now, oid, typ)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// TxIncrSubCount incr subject count and rcount by transaction.
func (d *Dao) TxIncrSubCount(tx *xsql.Tx, oid int64, typ int32, now time.Time) (rows int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_incrSubCountSQL, subHit(oid)), now, oid, typ)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// TxIncrSubFCount incr subject count and rcount by transaction.
func (d *Dao) TxIncrSubFCount(tx *xsql.Tx, oid int64, typ int32, now time.Time) (rows int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_incrSubFCountSQL, subHit(oid)), now, oid, typ)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// TxIncrSubRCount incr subject rcount by transaction
func (d *Dao) TxIncrSubRCount(tx *xsql.Tx, oid int64, typ int32, now time.Time) (rows int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_incrSubRCountSQL, subHit(oid)), now, oid, typ)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// TxDecrSubRCount decr subject count by transaction.
func (d *Dao) TxDecrSubRCount(tx *xsql.Tx, oid int64, typ int32, now time.Time) (rows int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_decrSubRCountSQL, subHit(oid)), now, oid, typ)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// TxIncrSubACount incr subject acount by transaction.
func (d *Dao) TxIncrSubACount(tx *xsql.Tx, oid int64, typ int32, count int32, now time.Time) (rows int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_incrSubACountSQL, subHit(oid)), count, now, oid, typ)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// TxSubDecrACount decr subject rcount by transaction.
func (d *Dao) TxSubDecrACount(tx *xsql.Tx, oid int64, typ int32, count int32, now time.Time) (rows int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_decrSubACountSQL, subHit(oid)), count, now, oid, typ)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// TxSubDecrMCount decr subject mcount by transaction.
func (d *Dao) TxSubDecrMCount(tx *xsql.Tx, oid int64, typ int32, now time.Time) (rows int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_decrSubMCountSQL, subHit(oid)), now, oid, typ)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// TxUpSubAttr update subject attr.
func (d *Dao) TxUpSubAttr(tx *xsql.Tx, oid int64, tp int32, attr uint32, now time.Time) (rows int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_upSubAttrSQL, subHit(oid)), attr, now, oid, tp)
	if err != nil {
		log.Error("mysqlDB.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// TxUpSubMeta update subject meta.
func (d *Dao) TxUpSubMeta(tx *xsql.Tx, oid int64, tp int32, meta string, now time.Time) (rows int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_upSubMetaSQL, subHit(oid)), meta, now, oid, tp)
	if err != nil {
		log.Error("mysqlDB.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}
