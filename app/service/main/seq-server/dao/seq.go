package dao

import (
	"context"
	"database/sql"
	"time"

	"go-common/app/service/main/seq-server/model"
	"go-common/library/log"
)

const (
	_allSQL           = "SELECT id,max_seq,step,token,perch,ctime,mtime FROM business"
	_maxSeqSQL        = "SELECT max_seq FROM business WHERE id=?"
	_upMaxSeqSQL      = "UPDATE business SET max_seq=? WHERE id=? AND max_seq=?"
	_upMaxSeqTokenSQL = "UPDATE business SET max_seq=?,step=? WHERE id=? AND token=?"
)

// All get all seq
func (d *Dao) All(c context.Context) (bs map[int64]*model.Business, err error) {
	rows, err := d.db.Query(c, _allSQL)
	if err != nil {
		log.Error("d.db.Query(%s) error(%v)", _allSQL, err)
		return
	}
	bs = make(map[int64]*model.Business)
	defer rows.Close()
	for rows.Next() {
		b := new(model.Business)
		if err = rows.Scan(&b.ID, &b.MaxSeq, &b.Step, &b.Token, &b.Perch, &b.CTime, &b.MTime); err != nil {
			log.Error("rows.Scan error(%v)")
			return
		}
		b.BenchTime = b.CTime.Time().UnixNano() / int64(time.Millisecond)
		b.LastTimestamp = time.Now().UnixNano() / int64(time.Millisecond)
		bs[b.ID] = b
	}
	return
}

// MaxSeq return current max seq.
func (d *Dao) MaxSeq(c context.Context, businessID int64) (maxSeq int64, err error) {
	row := d.db.QueryRow(c, _maxSeqSQL, businessID)
	if err = row.Scan(&maxSeq); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// UpMaxSeq update max seq by businessID.
func (d *Dao) UpMaxSeq(c context.Context, businessID, maxSeq, lastSeq int64) (rows int64, err error) {
	res, err := d.db.Exec(c, _upMaxSeqSQL, maxSeq, businessID, lastSeq)
	if err != nil {
		log.Error("d.db.Exec(%s, %d, %d) error(%v)", _upMaxSeqSQL, maxSeq, businessID)
		return
	}
	return res.RowsAffected()
}

// UpMaxSeqToken update max seq by businessID and token.
func (d *Dao) UpMaxSeqToken(c context.Context, businessID, maxSeq, step int64, token string) (rows int64, err error) {
	res, err := d.db.Exec(c, _upMaxSeqTokenSQL, maxSeq, step, businessID, token)
	if err != nil {
		log.Error("d.db.Exec(%s, %d, %d) error(%v)", _upMaxSeqSQL, maxSeq, businessID)
		return
	}
	return res.RowsAffected()
}
